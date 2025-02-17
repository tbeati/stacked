package analyzer

import "go/types"

var errorType *types.Interface

func init() {
	errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
}

func isError(t types.Type) bool {
	return types.Implements(t, errorType)
}
