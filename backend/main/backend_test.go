package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
					"errors"

    "github.com/stretchr/testify/assert"

)

func TestHandlerGood(t *testing.T) {

    req, err := http.NewRequest("GET", "/health-check", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Handle(func() (interface{},error) {return []string{"test"}, nil }))
    handler.ServeHTTP(rr, req)

    assert.Equal(t, rr.Code, http.StatusOK, "wrong status")

    assert.Equal(t, rr.Body.String(), `{"success":true,"data":["test"]}`, "wrong body")
}

func TestHandlerBad(t *testing.T) {

    req, err := http.NewRequest("GET", "/health-check", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Handle(func() (interface{},error) {return nil, errors.New("error") }))
    handler.ServeHTTP(rr, req)

    assert.Equal(t, rr.Code, http.StatusOK, "wrong status")

    assert.Equal(t, rr.Body.String(), `{"success":false,"message":"Cannot get data"}`, "wrong body")
}
