package handlers

import "net/http"

type Handler func(request *http.Request, headers http.Header) (string, error)
