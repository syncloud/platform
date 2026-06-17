package auth

type Group struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}
