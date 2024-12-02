package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/syncloud/platform/activation"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ManagedActivationStub struct {
	error bool
}

func (a *ManagedActivationStub) Activate(_ string, _ string, _ string, _ string, _ string) error {
	if a.error {
		return fmt.Errorf("error")
	}
	return nil
}

type CustomActivationStub struct{}

func (a CustomActivationStub) Activate(_ string, _ string, _ string) error {
	return nil
}

func TestActivate_CustomLoginShort(t *testing.T) {
	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{})
	request := &activation.CustomActivateRequest{Domain: "example.com", DeviceUsername: "a", DevicePassword: "password123"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Custom(req)
	assert.Nil(t, message)
	assert.NotNil(t, err)
}

func TestActivate_CustomPasswordShort(t *testing.T) {
	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{})
	request := &activation.CustomActivateRequest{Domain: "example.com", DeviceUsername: "username", DevicePassword: "pass"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Custom(req)
	assert.Nil(t, message)
	assert.NotNil(t, err)
}

func TestActivate_CustomGood(t *testing.T) {
	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{})
	request := &activation.CustomActivateRequest{Domain: "example.com", DeviceUsername: "username", DevicePassword: "password"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Custom(req)
	assert.Equal(t, "ok", message)
	assert.Nil(t, err)
}

func TestActivate_ManagedLoginShort(t *testing.T) {
	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{})
	request := &activation.ManagedActivateRequest{Domain: "example.com", DeviceUsername: "a", DevicePassword: "password"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Managed(req)
	assert.Nil(t, message)
	assert.NotNil(t, err)
}

func TestActivate_ManagedPasswordShort(t *testing.T) {
	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{})
	request := &activation.ManagedActivateRequest{Domain: "example.com", DeviceUsername: "username", DevicePassword: "pass"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Managed(req)
	assert.Nil(t, message)
	assert.NotNil(t, err)
}

func TestActivate_ManagedGood(t *testing.T) {
	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{})
	request := &activation.ManagedActivateRequest{Domain: "example.com", DeviceUsername: "username", DevicePassword: "password"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Managed(req)
	assert.Equal(t, "ok", message)
	assert.Nil(t, err)
}

func TestActivate_ManagedRedirectError(t *testing.T) {
	managed := &ManagedActivationStub{error: true}
	activate := NewActivateBackend(managed, &CustomActivationStub{})
	request := &activation.ManagedActivateRequest{Domain: "example.com", DeviceUsername: "username", DevicePassword: "password"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	_, err = activate.Managed(req)
	assert.NotNil(t, err)
}

func TestActivate_Managed_LowerCase(t *testing.T) {
	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{})
	request := &activation.ManagedActivateRequest{Domain: "example.com", DeviceUsername: "Boris@example.com", DevicePassword: "password"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Managed(req)
	assert.Nil(t, message)
	assert.NotNil(t, err)
}
