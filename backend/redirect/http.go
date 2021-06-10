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
	var redirectResponse Response
	err := json.Unmarshal(body, &redirectResponse)
	if err != nil {
		log.Printf("error parsing redirect response: %v\n", err)
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
