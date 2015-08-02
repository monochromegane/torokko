package main

type buildError struct {
	err error
}

func (e buildError) Error() string {
	return e.err.Error()
}

type aleadyExistsError struct {
}

func (e aleadyExistsError) Error() string {
	return "already exists"
}
