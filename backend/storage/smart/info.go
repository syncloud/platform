package smart

type Info struct {
	ExitStatus int       `json:"exit_status"`
	Messages   []Message `json:"messages"`
}
