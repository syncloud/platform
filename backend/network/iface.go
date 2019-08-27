package network

import (
	"net"
 "fmt"
)

func LocalIp(isIpv6 bool) (net.IP, error) {
	ifaces, err := net.Interfaces()
 if err != nil {
  return nil, err
 }
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
  return nil, err
 }
		for _, addr := range addrs {
  var ip net.IP
  
  switch v := addr.(type) { 
   case *net.IPNet:
    ip = v.IP 
   case *net.IPAddr: 
    ip = v.IP 
  }
  if isIpv6 {
    ip = ip.To16()
  } else {
    ip = ip.To4()
  }
  if ip != nil && !ip.IsLoopback() {
   return ip, nil
  }
		}
	}
 return nil, fmt.Errorf("no address found")
}


