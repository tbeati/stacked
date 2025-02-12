package a

type errStruct struct {
	err error
}

func errArgument(n int, err error) {
}

func F() error {
	return nil
}

type S struct{}

func (s *S) F() error {
	return nil
}
