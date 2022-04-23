package access

import (
	"encoding/json"
	"fmt"
	"github.com/syncloud/platform/http"
	"go.uber.org/zap"
	"io"
	"net"
)

type ProbeUserConfig interface {
	GetRedirectApiUrl() string
	GetDomainUpdateToken() *string
}

type PortProbe struct {
	userConfig ProbeUserConfig
	client     http.Client
	logger     *zap.Logger
}

func NewProbe(userConfig ProbeUserConfig, client http.Client, logger *zap.Logger) *PortProbe {
	return &PortProbe{
		userConfig: userConfig,
		client:     client,
		logger:     logger,
	}
}

func (p *PortProbe) Probe(ip *string, port int) error {
	p.logger.Info(fmt.Sprintf("probing %v", port))

	url := fmt.Sprintf("%s/%s", p.userConfig.GetRedirectApiUrl(), "probe/port_v3")
	token := p.userConfig.GetDomainUpdateToken()
	if token == nil {
		return fmt.Errorf("token is not set")
	}

	request := &PortProbeRequest{Token: *token, Port: port}
	if ip != nil {
		addr := net.ParseIP(*ip)
		if addr.IsPrivate() {
			return fmt.Errorf("IP: %v is not public", *ip)
		}
		request.Ip = ip
	}
	requestJson, err := json.Marshal(request)
	if err != nil {
		return err
	}

	response, err := p.client.Post(url, "application/json", requestJson)
	if err != nil {
		return err
	}
	p.logger.Info(fmt.Sprintf("response status: %v", response.StatusCode))
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	p.logger.Info(fmt.Sprintf("response text: %v", string(body)))

	var probeResponse Response
	err = json.Unmarshal(body, &probeResponse)
	if err != nil {
		return err
	}

	if !probeResponse.Success {
		message := "Unable to verify open ports"
		if probeResponse.Message != nil {
			message = fmt.Sprintf("%v, %v", message, *probeResponse.Message)
		}
		return fmt.Errorf(message)
	}

	return nil
}
