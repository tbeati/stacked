package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	config := Config{
		GeneratedPackages: []string{"testdata/generated"},
	}

	analysistest.Run(t, analysistest.TestData(), NewAnalyzer(&config), "testdata/a")
}
