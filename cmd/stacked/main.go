package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/beati/stacked/cmd/stacked/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
