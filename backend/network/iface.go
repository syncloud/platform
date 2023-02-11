package network

type Interface struct {
	Name      string   `json:"name"`
	Addresses []string `json:"addresses"`
}
