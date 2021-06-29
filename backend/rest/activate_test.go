package rest

import (
	"bytes"
	"encoding/json"
	"github.com/syncloud/platform/activation"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ActivateFreeStub struct{}

func (a *ActivateFreeStub) Activate(redirectEmail string, redirectPassword string, requestDomain string, deviceUsername string, devicePassword string) error {
	return nil
}

type ActivateCustomStub struct{}

func (a ActivateCustomStub) Activate(requestDomain string, deviceUsername string, devicePassword string) error {
	return nil
}

func TestActivate_CustomLoginShort(t *testing.T) {
	activate := NewActivateBackend(&ActivateFreeStub{}, &ActivateCustomStub{})
	request := &activation.CustomActivateRequest{Domain: "example.com", DeviceUsername: "a", DevicePassword: "password123"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Custom(req)
	assert.Nil(t, message)
	assert.NotNil(t, err)
}

func TestActivate_CustomPasswordShort(t *testing.T) {
	activate := NewActivateBackend(&ActivateFreeStub{}, &ActivateCustomStub{})
	request := &activation.CustomActivateRequest{Domain: "example.com", DeviceUsername: "username", DevicePassword: "pass"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Custom(req)
	assert.Nil(t, message)
	assert.NotNil(t, err)
}

func TestActivate_CustomGood(t *testing.T) {
	activate := NewActivateBackend(&ActivateFreeStub{}, &ActivateCustomStub{})
	request := &activation.CustomActivateRequest{Domain: "example.com", DeviceUsername: "username", DevicePassword: "password"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Custom(req)
	assert.Equal(t, "ok", message)
	assert.Nil(t, err)
}

func TestActivate_FreeLoginShort(t *testing.T) {
	activate := NewActivateBackend(&ActivateFreeStub{}, &ActivateCustomStub{})
	request := &activation.FreeActivateRequest{Domain: "example.com", DeviceUsername: "a", DevicePassword: "password"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Free(req)
	assert.Nil(t, message)
	assert.NotNil(t, err)
}

func TestActivate_FreePasswordShort(t *testing.T) {
	activate := NewActivateBackend(&ActivateFreeStub{}, &ActivateCustomStub{})
	request := &activation.FreeActivateRequest{Domain: "example.com", DeviceUsername: "username", DevicePassword: "pass"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Free(req)
	assert.Nil(t, message)
	assert.NotNil(t, err)
}

func TestActivate_FreeGood(t *testing.T) {
	activate := NewActivateBackend(&ActivateFreeStub{}, &ActivateCustomStub{})
	request := &activation.FreeActivateRequest{Domain: "example.com", DeviceUsername: "username", DevicePassword: "password"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Free(req)
	assert.Equal(t, "ok", message)
	assert.Nil(t, err)
}
