package redirect

type PassThroughJsonError struct {
	Message string
	Json    string
}

func (p *PassThroughJsonError) Error() string {
	return p.Message
}
