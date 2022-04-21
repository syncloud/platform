package access

import (
	"encoding/json"
	"fmt"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/event"
	"github.com/syncloud/platform/http"
	"github.com/syncloud/platform/redirect"
	"github.com/syncloud/platform/rest/model"
	"go.uber.org/zap"
	"io"
	"net"
)

type ExternalAddress struct {
	userConfig *config.UserConfig
	redirect   *redirect.Service
	trigger    *event.Trigger
	client     http.Client
	logger     *zap.Logger
}

type Response struct {
	Message string `json:"message"`
	Ip      string `json:"device_ip,omitempty"`
}

func New(userConfig *config.UserConfig, redirect *redirect.Service, trigger *event.Trigger, client http.Client, logger *zap.Logger) *ExternalAddress {
	return &ExternalAddress{
		userConfig: userConfig,
		redirect:   redirect,
		trigger:    trigger,
		client:     client,
		logger:     logger,
	}
}

func (a *ExternalAddress) Update(request model.Access) error {

	a.logger.Info(fmt.Sprintf("update ipv4: %v, ipb4 public: %v, ipv6: %v", request.Ipv4Enabled, request.Ipv4Public, request.Ipv6Enabled))
	if a.userConfig.IsRedirectEnabled() {
		err := a.redirect.Update(
			request.Ipv4,
			request.AccessPort,
			request.Ipv4Enabled,
			request.Ipv4Public,
			request.Ipv6Enabled)
		if err != nil {
			return err
		}
	}
	a.userConfig.SetIpv4Enabled(request.Ipv4Enabled)
	a.userConfig.SetIpv4Public(request.Ipv4Public)
	a.userConfig.SetPublicIp(request.Ipv4)
	a.userConfig.SetIpv6Enabled(request.Ipv6Enabled)
	a.userConfig.SetPublicPort(request.AccessPort)

	return a.trigger.RunAccessChangeEvent()

}

func (a *ExternalAddress) Sync() error {

	if a.userConfig.IsRedirectEnabled() {
		err := a.redirect.Update(
			a.userConfig.GetPublicIp(),
			a.userConfig.GetPublicPort(),
			a.userConfig.IsIpv4Enabled(),
			a.userConfig.IsIpv4Public(),
			a.userConfig.IsIpv6Enabled())
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *ExternalAddress) probe(ip *string, port int) error {
	a.logger.Info(fmt.Sprintf("probing %v", port))

	url := fmt.Sprintf("%s/%s", a.userConfig.GetRedirectApiUrl(), "probe/port_v3")
	token := a.userConfig.GetDomainUpdateToken()
	if token == nil {
		return fmt.Errorf("token is not set")
	}

	request := &PortProbeRequest{Token: *token, Port: port}
	if ip != nil {
		addr := net.ParseIP(*ip)
		if addr.IsPrivate() {
			return fmt.Errorf("IP: %v is not public", ip)
		}
		request.Ip = ip
	}
	requestJson, err := json.Marshal(request)
	if err != nil {
		return err
	}

	response, err := a.client.Post(url, "application/json", requestJson)
	if err != nil {
		return err
	}
	a.logger.Info(fmt.Sprintf("response status: %v", response.StatusCode))
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	a.logger.Info(fmt.Sprintf("response text: %v", string(body)))

	var probeResponse Response
	err = json.Unmarshal(body, &probeResponse)
	if err != nil {
		return err
	}

	if response.StatusCode == 200 && probeResponse.Message == "OK" {
		return nil
	} else {
		resultAddr := net.ParseIP(probeResponse.Ip)
		ipType := "4"
		resultIp := ""
		if resultAddr.To4() == nil {
			ipType = "6"
			resultIp = resultAddr.To16().String()
		} else {
			resultIp = resultAddr.To4().String()
		}

		return fmt.Errorf("using device public IP: '%v' which is IPv%v", resultIp, ipType)
	}

}
