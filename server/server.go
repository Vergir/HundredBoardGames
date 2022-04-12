package server

import (
	"fmt"
	"net/http"
)

func AddHandler(path string, handler func(request *http.Request) string) {
	fullPath := "/" + path + "/"
	if path == PATH_INDEX {
		fullPath = "/"
	}
	http.HandleFunc(fullPath, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, handler(request))
	})
}

func Start() error {
	return http.ListenAndServe(":80", nil)
}
