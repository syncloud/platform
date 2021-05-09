package redirect

type Response struct {
	Message string `json:"message"`
}

type UserResponse struct {
	Data User `json:"data"`
}

type User struct {
	UpdateToken string `json:"update_token"`
}
