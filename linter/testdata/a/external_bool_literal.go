package a

func boolLiteralAssignmentExternal() {
	var err boolError
	_ = err

	err = true //  want "^basic literal true is not wrapped with stacked$"

	_, err = 0, true // want "^basic literal true is not wrapped with stacked$"
}

func boolLiteralDeclarationExternal() {
	var err boolError = true // want "^basic literal true is not wrapped with stacked$"
	_ = err
}

func boolLiteralReturnSingleExternal() boolError {
	return true // want "^basic literal true is not wrapped with stacked$"
}

func boolLiteralReturnMultipleExternal() (int, boolError) {
	return 0, true // want "^basic literal true is not wrapped with stacked$"
}

func boolLiteralArgumentExternal() {
	functionWithIntBoolTypeErrorArgument(0, true) // want "^basic literal true is not wrapped with stacked$"
}

func boolLiteralCompositeLiteralExternal() {
	_ = structWithBoolTypeErrorField{
		err: true, // want "^basic literal true is not wrapped with stacked$"
	}

	_ = []boolError{true} // want "^basic literal true is not wrapped with stacked$"

	_ = map[string]boolError{"": true} // want "^basic literal true is not wrapped with stacked$"
}

func boolLiteralChannelSendExternal() {
	var errChan chan boolError

	errChan <- true // want "^basic literal true is not wrapped with stacked$"
}
