package handlers

import "net/http"

func (a *Handlers) Hello(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello-world"))
}
