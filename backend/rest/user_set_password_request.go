package rest

type UserSetPasswordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
