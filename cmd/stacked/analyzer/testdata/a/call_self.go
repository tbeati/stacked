package a

func callSelfFuncAssignment() error {
	err := F()
	if err != nil {
		return err
	}

	return F()
}

func callSelfMethodAssignment() error {
	s := S{}
	err := s.F()
	if err != nil {
		return err
	}

	return s.F()
}
