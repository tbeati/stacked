package gclplugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/tbeati/stacked/linter"
)

func init() {
	register.Plugin("stacked", New)
}

type StackedPlugin struct {
	config *linter.Config
}

func New(settings any) (register.LinterPlugin, error) {
	config, err := register.DecodeSettings[linter.Config](settings)
	if err != nil {
		return nil, err
	}

	return &StackedPlugin{
		config: &config,
	}, nil
}

func (sp *StackedPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{linter.NewAnalyzer(sp.config)}, nil
}

func (sp *StackedPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
