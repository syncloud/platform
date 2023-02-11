package support

type Sender struct {
	aggregator *LogAggregator
	redirect   Redirect
}

type Redirect interface {
	SendLogs(logs string, includeSupport bool) error
}

func NewSender(aggregator *LogAggregator, redirect Redirect) *Sender {
	return &Sender{
		aggregator: aggregator,
		redirect:   redirect,
	}
}

func (s *Sender) Send(includeSupport bool) error {
	logs := s.aggregator.GetLogs()
	return s.redirect.SendLogs(logs, includeSupport)
}
