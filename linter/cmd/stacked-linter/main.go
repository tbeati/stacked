package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/tbeati/stacked/linter"
)

func main() {
	var config linter.Config

	configContent, err := os.ReadFile("./stacked.json")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("error reading config file: %v", err)
	}
	if err == nil {
		err = json.Unmarshal(configContent, &config)
		if err != nil {
			log.Fatalf("error parsing config file: %v", err)
		}
	}

	singlechecker.Main(linter.NewAnalyzer(&config))
}
