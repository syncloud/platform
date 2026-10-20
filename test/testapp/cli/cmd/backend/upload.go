package main

import (
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
)

type upload struct {
	received atomic.Int64
}

func (u *upload) Receive(w http.ResponseWriter, r *http.Request) {
	u.received.Store(0)
	buffer := make([]byte, 65536)
	for {
		n, err := r.Body.Read(buffer)
		u.received.Add(int64(n))
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	fmt.Fprintf(w, "%d", u.received.Load())
}

func (u *upload) Progress(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "%d", u.received.Load())
}
