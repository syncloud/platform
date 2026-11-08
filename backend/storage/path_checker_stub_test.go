package storage

type PathCheckerStub struct {
	exists bool
}

func (p *PathCheckerStub) ExternalDiskLinkExists() bool {
	return p.exists
}
