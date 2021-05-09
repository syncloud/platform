package redirect

import (
	"encoding/json"
	"github.com/syncloud/platform/util"
	"log"
)

func CheckHttpError(status int, body []byte) error {
	if status == 200 {
		return nil
	}
	log.Println(body)
	var redirectResponse Response
	err := json.Unmarshal(body, &redirectResponse)
	if err != nil {
		log.Println(err)
		return &util.PassThroughJsonError{
			Message: "Unable to parse Redirect response",
			Json:    string(body),
		}
	}
	return &util.PassThroughJsonError{
		Message: redirectResponse.Message,
		Json:    string(body),
	}
}
