package handlers

import (
	"hundred-board-games/code/pages"
	"hundred-board-games/code/templates"
	"net/http"
)

func HandleIndexPage(request *http.Request, headers http.Header) (string, error) {
	response, err := templates.RenderPage(pages.INDEX_PAGE, nil)
	if err != nil {
		return "", err
	}

	return response, nil
}
