package access

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/rest/model"
	"testing"
)

type PoptProbeStub struct {
}

func (p PoptProbeStub) Probe(ip string, port int) error {
	return nil
}

type RedirectStub struct {
}

func (r RedirectStub) Update(ipv4 *string, port *int, ipv4Enabled bool, ipv4Public bool, ipv6Enabled bool) error {
	return nil
}

type TriggerStub struct {
}

func (t TriggerStub) RunAccessChangeEvent() error {
	return nil
}

type NetworkInfoStub struct {
	ipv4called bool
}

func (n NetworkInfoStub) IPv6() (*string, error) {
	//TODO implement me
	panic("implement me")
}

func (n *NetworkInfoStub) PublicIPv4() (*string, error) {
	n.ipv4called = true
	ip := "1.1.1.1"
	return &ip, nil
}

type ExternalAddressUserConfigStub struct {
}

func (u ExternalAddressUserConfigStub) IsRedirectEnabled() bool {
	return true
}

func (u ExternalAddressUserConfigStub) SetIpv4Enabled(enabled bool) {
}

func (u ExternalAddressUserConfigStub) SetIpv4Public(enabled bool) {
}

func (u ExternalAddressUserConfigStub) SetIpv6Enabled(enabled bool) {
}

func (u ExternalAddressUserConfigStub) SetPublicIp(publicIp *string) {
}

func (u ExternalAddressUserConfigStub) SetPublicPort(port *int) {
}

func (u ExternalAddressUserConfigStub) GetPublicIp() *string {
	//TODO implement me
	panic("implement me")
}

func (u ExternalAddressUserConfigStub) GetPublicPort() *int {
	//TODO implement me
	panic("implement me")
}

func (u ExternalAddressUserConfigStub) IsIpv6Enabled() bool {
	//TODO implement me
	panic("implement me")
}

func (u ExternalAddressUserConfigStub) IsIpv4Public() bool {
	//TODO implement me
	panic("implement me")
}

func (u ExternalAddressUserConfigStub) IsIpv4Enabled() bool {
	//TODO implement me
	panic("implement me")
}

func TestExternalAddress_UpdateWithIpv4(t *testing.T) {
	network := &NetworkInfoStub{}
	access := New(&PoptProbeStub{}, &ExternalAddressUserConfigStub{}, &RedirectStub{}, &TriggerStub{}, network, log.Default())
	ip := "1.1.1.1"
	port := 443
	request := model.Access{
		Ipv4:        &ip,
		Ipv4Enabled: true,
		Ipv4Public:  true,
		AccessPort:  &port,
		Ipv6Enabled: false,
	}
	err := access.Update(request)
	assert.Nil(t, err)
	assert.False(t, network.ipv4called)

}

func TestExternalAddress_UpdateWithInvalidIpv4_Reset(t *testing.T) {
	network := &NetworkInfoStub{}
	access := New(&PoptProbeStub{}, &ExternalAddressUserConfigStub{}, &RedirectStub{}, &TriggerStub{}, network, log.Default())
	ip := ""
	port := 443
	request := model.Access{
		Ipv4:        &ip,
		Ipv4Enabled: true,
		Ipv4Public:  true,
		AccessPort:  &port,
		Ipv6Enabled: false,
	}
	err := access.Update(request)
	assert.Nil(t, err)
	assert.True(t, network.ipv4called)

}

