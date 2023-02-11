package rest

import (
	"encoding/json"
	"fmt"
	"github.com/syncloud/platform/rest/model"
	"go.uber.org/zap"
	"net/http"
)

type Cookies interface {
	GetSessionUser(r *http.Request) (string, error)
}

type UserConfig interface {
	IsActivated() bool
}

type Middleware struct {
	cookies    Cookies
	userConfig UserConfig
	logger     *zap.Logger
}

func NewMiddleware(cookies Cookies, userConfig UserConfig, logger *zap.Logger) *Middleware {
	return &Middleware{
		cookies:    cookies,
		userConfig: userConfig,
		logger:     logger,
	}
}

func (m *Middleware) FailIfNotActivated(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !m.userConfig.IsActivated() {
			http.Error(w, "Device is not activated", 501)
			return
		}
		f(w, r)
	}
}

func (m *Middleware) FailIfActivated(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if m.userConfig.IsActivated() {
			http.Error(w, "Device is activated", 502)
			return
		}
		f(w, r)
	}
}

func (m *Middleware) SecuredHandle(f func(*http.Request) (interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return m.Secured(m.Handle(f))
}

func (m *Middleware) Secured(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := m.cookies.GetSessionUser(r)
		if err != nil {
			m.logger.Info("error", zap.Error(err))
			http.Error(w, "Unauthorized", 401)
			return
		}
		f(w, r)
	}
}

func (m *Middleware) JsonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info(fmt.Sprintf("%s: %s", r.Method, r.RequestURI))
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	m.logger.Warn(fmt.Sprintf("404 %s: %s", r.Method, r.RequestURI))
	http.NotFound(w, r)
}

func (m *Middleware) Handle(f func(*http.Request) (interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := f(r)
		if err != nil {
			m.Fail(w, err)
		} else {
			m.success(w, data)
		}
	}
}

func (m *Middleware) Fail(w http.ResponseWriter, err error) {
	m.logger.Error("error", zap.Error(err))
	response := model.Response{
		Success: false,
		Message: err.Error(),
	}
	statusCode := http.StatusInternalServerError
	switch v := err.(type) {
	case *model.ParameterError:
		m.logger.Warn("error", zap.Error(v))
		response.ParametersMessages = v.ParameterErrors
		statusCode = 400
	case *model.ServiceError:
		statusCode = v.StatusCode
	}
	responseJson, err := json.Marshal(response)
	responseText := ""
	if err != nil {
		responseText = err.Error()
	} else {
		responseText = string(responseJson)
	}
	http.Error(w, responseText, statusCode)
}

func (m *Middleware) success(w http.ResponseWriter, data interface{}) {
	response := model.Response{
		Success: true,
		Data:    &data,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		m.Fail(w, err)
	} else {
		_, _ = fmt.Fprint(w, string(responseJson))
	}
}
