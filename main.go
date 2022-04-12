package main

import (
	"hundred-board-games/games"
	"hundred-board-games/server"
	"hundred-board-games/server/pages"
	"net/http"
)

func handleListPage(r *http.Request) string {

	gamesList, _ := games.ReadGamesFromStorage()

	response, _ := pages.RenderListPage(gamesList)

	return response
}

func handleIndexPage(r *http.Request) string {
	response, _ := pages.RenderIndexPage()

	return response
}

func main() {
	server.AddHandler(server.PATH_INDEX, handleIndexPage)
	server.AddHandler(server.PATH_LIST, handleListPage)

	server.Start()
}
