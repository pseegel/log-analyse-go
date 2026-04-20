package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

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

	for line := range lines {
		log.Println(line)
	}
}
