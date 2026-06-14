package linter

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	config := Config{
		IgnoredFunctions: []string{
			"testdata/generated.IgnoredFunction",
			"testdata/generated.IgnoredStruct.IgnoredMethod",
			"testdata/generated.IgnoredInterface.IgnoredMethod",
		},
		IgnoredTypes: []string{
			"testdata/generated.WrappedError",
			"testdata/generated.GenericWrappedError",
		},
		IgnoredInterfaces: []string{
			"testdata/b.IgnoredInterface",
		},
		GeneratedFiles: []string{
			"**/generated/*",
		},
		WrapChannelReceives: true,
	}

	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), NewAnalyzer(&config),
		"testdata/a",
	)
}
