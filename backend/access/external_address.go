package access

import (
	"fmt"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/rest/model"
	"go.uber.org/zap"
	"net"
)

const (
	CodeIpv4NotDetected  = "ipv4NotDetected"
	CodeIpv4NotReachable = "ipv4NotReachable"
	CodeIpv6NotAvailable = "ipv6NotAvailable"
	CodeIpv6NotReachable = "ipv6NotReachable"
)

type UserConfig interface {
	IsRedirectEnabled() bool
	SetIpv4Enabled(enabled bool)
	SetIpv4Public(enabled bool)
	SetIpv6Enabled(enabled bool)
	SetRelayEnabled(enabled bool)
	SetPublicIp(publicIp *string)
	SetPublicPort(port *int)
	GetPublicIp() *string
	GetPublicPort() *int
	IsIpv6Enabled() bool
	IsIpv4Public() bool
	IsIpv4Enabled() bool
	IsRelayEnabled() bool
}

type Redirect interface {
	Update(relay bool, ipv4 *string, port *int, ipv4Enabled bool, ipv4Public bool, ipv6Enabled bool) error
}

type Relay interface {
	Apply(relayEnabled bool) error
}

type Trigger interface {
	Trigger()
}

type NetworkInfo interface {
	IPv6() (*string, error)
	PublicIPv4() (*string, error)
}

type Response struct {
	Success bool    `json:"success"`
	Message *string `json:"message"`
}

type Probe interface {
	Probe(ip string, port int) error
}

type ExternalAddress struct {
	probe      Probe
	userConfig UserConfig
	redirect   Redirect
	relay      Relay
	trigger    Trigger
	network    NetworkInfo
	logger     *zap.Logger
}

func New(probe Probe, userConfig UserConfig, redirect Redirect, relay Relay, trigger Trigger, network NetworkInfo, logger *zap.Logger) *ExternalAddress {
	return &ExternalAddress{
		probe:      probe,
		userConfig: userConfig,
		redirect:   redirect,
		relay:      relay,
		trigger:    trigger,
		network:    network,
		logger:     logger,
	}
}

func (a *ExternalAddress) Update(request model.Access) error {

	a.logger.Info(fmt.Sprintf("update relay: %v, ipv4 enabled: %v, ipv4 public: %v, ipv6 enabled: %v",
		request.RelayEnabled, request.Ipv4Enabled, request.Ipv4Public, request.Ipv6Enabled))

	a.userConfig.SetRelayEnabled(request.RelayEnabled)

	if err := a.relay.Apply(request.RelayEnabled); err != nil {
		return err
	}

	ipv4 := request.Ipv4
	ipv4ToSave := ipv4
	if request.Ipv4Manual() {

		port := config.WebAccessPort
		if request.AccessPort != nil {
			port = *request.AccessPort
		}
		if ipv4 != nil {
			addr := net.ParseIP(*ipv4)
			if addr.To4() == nil {
				ipv4 = nil
				ipv4ToSave = nil
			}
		}

		if request.Ipv4Public {
			if ipv4 == nil {
				publicIp, err := a.network.PublicIPv4()
				if err != nil {
					return model.Coded(CodeIpv4NotDetected, err)
				}
				ipv4 = publicIp
			}
			err := a.probe.Probe(*ipv4, port)
			if err != nil {
				return model.Coded(CodeIpv4NotReachable, err)
			}
		}
	} else {
		ipv4 = nil
		ipv4ToSave = nil
	}

	if request.Ipv6Enabled {
		ipv6, err := a.network.IPv6()
		if err != nil {
			return model.Coded(CodeIpv6NotAvailable, err)
		}
		if ipv6 == nil {
			return model.Coded(CodeIpv6NotAvailable, fmt.Errorf("no ipv6 address on this device"))
		}
		err = a.probe.Probe(*ipv6, config.WebAccessPort)
		if err != nil {
			return model.Coded(CodeIpv6NotReachable, err)
		}
	}

	if a.userConfig.IsRedirectEnabled() {
		err := a.redirect.Update(
			request.RelayEnabled,
			ipv4,
			request.AccessPort,
			request.Ipv4Active(),
			request.Ipv4PublicDirect(),
			request.Ipv6Enabled)
		if err != nil {
			return err
		}
	}
	a.userConfig.SetIpv4Enabled(request.Ipv4Enabled)
	a.userConfig.SetIpv4Public(request.Ipv4Public)
	a.userConfig.SetPublicIp(ipv4ToSave)
	a.userConfig.SetIpv6Enabled(request.Ipv6Enabled)
	a.userConfig.SetPublicPort(request.AccessPort)

	if err := a.relay.Apply(request.RelayEnabled); err != nil {
		return err
	}

	a.trigger.Trigger()
	return nil

}

func (a *ExternalAddress) Sync() error {

	if a.userConfig.IsRedirectEnabled() {
		err := a.redirect.Update(
			a.userConfig.IsRelayEnabled(),
			a.userConfig.GetPublicIp(),
			a.userConfig.GetPublicPort(),
			a.userConfig.IsIpv4Enabled(),
			a.userConfig.IsIpv4Public(),
			a.userConfig.IsIpv6Enabled())
		if err != nil {
			return err
		}
	}
	return a.relay.Apply(a.userConfig.IsRelayEnabled())
}
