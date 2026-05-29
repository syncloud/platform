package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ConfigStub struct {
}

func (c *ConfigStub) ApiUrl() string {
	return "https://proxy"
}

type IconResolverStub struct {
	url string
	err error
}

func (s *IconResolverStub) AppImageUrl(string) (string, error) {
	return s.url, s.err
}

func TestProxy_ProxyImage(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write([]byte("PNGDATA"))
	}))
	defer upstream.Close()

	proxy := NewProxy(&ConfigStub{}, &IconResolverStub{url: upstream.URL})
	req := httptest.NewRequest(http.MethodGet, "/rest/proxy/image?channel=stable&app=test", nil)
	rec := httptest.NewRecorder()

	proxy.ProxyImageFunc()(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "PNGDATA", rec.Body.String())
	assert.Equal(t, "image/png", rec.Result().Header.Get("Content-Type"))
}

func TestProxy_ProxyImage_NotFound(t *testing.T) {
	proxy := NewProxy(&ConfigStub{}, &IconResolverStub{err: errors.New("not found")})
	req := httptest.NewRequest(http.MethodGet, "/rest/proxy/image?app=test", nil)
	rec := httptest.NewRecorder()

	proxy.ProxyImageFunc()(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestProxy_ProxyImage_MissingApp(t *testing.T) {
	proxy := NewProxy(&ConfigStub{}, &IconResolverStub{})
	req := httptest.NewRequest(http.MethodGet, "/rest/proxy/image", nil)
	rec := httptest.NewRecorder()

	proxy.ProxyImageFunc()(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProxy_ProxyRedirect(t *testing.T) {
	proxy := NewProxy(&ConfigStub{}, &IconResolverStub{})
	proxyUrl, err := url.Parse("https://localhost/test?a=b")
	assert.Nil(t, err)
	reverseProxy, err := proxy.ProxyRedirect()
	assert.Nil(t, err)
	req := &http.Request{URL: proxyUrl, Host: proxyUrl.Host}
	reverseProxy.Director(req)
	assert.Equal(t, "https://proxy/test?a=b", req.URL.String())
}
