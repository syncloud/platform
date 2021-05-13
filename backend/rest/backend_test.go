package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerSuccess(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()
	Handle(func(req *http.Request) (interface{}, error) { return []string{"test"}, nil })(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"success":true,"data":["test"]}`, rr.Body.String())
}

func TestHandlerSuccessBoolData(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()
	Handle(func(req *http.Request) (interface{}, error) { return true, nil })(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"success":true,"data":true}`, rr.Body.String())
}

func TestHandlerFail(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()
	Handle(func(req *http.Request) (interface{}, error) { return nil, errors.New("error") })(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"success":false,"message":"error"}`, rr.Body.String())
}

func TestBackupCreateFail(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()
	Handle(func(req *http.Request) (interface{}, error) { return nil, errors.New("error") })(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"success":false,"message":"error"}`, rr.Body.String())
}
