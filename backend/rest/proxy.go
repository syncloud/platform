package rest

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	redirect Config
}

type Config interface {
	ApiUrl() string
}

func NewProxy(redirect Config) *Proxy {
	return &Proxy{
		redirect: redirect,
	}
}

func (p *Proxy) ProxyRedirect() (*httputil.ReverseProxy, error) {
	redirectApiUrl := p.redirect.ApiUrl()
	redirectUrl, err := url.Parse(redirectApiUrl)
	if err != nil {
		fmt.Printf("proxy url error: %v", err)
		return nil, err
	}
	director := func(req *http.Request) {
		req.URL.Scheme = redirectUrl.Scheme
		req.URL.Host = redirectUrl.Host
		req.Host = redirectUrl.Host
	}
	return &httputil.ReverseProxy{Director: director}, nil
}

func (p *Proxy) ProxyImage() *httputil.ReverseProxy {
	host := "apps.syncloud.org"
	director := func(req *http.Request) {
		query := req.URL.Query()
		if !query.Has("channel") {
			return
		}
		if !query.Has("app") {
			return
		}
		req.URL.Scheme = "http"
		req.URL.RawQuery = ""
		req.URL.Host = host
		req.URL.Path = fmt.Sprintf("/releases/%s/images/%s-128.png", query.Get("channel"), query.Get("app"))
		req.Host = host
	}
	return &httputil.ReverseProxy{Director: director}
}

func (p *Proxy) ProxyImageFunc() func(http.ResponseWriter, *http.Request) {
	proxy := p.ProxyImage()
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}
