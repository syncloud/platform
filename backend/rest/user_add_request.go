package rest

type UserAddRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Admin    bool   `json:"admin"`
}
