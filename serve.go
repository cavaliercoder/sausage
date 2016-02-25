package main

import (
	"net/http"
)

func serveSausages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "sausage/0.0.1")
	w.Write([]byte("sausages\n"))
}
