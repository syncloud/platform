package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/activation"
	"github.com/syncloud/platform/timezone"
)

type fakeExecutor struct {
	calls [][]string
}

func (f *fakeExecutor) CombinedOutput(name string, arg ...string) ([]byte, error) {
	f.calls = append(f.calls, append([]string{name}, arg...))
	return nil, nil
}

type fakeTimezoneStore struct{ tz string }

func (f *fakeTimezoneStore) SetTimezone(tz string) { f.tz = tz }
func (f *fakeTimezoneStore) GetTimezone() string   { return f.tz }

func pickValidTimezone(t *testing.T) string {
	t.Helper()
	for _, name := range []string{"UTC", "Europe/London", "America/New_York"} {
		if info, err := os.Stat(filepath.Join("/usr/share/zoneinfo", name)); err == nil && !info.IsDir() {
			return name
		}
	}
	t.Skip("no zoneinfo database available on this host")
	return ""
}

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
	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{}, nil)
	request := &activation.CustomActivateRequest{Domain: "example.com", DeviceUsername: "a", DevicePassword: "password123"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Custom(req)
	assert.Nil(t, message)
	assert.NotNil(t, err)
}

func TestActivate_CustomPasswordShort(t *testing.T) {
	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{}, nil)
	request := &activation.CustomActivateRequest{Domain: "example.com", DeviceUsername: "username", DevicePassword: "pass"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Custom(req)
	assert.Nil(t, message)
	assert.NotNil(t, err)
}

func TestActivate_CustomGood(t *testing.T) {
	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{}, nil)
	request := &activation.CustomActivateRequest{Domain: "example.com", DeviceUsername: "username", DevicePassword: "password"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Custom(req)
	assert.Equal(t, "ok", message)
	assert.Nil(t, err)
}

func TestActivate_ManagedLoginShort(t *testing.T) {
	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{}, nil)
	request := &activation.ManagedActivateRequest{Domain: "example.com", DeviceUsername: "a", DevicePassword: "password"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Managed(req)
	assert.Nil(t, message)
	assert.NotNil(t, err)
}

func TestActivate_ManagedPasswordShort(t *testing.T) {
	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{}, nil)
	request := &activation.ManagedActivateRequest{Domain: "example.com", DeviceUsername: "username", DevicePassword: "pass"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Managed(req)
	assert.Nil(t, message)
	assert.NotNil(t, err)
}

func TestActivate_ManagedGood(t *testing.T) {
	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{}, nil)
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
	activate := NewActivateBackend(managed, &CustomActivationStub{}, nil)
	request := &activation.ManagedActivateRequest{Domain: "example.com", DeviceUsername: "username", DevicePassword: "password"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	_, err = activate.Managed(req)
	assert.NotNil(t, err)
}

func TestActivate_Managed_AppliesOptionalTimezone(t *testing.T) {
	tz := pickValidTimezone(t)
	exec := &fakeExecutor{}
	store := &fakeTimezoneStore{}
	applier := timezone.NewApplier(exec, store)

	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{}, applier)
	request := &activation.ManagedActivateRequest{
		Domain: "example.com", DeviceUsername: "username", DevicePassword: "password",
		Timezone: tz,
	}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))

	message, err := activate.Managed(req)

	assert.Nil(t, err)
	assert.Equal(t, "ok", message)
	assert.Equal(t, tz, store.tz)
	assert.Equal(t, [][]string{{"timedatectl", "set-timezone", tz}}, exec.calls)
}

func TestActivate_Managed_RejectsInvalidTimezone(t *testing.T) {
	exec := &fakeExecutor{}
	applier := timezone.NewApplier(exec, &fakeTimezoneStore{})

	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{}, applier)
	request := &activation.ManagedActivateRequest{
		Domain: "example.com", DeviceUsername: "username", DevicePassword: "password",
		Timezone: "Not/A/Real/Zone",
	}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))

	_, err = activate.Managed(req)

	assert.Error(t, err)
	assert.Empty(t, exec.calls, "timedatectl should not run for invalid tz")
}

func TestActivate_Managed_LowerCase(t *testing.T) {
	activate := NewActivateBackend(&ManagedActivationStub{}, &CustomActivationStub{}, nil)
	request := &activation.ManagedActivateRequest{Domain: "example.com", DeviceUsername: "Boris@example.com", DevicePassword: "password"}
	body, err := json.Marshal(request)
	assert.Nil(t, err)
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(body))
	message, err := activate.Managed(req)
	assert.Nil(t, message)
	assert.NotNil(t, err)
}
