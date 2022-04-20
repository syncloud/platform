package redirect

import "time"

type Response struct {
	Message string `json:"message"`
}

type UserResponse struct {
	Data User `json:"data"`
}

type User struct {
	UpdateToken string `json:"update_token"`
}

type DomainAvailabilityRequest struct {
	Domain   *string `json:"domain,omitempty"`
	Password *string `json:"password,omitempty"`
	Email    *string `json:"email,omitempty"`
}

type FreeDomainAcquireRequest struct {
	Domain           string `json:"domain,omitempty"`
	Password         string `json:"password,omitempty"`
	Email            string `json:"email,omitempty"`
	DeviceMacAddress string `json:"device_mac_address,omitempty"`
	DeviceName       string `json:"device_name,omitempty"`
	DeviceTitle      string `json:"device_title,omitempty"`
}

type UserCredentials struct {
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type CertbotPresentRequest struct {
	Fqdn   string   `json:"fqdn,omitempty"`
	Token  string   `json:"token,omitempty"`
	Values []string `json:"values,omitempty"`
}

type CertbotCleanUpRequest struct {
	Fqdn  string `json:"fqdn,omitempty"`
	Token string `json:"token,omitempty"`
}

type FreeDomainUpdateRequest struct {
	Ip              *string `json:"ip,omitempty"`
	LocalIp         *string `json:"local_ip,omitempty"`
	MapLocalAddress bool    `json:"map_local_address,omitempty"`
	Token           string  `json:"token"`
	Ipv6            *string `json:"ipv6,omitempty"`
	DkimKey         *string `json:"dkim_key,omitempty"`
	PlatformVersion string  `json:"platform_version"`
	WebProtocol     string  `json:"web_protocol"`
	WebLocalPort    int     `json:"web_local_port"`
	WebPort         *int    `json:"web_port,omitempty"`
	Ipv4Enabled     bool    `json:"ipv4_enabled"`
	Ipv6Enabled     bool    `json:"ipv6_enabled"`
}

type FreeDomainAcquireResponse struct {
	Success bool    `json:"success"`
	Data    *Domain `json:"data,omitempty"`
}

type Domain struct {
	Name             string     `json:"name,omitempty"`
	Ip               *string    `json:"ip,omitempty"`
	Ipv6             *string    `json:"ipv6,omitempty"`
	DkimKey          *string    `json:"dkim_key,omitempty"`
	LocalIp          *string    `json:"local_ip,omitempty"`
	MapLocalAddress  bool       `json:"map_local_address,omitempty"`
	UpdateToken      string     `json:"update_token"`
	LastUpdate       *time.Time `json:"last_update,omitempty"`
	DeviceMacAddress *string    `json:"device_mac_address,omitempty"`
	DeviceName       *string    `json:"device_name,omitempty"`
	DeviceTitle      *string    `json:"device_title,omitempty"`
	PlatformVersion  *string    `json:"platform_version,omitempty"`
	WebProtocol      *string    `json:"web_protocol,omitempty"`
	WebPort          *int       `json:"web_port,omitempty"`
	WebLocalPort     *int       `json:"web_local_port,omitempty"`
}
