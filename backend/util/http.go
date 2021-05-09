package util

import (
	"encoding/json"
	"github.com/syncloud/platform/redirect"
	"log"
)

func CheckHttpError(status int, body []byte) error {
	if status == 200 {
		return nil
	}
	log.Println(body)
	var redirectResponse redirect.Response
	err := json.Unmarshal(body, &redirectResponse)
	if err != nil {
		log.Println(err)
		return &PassThroughJsonError{
			message: "Unable to parse Redirect response",
			json:    string(body),
		}
	}
	return &PassThroughJsonError{
		message: redirectResponse.Message,
		json:    string(body),
	}
}
