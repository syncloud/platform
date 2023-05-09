package access

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/rest/model"
	"testing"
)

type PoptProbeStub struct {
	probed map[string]int
}

func NewPoptProbeStub() *PoptProbeStub {
	return &PoptProbeStub{probed: make(map[string]int)}
}

func (p *PoptProbeStub) Probe(ip string, port int) error {
	p.probed[ip]++
	return nil
}

type RedirectStub struct {
}

func (r *RedirectStub) Update(ipv4 *string, port *int, ipv4Enabled bool, ipv4Public bool, ipv6Enabled bool) error {
	return nil
}

type TriggerStub struct {
}

func (t TriggerStub) RunAccessChangeEvent() error {
	return nil
}

type NetworkInfoStub struct {
	publicIPv4 string
}

func (n *NetworkInfoStub) IPv6() (*string, error) {
	//TODO implement me
	panic("implement me")
}

func (n *NetworkInfoStub) PublicIPv4() (*string, error) {
	return &n.publicIPv4, nil
}

type ExternalAddressUserConfigStub struct {
	publicIp *string
}

func (u *ExternalAddressUserConfigStub) IsRedirectEnabled() bool {
	return true
}

func (u *ExternalAddressUserConfigStub) SetIpv4Enabled(enabled bool) {
}

func (u *ExternalAddressUserConfigStub) SetIpv4Public(enabled bool) {
}

func (u *ExternalAddressUserConfigStub) SetIpv6Enabled(enabled bool) {
}

func (u *ExternalAddressUserConfigStub) SetPublicIp(publicIp *string) {
	u.publicIp = publicIp
}

func (u *ExternalAddressUserConfigStub) SetPublicPort(port *int) {
}

func (u *ExternalAddressUserConfigStub) GetPublicIp() *string {
	//TODO implement me
	panic("implement me")
}

func (u *ExternalAddressUserConfigStub) GetPublicPort() *int {
	//TODO implement me
	panic("implement me")
}

func (u *ExternalAddressUserConfigStub) IsIpv6Enabled() bool {
	//TODO implement me
	panic("implement me")
}

func (u *ExternalAddressUserConfigStub) IsIpv4Public() bool {
	//TODO implement me
	panic("implement me")
}

func (u *ExternalAddressUserConfigStub) IsIpv4Enabled() bool {
	//TODO implement me
	panic("implement me")
}

func TestExternalAddress_UpdateWithIpv4(t *testing.T) {
	network := &NetworkInfoStub{publicIPv4: "2.2.2.2"}
	config := &ExternalAddressUserConfigStub{}
	probe := NewPoptProbeStub()
	access := New(probe, config, &RedirectStub{}, &TriggerStub{}, network, log.Default())
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
	assert.Equal(t, *config.publicIp, "1.1.1.1")
	assert.Equal(t, 1, probe.probed["1.1.1.1"])
}

func TestExternalAddress_UpdateWithInvalidIpv4_Reset(t *testing.T) {
	network := &NetworkInfoStub{publicIPv4: "2.2.2.2"}
	config := &ExternalAddressUserConfigStub{}
	probe := NewPoptProbeStub()
	access := New(probe, config, &RedirectStub{}, &TriggerStub{}, network, log.Default())
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
	assert.Nil(t, config.publicIp)
	assert.Equal(t, 1, probe.probed["2.2.2.2"])
}

func TestExternalAddress_Ipv4Private_NoProbe(t *testing.T) {
	network := &NetworkInfoStub{publicIPv4: "2.2.2.2"}
	config := &ExternalAddressUserConfigStub{}
	probe := NewPoptProbeStub()
	access := New(probe, config, &RedirectStub{}, &TriggerStub{}, network, log.Default())
	ip := ""
	port := 443
	request := model.Access{
		Ipv4:        &ip,
		Ipv4Enabled: true,
		Ipv4Public:  false,
		AccessPort:  &port,
		Ipv6Enabled: false,
	}
	err := access.Update(request)
	assert.Nil(t, err)
	assert.Nil(t, config.publicIp)
	assert.Equal(t, 0, len(probe.probed))
}
