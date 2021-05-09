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

type DomainAvailabilityRequest struct {
	UserDomain *string `json:"user_domain,omitempty"`
	Password   *string `json:"password,omitempty"`
	Email      *string `json:"email,omitempty"`
}
