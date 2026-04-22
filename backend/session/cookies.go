package session

import (
	"fmt"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"net/http"
)

const UserKey = "user"
const OIDCStateKey = "oidc_state"
const OIDCCodeVerifierKey = "oidc_code_verifier"

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
	c.store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (c *Cookies) getSession(r *http.Request) *sessions.Session {
	session, err := c.store.Get(r, "session")
	if err != nil {
		c.logger.Warn("session decode error, using new session", zap.Error(err))
	}
	return session
}

func (c *Cookies) SetSessionUser(w http.ResponseWriter, r *http.Request, user string) error {
	session := c.getSession(r)
	session.Values[UserKey] = user
	return session.Save(r, w)
}

func (c *Cookies) ClearSessionUser(w http.ResponseWriter, r *http.Request) error {
	r.Header.Del("Cookie")
	session := c.getSession(r)
	delete(session.Values, UserKey)
	return session.Save(r, w)
}

func (c *Cookies) GetSessionUser(r *http.Request) (string, error) {
	session := c.getSession(r)
	user, found := session.Values[UserKey]
	if !found {
		return "", fmt.Errorf("no session found")
	}
	return user.(string), nil
}


func (c *Cookies) SetOIDCState(w http.ResponseWriter, r *http.Request, state string, codeVerifier string) error {
	session := c.getSession(r)
	session.Values[OIDCStateKey] = state
	session.Values[OIDCCodeVerifierKey] = codeVerifier
	return session.Save(r, w)
}

func (c *Cookies) GetOIDCState(r *http.Request) (string, string, error) {
	session := c.getSession(r)
	state, found := session.Values[OIDCStateKey]
	if !found {
		return "", "", fmt.Errorf("no OIDC state found")
	}
	codeVerifier, found := session.Values[OIDCCodeVerifierKey]
	if !found {
		return "", "", fmt.Errorf("no OIDC code verifier found")
	}
	return state.(string), codeVerifier.(string), nil
}

func (c *Cookies) ClearOIDCState(w http.ResponseWriter, r *http.Request) error {
	session := c.getSession(r)
	delete(session.Values, OIDCStateKey)
	delete(session.Values, OIDCCodeVerifierKey)
	return session.Save(r, w)
}
