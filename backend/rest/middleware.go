package rest

import (
	"encoding/json"
	"fmt"
	"github.com/syncloud/platform/rest/model"
	"net/http"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s: %s\n", r.Method, r.RequestURI)
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("404 %s: %s\n", r.Method, r.RequestURI)
	http.NotFound(w, r)
}

func Handle(f func(req *http.Request) (interface{}, error)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		data, err := f(req)
		if err != nil {
			fail(w, err)
		} else {
			success(w, data)
		}
	}
}

func fail(w http.ResponseWriter, err error) {
	fmt.Println("error: ", err)
	response := model.Response{
		Success: false,
		Message: err.Error(),
	}
	statusCode := http.StatusInternalServerError
	switch v := err.(type) {
	case *model.ParameterError:
		fmt.Println("parameter error: ", v.ParameterErrors)
		response.ParametersMessages = v.ParameterErrors
		statusCode = 400
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

func success(w http.ResponseWriter, data interface{}) {
	response := model.Response{
		Success: true,
		Data:    &data,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		fail(w, err)
	} else {
		_, _ = fmt.Fprint(w, string(responseJson))
	}
}
