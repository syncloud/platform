package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"

	"github.com/go-ldap/ldap/v3"
	"github.com/syncloud/platform/auth"
	"github.com/syncloud/platform/cli"
	zaplog "github.com/syncloud/platform/log"
	"go.uber.org/zap"
)

type LoginService struct {
	totp      *auth.TOTP
	elevation *auth.ElevationCode
	logger    *zap.Logger
}

type CredentialsRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func authenticateLDAP(username, password string) (string, error) {
	conn, err := ldap.DialURL("ldap://localhost:389")
	if err != nil {
		return "", fmt.Errorf("ldap connect: %w", err)
	}
	defer conn.Close()
	dn := fmt.Sprintf("cn=%s,ou=users,%s", username, auth.Domain)
	err = conn.Bind(dn, password)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	result, err := conn.Search(ldap.NewSearchRequest(dn, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 1, 0, false, "(objectClass=*)", []string{"mail"}, nil))
	if err != nil || len(result.Entries) == 0 {
		return "", fmt.Errorf("ldap lookup: %w", err)
	}
	return result.Entries[0].GetAttributeValue("mail"), nil
}

func (s *LoginService) authenticate(w http.ResponseWriter, r *http.Request) (*CredentialsRequest, string, bool) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return nil, "", false
	}
	var req CredentialsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.Password == "" {
		writeError(w, "username and password are required", http.StatusBadRequest)
		return nil, "", false
	}
	email, err := authenticateLDAP(req.Username, req.Password)
	if err != nil {
		writeError(w, "invalid credentials", http.StatusUnauthorized)
		return nil, "", false
	}
	return &req, email, true
}

func (s *LoginService) handleTOTPStatus(w http.ResponseWriter, r *http.Request) {
	req, _, ok := s.authenticate(w, r)
	if !ok {
		return
	}

	configured, err := s.totp.Has(req.Username)
	if err != nil {
		s.logger.Error("TOTP status check failed", zap.Error(err))
		writeError(w, "failed to check TOTP status", http.StatusInternalServerError)
		return
	}
	writeSuccess(w, map[string]bool{"configured": configured})
}

func (s *LoginService) handleTOTPSetup(w http.ResponseWriter, r *http.Request) {
	req, _, ok := s.authenticate(w, r)
	if !ok {
		return
	}

	s.logger.Info("generating TOTP", zap.String("username", req.Username))
	uri, err := s.totp.Generate(req.Username)
	if err != nil {
		s.logger.Error("TOTP generation failed", zap.Error(err))
		writeError(w, "failed to generate TOTP", http.StatusInternalServerError)
		return
	}
	writeSuccess(w, map[string]string{"uri": uri})
}

func (s *LoginService) handlePasskeyElevation(w http.ResponseWriter, r *http.Request) {
	req, email, ok := s.authenticate(w, r)
	if !ok {
		return
	}
	if email == "" {
		writeError(w, "user has no email", http.StatusBadRequest)
		return
	}

	code, err := s.elevation.Read(email)
	if err != nil {
		s.logger.Error("elevation code read failed", zap.String("username", req.Username), zap.Error(err))
		writeError(w, "failed to read elevation code", http.StatusInternalServerError)
		return
	}
	writeSuccess(w, ElevationCodeResponse{Code: code})
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Success: true, Data: data})
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Success: false, Message: message})
}

func main() {
	socketPath := "/var/snap/platform/current/login.sock"
	if len(os.Args) > 1 {
		socketPath = os.Args[1]
	}

	logger := zaplog.Default()

	dataDir := "/var/snap/platform/current"
	executor := cli.New(logger)

	svc := &LoginService{
		totp:      auth.NewTOTP(executor, dataDir, logger),
		elevation: auth.NewElevationCode(path.Join(dataDir, "authelia-notification.txt")),
		logger:    logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/login/totp/status", svc.handleTOTPStatus)
	mux.HandleFunc("/login/totp/setup", svc.handleTOTPSetup)
	mux.HandleFunc("/login/passkey/elevation", svc.handlePasskeyElevation)

	os.Remove(socketPath)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", socketPath, err)
	}
	defer listener.Close()

	logger.Info("login service started", zap.String("socket", socketPath))
	if err := http.Serve(listener, mux); err != nil {
		log.Fatal(err)
	}
}
