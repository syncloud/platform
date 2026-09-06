package auth

type TokenResponse struct {
	IDToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
}
