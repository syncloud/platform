package rest

type UserSetEmailRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
