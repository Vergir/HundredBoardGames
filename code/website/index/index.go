package index

import (
	"hundred-board-games/code/templates"
	"hundred-board-games/code/utils"
	"hundred-board-games/code/website"
	"net/http"
)

var ENDPOINT = website.NewEndpoint("index", utils.StaticJs("common.js"), "index.css")
var HANDLER = handleRequest

func handleRequest(request *http.Request, headers http.Header) (string, error) {
	response, err := templates.RenderEndpoint(ENDPOINT, nil)
	if err != nil {
		return "", err
	}

	return response, nil
}
