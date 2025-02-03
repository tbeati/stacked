package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), NewAnalyzer([]string{"testdata/generated"}), "testdata/a")
}
