package rest

import (
	"errors"
	"github.com/syncloud/platform/log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CookiesStub struct {
}

func (c *CookiesStub) GetSessionUser(_ *http.Request) (string, error) {
	return "user", nil
}

type UserConfigStub struct {
}

func (u *UserConfigStub) IsActivated() bool {
	return true
}

func TestHandlerSuccess(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()

	m := NewMiddleware(&CookiesStub{}, &UserConfigStub{}, log.Default())
	m.Handle(func(req *http.Request) (interface{}, error) { return []string{"test"}, nil })(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"success":true,"data":["test"]}`, rr.Body.String())
}

func TestHandlerSuccessBoolData(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()
	m := NewMiddleware(&CookiesStub{}, &UserConfigStub{}, log.Default())
	m.Handle(func(req *http.Request) (interface{}, error) { return true, nil })(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"success":true,"data":true}`, rr.Body.String())
}

func TestHandlerSuccess_Percent(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()
	m := NewMiddleware(&CookiesStub{}, &UserConfigStub{}, log.Default())
	m.Handle(func(req *http.Request) (interface{}, error) { return []string{"test %123 "}, nil })(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"success":true,"data":["test %123 "]}`, rr.Body.String())
}

func TestHandlerFail(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()
	m := NewMiddleware(&CookiesStub{}, &UserConfigStub{}, log.Default())
	m.Handle(func(req *http.Request) (interface{}, error) { return nil, errors.New("error") })(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, `{"success":false,"message":"error"}`+"\n", rr.Body.String())
}

func TestBackupCreateFail(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()
	m := NewMiddleware(&CookiesStub{}, &UserConfigStub{}, log.Default())
	m.Handle(func(req *http.Request) (interface{}, error) { return nil, errors.New("error") })(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, `{"success":false,"message":"error"}`+"\n", rr.Body.String())
}
