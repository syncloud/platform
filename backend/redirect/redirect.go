package redirect

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/identification"
	"github.com/syncloud/platform/network"
	"github.com/syncloud/platform/version"
	"go.uber.org/zap"
	"io"
	"log"
)

type Service struct {
	userConfig     *config.UserConfig
	identification *identification.Parser
	networkIface   *network.Interface
	certbotLogger  *zap.Logger
}

func New(userPlatformConfig *config.UserConfig, identification *identification.Parser, networkIface *network.Interface, certbotLogger *zap.Logger) *Service {
	return &Service{
		userConfig:     userPlatformConfig,
		identification: identification,
		networkIface:   networkIface,
		certbotLogger:  certbotLogger,
	}
}

func (r *Service) Authenticate(email string, password string) (*User, error) {
	request := &UserCredentials{Email: email, Password: password}
	url := fmt.Sprintf("%s/user", r.userConfig.GetRedirectApiUrl())
	body, err := r.postAndCheck(url, request)
	if err != nil {
		return nil, err
	}
	var redirectUserResponse UserResponse
	err = json.Unmarshal(*body, &redirectUserResponse)
	if err != nil {
		return nil, err
	}
	return &redirectUserResponse.Data, nil
}

func (r *Service) CertbotPresent(token, fqdn string, value ...string) error {
	request := &CertbotPresentRequest{Token: token, Fqdn: fqdn, Values: value}
	url := fmt.Sprintf("%s/certbot/present", r.userConfig.GetRedirectApiUrl())
	r.certbotLogger.Info(fmt.Sprintf("dns present: %s", url))
	_, err := r.postAndCheck(url, request)
	return err
}

func (r *Service) CertbotCleanUp(token, fqdn string) error {
	request := &CertbotCleanUpRequest{Token: token, Fqdn: fqdn}
	url := fmt.Sprintf("%s/certbot/cleanup", r.userConfig.GetRedirectApiUrl())
	r.certbotLogger.Info(fmt.Sprintf("dns cleanup: %s", url))
	_, err := r.postAndCheck(url, request)
	return err
}

func (r *Service) Acquire(email string, password string, domain string) (*Domain, error) {

	deviceId, err := r.identification.Id()
	if err != nil {
		return nil, err
	}

	request := &FreeDomainAcquireRequest{
		Email:            email,
		Password:         password,
		Domain:           domain,
		DeviceMacAddress: deviceId.MacAddress,
		DeviceName:       deviceId.Name,
		DeviceTitle:      deviceId.Title}
	url := fmt.Sprintf("%s/%s", r.userConfig.GetRedirectApiUrl(), "domain/acquire_v2")

	body, err := r.postAndCheck(url, request)
	if err != nil {
		return nil, err
	}
	var response FreeDomainAcquireResponse
	err = json.Unmarshal(*body, &response)
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, fmt.Errorf("failed to acquire domain")
	}
	return response.Data, nil
}

func (r *Service) Reset() error {
	return r.Update(nil, nil, true, false, true)
}

func (r *Service) Update(ipv4 *string, port *int, ipv4Enabled bool, ipv4Public bool, ipv6Enabled bool) error {

	platformVersion, err := version.PlatformVersion()
	if err != nil {
		return err
	}
	updateToken := r.userConfig.GetDomainUpdateToken()
	if updateToken == nil {
		return fmt.Errorf("domain update token is not evailable")
	}

	request := &FreeDomainUpdateRequest{
		Token:           *updateToken,
		PlatformVersion: platformVersion,
		WebProtocol:     config.WebProtocol,
		WebPort:         port,
		WebLocalPort:    config.WebAccessPort,
 Ipv4Enabled:      ipv4Enabled,
	Ipv6Enabled: ipv6Enabled,
	}

	if ipv4Enabled {
		localIpAddr, err := r.networkIface.LocalIPv4()
		if err != nil {
			return err
		}
		localIp := localIpAddr.String()
		request.LocalIp = &localIp

		if ipv4Public {
			if ipv4 == nil {
				ipv4, err := r.networkIface.PublicIPv4()
				if err != nil {
					return err
				}
				log.Printf("public ipv4: %s", ipv4)
			}
			request.Ip = ipv4
		}
		request.MapLocalAddress = !ipv4Public
	}

	if ipv6Enabled {
		ipv6Addr, err := r.networkIface.IPv6()
		if err == nil {
			ipv6 := ipv6Addr.String()
			request.Ipv6 = &ipv6
		}
	}

	dkimKey := r.userConfig.GetDkimKey()
	if dkimKey != nil {
		request.DkimKey = dkimKey
	}

	url := fmt.Sprintf("%s/%s", r.userConfig.GetRedirectApiUrl(), "domain/update")
	_, err = r.postAndCheck(url, request)
	return err
}

func (r *Service) postAndCheck(url string, request interface{}) (*[]byte, error) {
	requestJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	client := retryablehttp.NewClient()
	resp, err := client.Post(url, "application/json", requestJson)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = CheckHttpError(resp.StatusCode, body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

