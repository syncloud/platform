package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-ldap/ldap/v3"
	"github.com/syncloud/platform/auth"
	"github.com/syncloud/platform/cli"
	zaplog "github.com/syncloud/platform/log"
	"go.uber.org/zap"
)

type LoginService struct {
	authelia *auth.Authelia
	logger   *zap.Logger
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

func authenticateLDAP(username, password string) error {
	conn, err := ldap.DialURL("ldap://localhost:389")
	if err != nil {
		return fmt.Errorf("ldap connect: %w", err)
	}
	defer conn.Close()
	err = conn.Bind(fmt.Sprintf("cn=%s,ou=users,%s", username, auth.Domain), password)
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}
	return nil
}

func (s *LoginService) authenticate(w http.ResponseWriter, r *http.Request) (*CredentialsRequest, bool) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return nil, false
	}
	var req CredentialsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.Password == "" {
		writeError(w, "username and password are required", http.StatusBadRequest)
		return nil, false
	}
	if err := authenticateLDAP(req.Username, req.Password); err != nil {
		writeError(w, "invalid credentials", http.StatusUnauthorized)
		return nil, false
	}
	return &req, true
}

func (s *LoginService) handleTOTPStatus(w http.ResponseWriter, r *http.Request) {
	req, ok := s.authenticate(w, r)
	if !ok {
		return
	}

	configured, err := s.authelia.HasTOTP(req.Username)
	if err != nil {
		s.logger.Error("TOTP status check failed", zap.Error(err))
		writeError(w, "failed to check TOTP status", http.StatusInternalServerError)
		return
	}
	writeSuccess(w, map[string]bool{"configured": configured})
}

func (s *LoginService) handleTOTPSetup(w http.ResponseWriter, r *http.Request) {
	req, ok := s.authenticate(w, r)
	if !ok {
		return
	}

	s.logger.Info("generating TOTP", zap.String("username", req.Username))
	uri, err := s.authelia.GenerateTOTP(req.Username)
	if err != nil {
		s.logger.Error("TOTP generation failed", zap.Error(err))
		writeError(w, "failed to generate TOTP", http.StatusInternalServerError)
		return
	}
	writeSuccess(w, map[string]string{"uri": uri})
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

	appDir := "/snap/platform/current"
	dataDir := "/var/snap/platform/current"
	configDir := appDir + "/config"
	executor := cli.New(logger)

	autheliaService := auth.NewAuthelia(
		configDir+"/authelia",
		dataDir+"/config/authelia",
		dataDir,
		dataDir+"/authelia.socket",
		nil, nil, nil, executor, nil, logger,
	)

	svc := &LoginService{
		authelia: autheliaService,
		logger:   logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/login/totp/status", svc.handleTOTPStatus)
	mux.HandleFunc("/login/totp/setup", svc.handleTOTPSetup)

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
