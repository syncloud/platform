package redirect

import (
	"encoding/json"
	"fmt"
	"github.com/syncloud/platform/config"
	"io"
	"log"
	"net/http"
)

type Redirect struct {
	UserPlatformConfig *config.PlatformUserConfig
}

func New(userPlatformConfig *config.PlatformUserConfig) *Redirect {
	return &Redirect{
		UserPlatformConfig: userPlatformConfig,
	}
}

func (r *Redirect) Authenticate(email string, password string) (*User, error) {
	url := fmt.Sprintf("%s/user/get?email=%s&password=%s", r.UserPlatformConfig.GetRedirectApiUrl(), email, password)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	err = CheckHttpError(resp.StatusCode, body)
	if err != nil {
		return nil, err
	}
	var redirectUserResponse UserResponse
	err = json.Unmarshal(body, &redirectUserResponse)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &redirectUserResponse.Data, nil
}

/*func Acquire(email string, password string, userDomain string) {
	uuid.NodeID()

	device_id = id.id()
	data =
	{
		'email': email,
		'password': password,
		'user_domain': user_domain,
		'device_mac_address': device_id.mac_address,
		'device_name': device_id.name,
		'device_title': device_id.title,
	}
	url = urljoin(self.user_platform_config.get_redirect_api_url(), "/domain/acquire")
	response = requests.post(url, data)
	util.check_http_error(response)
	response_data = convertible.from_json(response.text)
	return response_data
}*/
