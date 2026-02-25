package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type OIDCConfig interface {
	DeviceUrl() string
	Url(app string) string
}

type OIDCService struct {
	config OIDCConfig
	logger *zap.Logger
}

func NewOIDCService(config OIDCConfig, logger *zap.Logger) *OIDCService {
	return &OIDCService{
		config: config,
		logger: logger,
	}
}

func (s *OIDCService) GetAuthorizationURL() (authURL string, state string, codeVerifier string, err error) {
	state, err = randomString(32)
	if err != nil {
		return "", "", "", fmt.Errorf("generate state: %w", err)
	}

	codeVerifier, err = randomString(64)
	if err != nil {
		return "", "", "", fmt.Errorf("generate code verifier: %w", err)
	}

	codeChallenge := generateCodeChallenge(codeVerifier)
	redirectURI := s.config.DeviceUrl() + "/rest/oidc/callback"
	authEndpoint := s.config.Url("auth") + "/api/oidc/authorization"

	params := url.Values{
		"client_id":             {"syncloud"},
		"response_type":        {"code"},
		"redirect_uri":         {redirectURI},
		"scope":                {"openid profile email groups"},
		"state":                {state},
		"code_challenge":       {codeChallenge},
		"code_challenge_method": {"S256"},
	}

	authURL = authEndpoint + "?" + params.Encode()
	return authURL, state, codeVerifier, nil
}

func (s *OIDCService) ExchangeCode(code string, codeVerifier string) (string, error) {
	tokenEndpoint := s.config.Url("auth") + "/api/oidc/token"
	redirectURI := s.config.DeviceUrl() + "/rest/oidc/callback"

	data := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {"syncloud"},
		"code":          {code},
		"redirect_uri":  {redirectURI},
		"code_verifier": {codeVerifier},
	}

	resp, err := http.PostForm(tokenEndpoint, data)
	if err != nil {
		return "", fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse struct {
		IDToken string `json:"id_token"`
	}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", fmt.Errorf("parse token response: %w", err)
	}

	if tokenResponse.IDToken == "" {
		return "", fmt.Errorf("no id_token in response")
	}

	username, err := extractUsernameFromIDToken(tokenResponse.IDToken)
	if err != nil {
		return "", err
	}

	return username, nil
}

func extractUsernameFromIDToken(idToken string) (string, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid id_token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("decode id_token payload: %w", err)
	}

	var claims struct {
		PreferredUsername string `json:"preferred_username"`
		Subject          string `json:"sub"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("parse id_token claims: %w", err)
	}

	username := claims.PreferredUsername
	if username == "" {
		username = claims.Subject
	}
	if username == "" {
		return "", fmt.Errorf("no username in id_token")
	}

	return username, nil
}

func randomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}

func generateCodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
