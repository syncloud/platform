package redirect

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/identification"
	"github.com/syncloud/platform/log"
)

type UserConfigStub struct {
}

func (u *UserConfigStub) GetUserUpdateToken() (string, error) {
	return "token", nil
}

func (u *UserConfigStub) GetRedirectApiUrl() string {
	return "url"
}

func (u *UserConfigStub) GetDomainUpdateToken() *string {
	s := "token"
	return &s
}

func (u *UserConfigStub) GetDkimKey() *string {
	key := "dkim"
	return &key
}

type IpParserStub struct {
}

func (i *IpParserStub) Id() (*identification.Id, error) {
	panic("implement me")
}

type NetInfoStub struct {
	publicIp string
}

func (n *NetInfoStub) LocalIPv4() (net.IP, error) {
	return net.IPv4zero, nil
}

func (n *NetInfoStub) IPv6() (*string, error) {
	ip := "[::1]"
	return &ip, nil
}

func (n *NetInfoStub) PublicIPv4() (*string, error) {
	return &n.publicIp, nil
}

type ClientStub struct {
	request string
}

func (c *ClientStub) Get(_ string) (*http.Response, error) {
	panic("implement me")
}

func (c *ClientStub) Post(_, _ string, body interface{}) (*http.Response, error) {
	c.request = string(body.([]byte))
	r := io.NopCloser(bytes.NewReader([]byte(`
{
	"success": true,
	"data": ""
}
`)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

type VersionStub struct {
}

func (v *VersionStub) Get() (string, error) {
	return "", nil
}

func TestUpdate_Ipv4And6Enabled(t *testing.T) {
	client := &ClientStub{}
	service := New(&UserConfigStub{}, &IpParserStub{}, &NetInfoStub{}, client, &VersionStub{}, log.Default())
	ipv4 := "1.1.1.1"
	port := 1
	err := service.Update(&ipv4, &port, true, true, true)
	assert.Nil(t, err)
	assert.Equal(t, `{"ip":"1.1.1.1","local_ip":"0.0.0.0","token":"token","ipv6":"[::1]","dkim_key":"dkim","platform_version":"","web_protocol":"https","web_local_port":443,"web_port":1,"ipv4_enabled":true,"ipv6_enabled":true}`, client.request)
}

func TestUpdate_Ipv4Disabled(t *testing.T) {
	client := &ClientStub{}
	service := New(&UserConfigStub{}, &IpParserStub{}, &NetInfoStub{}, client, &VersionStub{}, log.Default())
	ipv4 := "1.1.1.1"
	port := 1
	err := service.Update(&ipv4, &port, false, true, true)
	assert.Nil(t, err)
	assert.Equal(t, `{"token":"token","ipv6":"[::1]","dkim_key":"dkim","platform_version":"","web_protocol":"https","web_local_port":443,"web_port":1,"ipv4_enabled":false,"ipv6_enabled":true}`, client.request)
}

func TestUpdate_Ipv6Disabled(t *testing.T) {
	client := &ClientStub{}
	service := New(&UserConfigStub{}, &IpParserStub{}, &NetInfoStub{}, client, &VersionStub{}, log.Default())
	ipv4 := "1.1.1.1"
	port := 1
	err := service.Update(&ipv4, &port, true, true, false)
	assert.Nil(t, err)
	assert.Equal(t, `{"ip":"1.1.1.1","local_ip":"0.0.0.0","token":"token","dkim_key":"dkim","platform_version":"","web_protocol":"https","web_local_port":443,"web_port":1,"ipv4_enabled":true,"ipv6_enabled":false}`, client.request)
}

func TestUpdate_Ipv4AutoDetect(t *testing.T) {
	client := &ClientStub{}
	netInfo := &NetInfoStub{publicIp: "2.2.2.2"}
	service := New(&UserConfigStub{}, &IpParserStub{}, netInfo, client, &VersionStub{}, log.Default())
	port := 1
	err := service.Update(nil, &port, true, true, false)
	assert.Nil(t, err)
	assert.Equal(t, `{"ip":"2.2.2.2","local_ip":"0.0.0.0","token":"token","dkim_key":"dkim","platform_version":"","web_protocol":"https","web_local_port":443,"web_port":1,"ipv4_enabled":true,"ipv6_enabled":false}`, client.request)
}
