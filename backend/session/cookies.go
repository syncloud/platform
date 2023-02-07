package session

import (
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

const UserKey = "user"

type Cookies struct {
	config Config
	store  *sessions.CookieStore
}

type Config interface {
	GetWebSecretKey() (string, error)
}

func New(config Config) *Cookies {
	return &Cookies{
		config: config,
	}
}

func (c *Cookies) Start() error {
	secretKey, err := c.config.GetWebSecretKey()
	if err != nil {
		return err
	}
	c.store = sessions.NewCookieStore([]byte(secretKey))
	return nil
}

func (c *Cookies) getSession(r *http.Request) (*sessions.Session, error) {
	return c.store.Get(r, "session")
}

func (c *Cookies) setSessionUser(w http.ResponseWriter, r *http.Request, user string) error {
	session, err := c.getSession(r)
	if err != nil {
		return err
	}
	session.Values[UserKey] = user
	return session.Save(r, w)
}

func (c *Cookies) clearSessionUser(w http.ResponseWriter, r *http.Request) error {
	r.Header.Del("Cookie")
	session, err := c.getSession(r)
	if err != nil {
		return err
	}
	delete(session.Values, UserKey)
	return session.Save(r, w)
}

func (c *Cookies) GetSessionUser(r *http.Request) (string, error) {
	session, err := c.getSession(r)
	if err != nil {
		return "", err
	}
	user, found := session.Values[UserKey]
	if !found {
		return "", fmt.Errorf("no session found")
	}

	return user.(string), nil
}
