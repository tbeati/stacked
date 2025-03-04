package linter

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	config := Config{
		PackagesTreatedAsExternal: []string{"testdata/generated"},
		IgnoredFunctions: []string{
			"testdata/generated.IgnoredFunction",
			"testdata/generated.IgnoredStruct.IgnoredMethod",
		},
	}

	analysistest.Run(t, analysistest.TestData(), NewAnalyzer(&config), "testdata/a")
}
