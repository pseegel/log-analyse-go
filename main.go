package main

import (
	"context"
	"flag"
	"log"
	"time"
)

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

	_ = ctx // Hier könnte der Kontext in einer echten Implementierung verwendet werden

	log.Printf("würde verarbeiten: input=%s, output=%s", inputPath, outputPath)
}
