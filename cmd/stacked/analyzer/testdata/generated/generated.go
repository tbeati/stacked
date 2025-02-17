package generated

type StringError string

func (err StringError) Error() string {
	return string(err)
}

func F() error {
	return nil
}

type S struct{}

func (s *S) F() error {
	return nil
}
