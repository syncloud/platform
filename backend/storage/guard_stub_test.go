package storage

type GuardStub struct {
	err error
}

func (g *GuardStub) Check() error {
	return g.err
}
