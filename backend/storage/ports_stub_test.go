package storage

type PortsStub struct {
	ports map[string]string
}

func (p *PortsStub) Port(device string) string {
	return p.ports[device]
}
