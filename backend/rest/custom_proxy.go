package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/syncloud/platform/config"
	"net/http"
)

type CustomProxyNginx interface {
	InitCustomProxyConfig() error
}

type CustomProxy struct {
	config *config.UserConfig
	nginx  CustomProxyNginx
}

func NewCustomProxy(config *config.UserConfig, nginx CustomProxyNginx) *CustomProxy {
	return &CustomProxy{
		config: config,
		nginx:  nginx,
	}
}

func (cp *CustomProxy) List(_ *http.Request) (interface{}, error) {
	return cp.config.CustomProxies()
}

type customProxyAddRequest struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (cp *CustomProxy) Add(req *http.Request) (interface{}, error) {
	var request customProxyAddRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, errors.New("bad request")
	}
	if request.Name == "" || request.Host == "" || request.Port == 0 {
		return nil, errors.New("name, host and port are required")
	}
	err = cp.config.AddCustomProxy(request.Name, request.Host, request.Port)
	if err != nil {
		return nil, err
	}
	return "OK", cp.nginx.InitCustomProxyConfig()
}

type customProxyRemoveRequest struct {
	Name string `json:"name"`
}

func (cp *CustomProxy) Remove(req *http.Request) (interface{}, error) {
	var request customProxyRemoveRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return nil, errors.New("bad request")
	}
	if request.Name == "" {
		return nil, errors.New("name is required")
	}
	err = cp.config.RemoveCustomProxy(request.Name)
	if err != nil {
		return nil, err
	}
	return "OK", cp.nginx.InitCustomProxyConfig()
}
