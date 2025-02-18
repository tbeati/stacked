package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/beati/stacked/cmd/stacked/analyzer"
)

func main() {
	configContent, err := os.ReadFile("./stacked.json")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("error reading config file: %v", err)
	}

	var config analyzer.Config
	err = json.Unmarshal(configContent, &config)
	if err != nil {
		log.Fatalf("error parsing config file: %v", err)
	}

	singlechecker.Main(analyzer.NewAnalyzer(&config))
}
