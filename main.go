package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const numWorkers = 5

type Result struct {
	Line   string
	Status int
}

func produce(ctx context.Context, path string, out chan<- string) error {

	defer close(out)

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		select {
		case out <- line:
			// gesendet, weiter
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan error: %w", err)
	}
	return nil
}

func worker(ctx context.Context, wg *sync.WaitGroup, in <-chan string, out chan<- Result) {
	defer wg.Done()

	for {
		select {
		case line, ok := <-in:
			if !ok {
				// Channel geschlossen, alle Zeilen gelesen
				return
			}
			// line verarbeiten
			result := Result{Line: line, Status: 0}
			select {
			case out <- result:
				// Ergebnis gesendet, weiter
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func main() {

	var (
		inputPath  string
		outputPath string
	)

	flag.StringVar(&inputPath, "input", "access.log", "Pfad zur einzulesenden Log Datei")
	flag.StringVar(&outputPath, "output", "report.json", "Pfad zur ausgegebenen JSON Datei")
	flag.Parse()

	log.Printf("input=%s, output=%s", inputPath, outputPath)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	lines := make(chan string, 100)

	go func() {
		if err := produce(ctx, inputPath, lines); err != nil {
			log.Printf("producer error: %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	results := make(chan Result, 100)

	for i := 0; i < numWorkers; i++ {
		go worker(ctx, &wg, lines, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		log.Println(result)
	}
}
