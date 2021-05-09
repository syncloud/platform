package util

type PassThroughJsonError struct {
	message string
	json    string
}

func (p *PassThroughJsonError) Error() string {
	return p.message
}
