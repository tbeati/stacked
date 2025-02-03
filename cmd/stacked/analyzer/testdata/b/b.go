package b

func F() error {
	return nil
}

type S struct{}

func (s *S) F() error {
	return nil
}
