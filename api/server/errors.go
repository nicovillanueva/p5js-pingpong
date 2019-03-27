package server

type errCannotSave struct{}

func (cs *errCannotSave) Error() string {
	return "could not save data"
}
