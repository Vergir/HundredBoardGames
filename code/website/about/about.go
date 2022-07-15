package about

import (
	"hundred-board-games/code/templates"
	"hundred-board-games/code/website"
	"net/http"
)

var ENDPOINT = website.NewEndpoint("about", "about.css")
var HANDLER = handleRequest

func handleRequest(request *http.Request, headers http.Header) (string, error) {
	props := buildPageProps()

	response, err := templates.RenderEndpoint(ENDPOINT, props)
	if err != nil {
		return "", err
	}

	return response, nil
}
