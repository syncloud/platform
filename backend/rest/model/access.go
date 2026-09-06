package model

type Access struct {
	RelayEnabled bool    `json:"relay_enabled"`
	Ipv4         *string `json:"ipv4,omitempty"`
	Ipv4Enabled  bool    `json:"ipv4_enabled"`
	Ipv4Public   bool    `json:"ipv4_public"`
	AccessPort   *int    `json:"access_port,omitempty"`
	Ipv6Enabled  bool    `json:"ipv6_enabled"`
}

func (a Access) Ipv4Manual() bool {
	return !a.RelayEnabled && a.Ipv4Enabled
}

func (a Access) Ipv4Active() bool {
	return a.RelayEnabled || a.Ipv4Enabled
}

func (a Access) Ipv4PublicDirect() bool {
	return a.Ipv4Public && !a.RelayEnabled
}
