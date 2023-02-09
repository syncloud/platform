package session

import (
	"fmt"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"net/http"
)

const UserKey = "user"

type Cookies struct {
	config Config
	store  *sessions.CookieStore
	logger *zap.Logger
}

type Config interface {
	GetWebSecretKey() string
}

func New(config Config, logger *zap.Logger) *Cookies {
	return &Cookies{
		config: config,
		logger: logger,
	}
}

func (c *Cookies) Start() error {
	c.Reset()
	return nil
}

func (c *Cookies) Reset() {
	c.store = sessions.NewCookieStore([]byte(c.config.GetWebSecretKey()))
}

func (c *Cookies) getSession(r *http.Request) (*sessions.Session, error) {
	return c.store.Get(r, "session")
}

func (c *Cookies) SetSessionUser(w http.ResponseWriter, r *http.Request, user string) error {
	session, err := c.getSession(r)
	if err != nil {
		c.logger.Error("cannot update session", zap.Error(err))
		return err
	}
	session.Values[UserKey] = user
	return session.Save(r, w)
}

func (c *Cookies) ClearSessionUser(w http.ResponseWriter, r *http.Request) error {
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
