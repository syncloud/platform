package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
)

type Api struct {
}

func NewApi() *Api {
	return &Api{}
}

func (b *Api) Start(network string, address string) {
	listener, err := net.Listen(network, address)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/app/install_path", Handle(b.AppInstallPath)).Methods("GET")
	r.HandleFunc("/app/data_path", Handle(b.AppDataPath)).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	r.Use(middleware)

	fmt.Println("Started api")
	_ = http.Serve(listener, r)

}

// TODO: Not sure if this is used

func (b *Api) AppInstallPath(req *http.Request) (interface{}, error) {
	keys, ok := req.URL.Query()["name"]
	if !ok {
		return nil, fmt.Errorf("no token")
	}
	return fmt.Sprintf("/snap/%s/current", keys[0]), nil
}

// TODO: Not sure if this is used

func (b *Api) AppDataPath(req *http.Request) (interface{}, error) {
	keys, ok := req.URL.Query()["name"]
	if !ok {
		return nil, fmt.Errorf("no token")
	}
	return fmt.Sprintf("/var/snap/%s/common", keys[0]), nil
}
