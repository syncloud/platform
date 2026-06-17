package rest

type UserSetAdminRequest struct {
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
}
