package config

import "fmt"

const DefaultRedirectDomain = "syncloud.it"

type Redirect struct {
	db *Db
}

func NewRedirect(db *Db) *Redirect {
	return &Redirect{db: db}
}

func (r *Redirect) Domain() string {
	return r.db.Get("redirect.domain", DefaultRedirectDomain)
}

func (r *Redirect) SetDomain(domain string) {
	r.db.Delete("redirect.api_url")
	r.db.Upsert("redirect.domain", domain)
}

func (r *Redirect) ApiUrl() string {
	return r.db.Get("redirect.api_url", fmt.Sprintf("https://api.%s", r.Domain()))
}

func (r *Redirect) UpdateApiUrl(apiUrl string) {
	r.db.Upsert("redirect.api_url", apiUrl)
}

func (r *Redirect) UserEmail() *string {
	return r.db.GetOrNilString("redirect.user_email")
}

func (r *Redirect) SetUserEmail(email string) {
	r.db.Upsert("redirect.user_email", email)
}

func (r *Redirect) UserUpdateToken() (string, error) {
	return r.db.GetStringOrError("redirect.user_update_token")
}

func (r *Redirect) SetUserUpdateToken(token string) {
	r.db.Upsert("redirect.user_update_token", token)
}
