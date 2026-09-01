package access

import (
	"errors"
	"fmt"
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
	relay       bool
	ipv4Enabled bool
	ipv4Public  bool
	ipv6Enabled bool
	called      bool
}

type RelayStub struct {
	enabled  bool
	disabled bool
	applied  int
}

func (r *RelayStub) Apply(relayEnabled bool) error {
	r.applied++
	if relayEnabled {
		r.enabled = true
	} else {
		r.disabled = true
	}
	return nil
}

func (r *RedirectStub) Update(relay bool, ipv4 *string, port *int, ipv4Enabled bool, ipv4Public bool, ipv6Enabled bool) error {
	r.relay = relay
	r.ipv4Enabled = ipv4Enabled
	r.ipv4Public = ipv4Public
	r.ipv6Enabled = ipv6Enabled
	r.called = true
	return nil
}

type TriggerStub struct {
	triggered int
}

func (t *TriggerStub) Trigger() {
	t.triggered++
}

type NetworkInfoStub struct {
	publicIPv4    string
	ipv6Err       error
	publicIPv4Err error
}

func (n *NetworkInfoStub) IPv6() (*string, error) {
	if n.ipv6Err != nil {
		return nil, n.ipv6Err
	}
	ipv6 := "2a0d:3344:2d9:1b00::1"
	return &ipv6, nil
}

func (n *NetworkInfoStub) PublicIPv4() (*string, error) {
	if n.publicIPv4Err != nil {
		return nil, n.publicIPv4Err
	}
	return &n.publicIPv4, nil
}

type FailingProbeStub struct{}

func (p *FailingProbeStub) Probe(_ string, _ int) error {
	return fmt.Errorf("connection refused")
}

func assertCode(t *testing.T, err error, code string) {
	t.Helper()
	var coded *model.CodedError
	assert.True(t, errors.As(err, &coded), "expected a coded error, got %v", err)
	assert.Equal(t, code, coded.Code)
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
	return u.publicIp
}

func (u *ExternalAddressUserConfigStub) GetPublicPort() *int {
	return nil
}

func (u *ExternalAddressUserConfigStub) IsIpv6Enabled() bool {
	return false
}

func (u *ExternalAddressUserConfigStub) IsIpv4Public() bool {
	return false
}

func (u *ExternalAddressUserConfigStub) IsRelayEnabled() bool {
	return false
}

func (u *ExternalAddressUserConfigStub) SetRelayEnabled(enabled bool) {
}

func (u *ExternalAddressUserConfigStub) IsIpv4Enabled() bool {
	return true
}

func TestExternalAddress_UpdateWithIpv4(t *testing.T) {
	network := &NetworkInfoStub{publicIPv4: "2.2.2.2"}
	config := &ExternalAddressUserConfigStub{}
	probe := NewPoptProbeStub()
	access := New(probe, config, &RedirectStub{}, &RelayStub{}, &TriggerStub{}, network, log.Default())
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
	access := New(probe, config, &RedirectStub{}, &RelayStub{}, &TriggerStub{}, network, log.Default())
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

func TestExternalAddress_RelayWithIpv6_Coexist(t *testing.T) {
	network := &NetworkInfoStub{publicIPv4: "2.2.2.2"}
	config := &ExternalAddressUserConfigStub{}
	probe := NewPoptProbeStub()
	relay := &RelayStub{}
	redirect := &RedirectStub{}
	access := New(probe, config, redirect, relay, &TriggerStub{}, network, log.Default())
	request := model.Access{
		RelayEnabled: true,
		Ipv6Enabled:  true,
	}
	err := access.Update(request)
	assert.Nil(t, err)
	assert.True(t, relay.enabled)
	assert.False(t, relay.disabled)
	assert.True(t, redirect.called)
	assert.True(t, redirect.relay)
	assert.True(t, redirect.ipv4Enabled)
	assert.False(t, redirect.ipv4Public)
	assert.True(t, redirect.ipv6Enabled)
	assert.Equal(t, 1, probe.probed["2a0d:3344:2d9:1b00::1"])
}

func TestExternalAddress_RelayDisabled_CallsDisable(t *testing.T) {
	network := &NetworkInfoStub{publicIPv4: "2.2.2.2"}
	config := &ExternalAddressUserConfigStub{}
	relay := &RelayStub{}
	redirect := &RedirectStub{}
	access := New(NewPoptProbeStub(), config, redirect, relay, &TriggerStub{}, network, log.Default())
	err := access.Update(model.Access{RelayEnabled: false, Ipv6Enabled: true})
	assert.Nil(t, err)
	assert.True(t, relay.disabled)
	assert.False(t, redirect.relay)
	assert.True(t, redirect.ipv6Enabled)
}

func TestExternalAddress_Ipv4Private_NoProbe(t *testing.T) {
	network := &NetworkInfoStub{publicIPv4: "2.2.2.2"}
	config := &ExternalAddressUserConfigStub{}
	probe := NewPoptProbeStub()
	access := New(probe, config, &RedirectStub{}, &RelayStub{}, &TriggerStub{}, network, log.Default())
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

func TestExternalAddress_UpdateAppliesTheTunnelAgainAfterTheAddressUpdate(t *testing.T) {
	relay := &RelayStub{}
	access := New(NewPoptProbeStub(), &ExternalAddressUserConfigStub{}, &RedirectStub{}, relay,
		&TriggerStub{}, &NetworkInfoStub{publicIPv4: "2.2.2.2"}, log.Default())

	err := access.Update(model.Access{Ipv4Enabled: true, Ipv4Public: false})

	assert.Nil(t, err)
	assert.Equal(t, 2, relay.applied)
}

func TestExternalAddress_SyncAppliesTheTunnel(t *testing.T) {
	relay := &RelayStub{}
	access := New(NewPoptProbeStub(), &ExternalAddressUserConfigStub{}, &RedirectStub{}, relay,
		&TriggerStub{}, &NetworkInfoStub{publicIPv4: "2.2.2.2"}, log.Default())

	err := access.Sync()

	assert.Nil(t, err)
	assert.Equal(t, 1, relay.applied)
}

func TestExternalAddress_NoIpv6OnNetwork_SaysIpv6IsUnavailable(t *testing.T) {
	network := &NetworkInfoStub{publicIPv4: "2.2.2.2", ipv6Err: fmt.Errorf("network is unreachable")}
	access := New(NewPoptProbeStub(), &ExternalAddressUserConfigStub{}, &RedirectStub{}, &RelayStub{}, &TriggerStub{}, network, log.Default())

	err := access.Update(model.Access{Ipv6Enabled: true})

	assertCode(t, err, "ipv6NotAvailable")
}

func TestExternalAddress_Ipv6NotProbeable_SaysIpv6IsUnreachable(t *testing.T) {
	network := &NetworkInfoStub{publicIPv4: "2.2.2.2"}
	access := New(&FailingProbeStub{}, &ExternalAddressUserConfigStub{}, &RedirectStub{}, &RelayStub{}, &TriggerStub{}, network, log.Default())

	err := access.Update(model.Access{Ipv6Enabled: true})

	assertCode(t, err, "ipv6NotReachable")
}

func TestExternalAddress_Ipv4NotProbeable_SaysIpv4IsUnreachable(t *testing.T) {
	network := &NetworkInfoStub{publicIPv4: "2.2.2.2"}
	access := New(&FailingProbeStub{}, &ExternalAddressUserConfigStub{}, &RedirectStub{}, &RelayStub{}, &TriggerStub{}, network, log.Default())
	ip := "1.1.1.1"

	err := access.Update(model.Access{Ipv4: &ip, Ipv4Enabled: true, Ipv4Public: true})

	assertCode(t, err, "ipv4NotReachable")
}

func TestExternalAddress_NoPublicIpv4_SaysIpv4WasNotDetected(t *testing.T) {
	network := &NetworkInfoStub{publicIPv4Err: fmt.Errorf("lookup failed")}
	access := New(NewPoptProbeStub(), &ExternalAddressUserConfigStub{}, &RedirectStub{}, &RelayStub{}, &TriggerStub{}, network, log.Default())

	err := access.Update(model.Access{Ipv4Enabled: true, Ipv4Public: true})

	assertCode(t, err, "ipv4NotDetected")
}
