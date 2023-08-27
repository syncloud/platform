package network

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"
)

type TcpInterfaces struct {
}

type Interfaces interface {
	LocalIPv4() (net.IP, error)
	IPv6() (*string, error)
	PublicIPv4() (*string, error)
}

func New() *TcpInterfaces {
	return &TcpInterfaces{}
}

func (i *TcpInterfaces) LocalIPv4() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	addr := conn.LocalAddr().(*net.UDPAddr)
	return addr.IP, nil
}

func (i *TcpInterfaces) IPv6() (*string, error) {
	addr, err := i.IPv6Addr()
	if err != nil {
		return nil, err
	}
	ip := addr.String()
	return &ip, nil
}

func (i *TcpInterfaces) IPv6Addr() (net.IP, error) {
	conn, err := net.Dial("udp", "[2001:4860:4860::8888]:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	addr := conn.LocalAddr().(*net.UDPAddr)
	return addr.IP, nil
}

func (i *TcpInterfaces) PublicIPv4() (*string, error) {
	var zeroDialer net.Dialer
	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return zeroDialer.DialContext(ctx, "tcp4", addr)
	}
	httpClient.Transport = transport
	//url := "https://api.ipify.org?format=text"
	url := "https://myexternalip.com/raw"
	// http://api.ident.me
	// http://whatismyipaddress.com/api
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ip := string(ipBytes)
	return &ip, nil
}

func (i *TcpInterfaces) List() ([]Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var result []Interface
	for _, iface := range interfaces {
		if iface.Name == "lo" {
			continue
		}
		var addresses []string
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			addresses = append(addresses, addr.String())
		}
		result = append(result, Interface{
			Name:      iface.Name,
			Addresses: addresses,
		})
	}
	return result, nil

}
