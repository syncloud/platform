package auth

type IDTokenClaims struct {
	PreferredUsername string `json:"preferred_username"`
	Subject           string `json:"sub"`
}
