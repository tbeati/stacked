package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/beati/stacked/cmd/stacked/analyser"
)

func Test(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "a")
}
