package access

type PortProbeRequest struct {
	Token string  `json:"token,omitempty"`
	Port  int     `json:"port,omitempty"`
	Ip    *string `json:"ip,omitempty"`
}
