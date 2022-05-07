package network

import (
	"io/ioutil"
	"net"
	"net/http"
)

type Interface struct {
}

type Info interface {
	LocalIPv4() (net.IP, error)
	IPv6() (*string, error)
	PublicIPv4() (*string, error)
}

func New() *Interface {
	return &Interface{}
}

func (i *Interface) LocalIPv4() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

func (i *Interface) IPv6() (*string, error) {
	addr, err := i.IPv6Addr()
	if err != nil {
		return nil, err
	}
	ip := addr.String()
	return &ip, nil
}

func (i *Interface) IPv6Addr() (net.IP, error) {
	conn, err := net.Dial("udp", "[2001:4860:4860::8888]:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

func (i *Interface) PublicIPv4() (*string, error) {
	//url := "https://api.ipify.org?format=text"
	url := "https://myexternalip.com/raw"
	// http://api.ident.me
	// http://whatismyipaddress.com/api
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ipBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ip := string(ipBytes)
	return &ip, nil
}
