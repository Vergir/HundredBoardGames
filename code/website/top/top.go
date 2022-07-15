package top

import (
	"hundred-board-games/code/datamining"
	"hundred-board-games/code/games"
	"hundred-board-games/code/templates"
	"hundred-board-games/code/utils"
	"hundred-board-games/code/website"
	"net/http"
)

var ENDPOINT = website.NewComplexEndpoint(
	"top", "list.tmpl", "list", "top",
	[]string{utils.StaticJs("common.js"), utils.StaticJs("lib/lazysizes.min.js"), GAMESEXTRAS_URL, utils.StaticJs("list.js")},
	[]string{"top.css"},
)
var HANDLER = handleRequest

var GAMESEXTRAS_URL = "gamesextras"

func handleRequest(request *http.Request, headers http.Header) (string, error) {
	gamesList, _ := datamining.ReadGamesFromStorage()

	gamesList = games.GetTopGames(gamesList, games.RATING_ID_WILSON, 100)

	dataPtr, err := buildTopPageProps(gamesList)
	if err != nil {
		return "", err
	}

	response, err := templates.RenderEndpoint(ENDPOINT, *dataPtr)
	if err != nil {
		return "", err
	}

	return response, nil
}
