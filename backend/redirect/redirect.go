package redirect

import (
	"encoding/json"
	"fmt"
	"github.com/syncloud/platform/config"
	"github.com/syncloud/platform/util"
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

func (r *Redirect) Login(email string, password string) (*User, error) {
	url := fmt.Sprintf("%s/user/get?email=%s&password=%s", r.userPlatformConfig.GetRedirectApiUrl(), email, password)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	err = util.CheckHttpError(resp.StatusCode, body)
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
