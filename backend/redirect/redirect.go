package redirect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/rest/model"
	"io"
	"log"
	"net/http"
)

type Redirect struct {
	userPlatformConfig *config.PlatformUserConfig
}

func New(userPlatformConfig *config.PlatformUserConfig) *Redirect {
	return &Redirect{
		userPlatformConfig: userPlatformConfig,
	}
}

func (r *Redirect) DomainAvailability(request model.RedirectCheckFreeDomainRequest) error {
	url := fmt.Sprintf("%s/domain/availability", r.userPlatformConfig.GetRedirectApiUrl())
	requestJson, err := json.Marshal(request)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestJson))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return CheckHttpError(resp.StatusCode, body)
}

func (r *Redirect) authenticate(email string, password string) (*User, error) {
	url := fmt.Sprintf("%s/user/get?email=%s&password=%s", r.userPlatformConfig.GetRedirectApiUrl(), email, password)
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
