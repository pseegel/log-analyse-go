package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	numWorkers = 5
	numFields  = 6
)

type Result struct {
	Line   string
	Status int
}

type Report struct {
	StatusCounts map[int]int `json:"status_counts"`
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

func parseLine(line string) (Result, error) {

	fields := strings.Fields(line)
	if len(fields) != numFields {
		return Result{Line: line}, fmt.Errorf("malformed line %q: expected %d fields, got %d", line, numFields, len(fields))
	}

	status, err := strconv.Atoi(fields[4])
	if err != nil {
		return Result{Line: line}, fmt.Errorf("status not valid: %w", err)
	}

	result := Result{
		Line:   line,
		Status: status,
	}

	return result, nil
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
			result, err := parseLine(line)
			if err != nil {
				log.Printf("parse error: %v", err)
				continue
			}
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

func collect(results <-chan Result, out string) error {
	report := Report{
		StatusCounts: make(map[int]int),
	}
	for result := range results {
		report.StatusCounts[result.Status]++
	}

	writer, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer writer.Close()

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("cannot write report: %w", err)
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

	err := collect(results, outputPath)
	if err != nil {
		log.Fatalf("collect error: %v", err)
	}
}
