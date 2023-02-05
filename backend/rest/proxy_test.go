package rest

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ConfigStub struct {
}

func (c *ConfigStub) GetRedirectApiUrl() string {
	return "https://proxy"
}

func TestProxy_ProxyImage(t *testing.T) {
	proxy := NewProxy(&ConfigStub{})
	imageUrl, err := url.Parse("https://localhost/image?channel=stable&app=test")
	assert.Nil(t, err)
	reverseProxy := proxy.ProxyImage()
	req := &http.Request{URL: imageUrl}
	reverseProxy.Director(req)
	assert.Equal(t, "http://apps.syncloud.org/releases/stable/images/test-128.png", req.URL.String())
}

func TestProxy_ProxyRedirect(t *testing.T) {
	proxy := NewProxy(&ConfigStub{})
	proxyUrl, err := url.Parse("https://localhost/test?a=b")
	assert.Nil(t, err)
	reverseProxy, err := proxy.ProxyRedirect()
	assert.Nil(t, err)
	req := &http.Request{URL: proxyUrl, Host: proxyUrl.Host}
	reverseProxy.Director(req)
	assert.Equal(t, "https://proxy/test?a=b", req.URL.String())
}
