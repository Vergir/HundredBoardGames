package main

import (
	"fmt"
	"hundred-board-games/code/datamining"
	"hundred-board-games/code/games"
	"hundred-board-games/code/handlers"
	"hundred-board-games/code/pages"
	"hundred-board-games/code/server"
	"hundred-board-games/code/server/paths"
	"hundred-board-games/code/templates"
	"net/http"
)

func handleListPage(r *http.Request, headers http.Header) (string, error) {
	gamesList, _ := datamining.ReadGamesFromStorage()
	gamesList = games.GetTopGames(gamesList, games.RATING_ID_WILSON, 100)
	dataPtr, err := pages.PrepareTopPageProps(gamesList)
	if err != nil {
		return "", err
	}

	response, err := templates.RenderPage(pages.TOP_PAGE, *dataPtr)
	if err != nil {
		return "", err
	}

	return response, nil
}

func handleIndexPage(r *http.Request, headers http.Header) (string, error) {
	indexPageProps := pages.PrepareIndexPageProps()

	response, err := templates.RenderPage(pages.INDEX_PAGE, indexPageProps)
	if err != nil {
		fmt.Println(err)
	}

	return response, nil
}

func main() {
	server.AddHandler(paths.PAGE_INDEX, handleIndexPage)
	server.AddHandler(paths.PAGE_TOP, handleListPage)

	server.AddHandler(paths.REQUEST_GAMES_EXTRAS, handlers.HandleGamesExtrasQuery)

	server.AddStatic()

	server.Start()
}
