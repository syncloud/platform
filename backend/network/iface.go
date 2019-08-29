package network

import (
	"net"
)

func LocalIPv4() (net.IP, error) {
 conn, err := net.Dial("udp", "8.8.8.8:80")
 if err != nil { 
  return nil, err
 }
 defer conn.Close()
 localAddr := conn.LocalAddr().(*net.UDPAddr)
 return localAddr.IP, nil
}

func LocalIPv6() (net.IP, error) {
 conn, err := net.Dial("udp", "[2001:4860:4860::8888]:80")
 if err != nil { 
  return nil, err
 }
 defer conn.Close()
 localAddr := conn.LocalAddr().(*net.UDPAddr)
 return localAddr.IP, nil
}


