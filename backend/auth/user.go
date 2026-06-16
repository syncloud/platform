package auth

type User struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Admin    bool     `json:"admin"`
	Groups   []string `json:"groups"`
}
