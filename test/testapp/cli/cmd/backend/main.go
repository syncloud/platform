package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func newHTTPClient() *http.Client {
	certPool, _ := x509.SystemCertPool()
	if certPool == nil {
		certPool = x509.NewCertPool()
	}
	caCert, err := os.ReadFile("/var/snap/platform/current/syncloud.ca.crt")
	if err == nil {
		certPool.AppendCertsFromPEM(caCert)
		block, _ := pem.Decode(caCert)
		if block != nil {
			if parsed, e := x509.ParseCertificate(block.Bytes); e == nil {
				fmt.Fprintf(os.Stderr, "backend: ca serial: %s\n", parsed.SerialNumber.String())
			}
		}
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: certPool},
		},
	}
}

const configDir = "/var/snap/testapp/current/config"

func readConfig(name string) string {
	data, err := os.ReadFile(configDir + "/" + name)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func main() {
	socketPath := os.Args[1]

	os.Remove(socketPath)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen: %v\n", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	var savedState, savedVerifier string

	mux.HandleFunc("/oidc/login", func(w http.ResponseWriter, r *http.Request) {
		authUrl := readConfig("auth_url")
		appUrl := readConfig("app_url")

		state, _ := randomString(32)
		verifier, _ := randomString(64)
		savedState = state
		savedVerifier = verifier

		challenge := generateCodeChallenge(verifier)
		redirectURI := appUrl + "/oidc/callback"
		authEndpoint := authUrl + "/api/oidc/authorization"

		params := url.Values{
			"client_id":             {"testapp"},
			"response_type":         {"code"},
			"redirect_uri":          {redirectURI},
			"scope":                 {"openid profile email"},
			"state":                 {state},
			"code_challenge":        {challenge},
			"code_challenge_method": {"S256"},
		}
		http.Redirect(w, r, authEndpoint+"?"+params.Encode(), http.StatusFound)
	})

	mux.HandleFunc("/oidc/callback", func(w http.ResponseWriter, r *http.Request) {
		authUrl := readConfig("auth_url")
		appUrl := readConfig("app_url")
		clientSecret := readConfig("client_secret")

		if errMsg := r.URL.Query().Get("error"); errMsg != "" {
			http.Error(w, fmt.Sprintf("oidc error: %s: %s", errMsg, r.URL.Query().Get("error_description")), http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		if code == "" || state != savedState {
			http.Error(w, "invalid state or missing code", http.StatusBadRequest)
			return
		}

		redirectURI := appUrl + "/oidc/callback"
		data := url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {redirectURI},
			"code_verifier": {savedVerifier},
		}

		tokenURL := authUrl + "/api/oidc/token"
		req, _ := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth("testapp", clientSecret)
		resp, err := newHTTPClient().Do(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("token request failed: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("token endpoint error %d: %s", resp.StatusCode, string(body)), http.StatusInternalServerError)
			return
		}

		var tokenResp struct {
			IDToken     string `json:"id_token"`
			AccessToken string `json:"access_token"`
		}
		json.Unmarshal(body, &tokenResp)

		var username string
		if tokenResp.AccessToken != "" {
			userinfoURL := authUrl + "/api/oidc/userinfo"
			uReq, _ := http.NewRequest("GET", userinfoURL, nil)
			uReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
			uResp, uErr := newHTTPClient().Do(uReq)
			if uErr == nil {
				defer uResp.Body.Close()
				uBody, _ := io.ReadAll(uResp.Body)
				var userInfo struct {
					PreferredUsername string `json:"preferred_username"`
				}
				json.Unmarshal(uBody, &userInfo)
				username = userInfo.PreferredUsername
			}
		}
		if username == "" {
			username, err = extractUsername(tokenResp.IDToken)
			if err != nil {
				http.Error(w, fmt.Sprintf("extract username: %v", err), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "OK %s", username)
	})

	http.Serve(listener, mux)
}

func extractUsername(idToken string) (string, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid id_token")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}
	var claims struct {
		PreferredUsername string `json:"preferred_username"`
		Subject          string `json:"sub"`
	}
	json.Unmarshal(payload, &claims)
	if claims.PreferredUsername != "" {
		return claims.PreferredUsername, nil
	}
	return claims.Subject, nil
}

func randomString(length int) (string, error) {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
