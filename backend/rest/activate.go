package rest

import (
	"encoding/json"
	"github.com/syncloud/platform/activation"
	"github.com/syncloud/platform/rest/model"
	"net/http"
)

type Activate struct {
	managed activation.ManagedActivation
	custom  activation.CustomActivation
}

func NewActivateBackend(managed activation.ManagedActivation, custom activation.CustomActivation) *Activate {
	return &Activate{
		managed: managed,
		custom:  custom,
	}
}

func (a *Activate) Custom(req *http.Request) (interface{}, error) {
	var request activation.CustomActivateRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	err = validate(request.DeviceUsername, request.DevicePassword)
	if err != nil {
		return nil, err
	}
	return "ok", a.custom.Activate(request.Domain, request.DeviceUsername, request.DevicePassword)
}

func (a *Activate) Managed(req *http.Request) (interface{}, error) {
	var request activation.ManagedActivateRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	err = validate(request.DeviceUsername, request.DevicePassword)
	if err != nil {
		return nil, err
	}
	return "ok", a.managed.Activate(request.RedirectEmail, request.RedirectPassword, request.Domain, request.DeviceUsername, request.DevicePassword)
}

func validate(username string, password string) error {
	if len(username) < 3 {
		return model.SingleParameterError("device_username", "less than 3 characters")
	}
	if len(password) < 7 {
		return model.SingleParameterError("device_password", "less than 7 characters")
	}
	return nil
}
