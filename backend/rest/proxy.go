package rest

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	redirect Config
	icons    IconResolver
	client   *http.Client
}

type Config interface {
	ApiUrl() string
}

type IconResolver interface {
	AppImageUrl(app string) (string, error)
}

func NewProxy(redirect Config, icons IconResolver) *Proxy {
	return &Proxy{
		redirect: redirect,
		icons:    icons,
		client:   http.DefaultClient,
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

func (p *Proxy) ProxyImageFunc() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		app := r.URL.Query().Get("app")
		if app == "" {
			http.Error(w, "app is required", http.StatusBadRequest)
			return
		}
		imageUrl, err := p.icons.AppImageUrl(app)
		if err != nil {
			http.Error(w, "icon not found", http.StatusNotFound)
			return
		}
		resp, err := p.client.Get(imageUrl)
		if err != nil {
			http.Error(w, "icon unavailable", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		if contentType := resp.Header.Get("Content-Type"); contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	}
}
