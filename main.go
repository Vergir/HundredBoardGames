package main

import (
	"fmt"
	"hundred-board-games/games"
	"hundred-board-games/server"
	"hundred-board-games/server/pages"
	"net/http"
)

func handleListPage(r *http.Request) string {
	gamesList, _ := games.ReadGamesFromStorage()
	gamesList = games.GetTopGames(gamesList, games.RATING_ID_WILSON, 250)
	data := pages.PrepareTopPageProps(gamesList)

	response, err := server.GetAndRenderTemplate(pages.TOP_PAGE, data)
	if err != nil {
		fmt.Println(err)
	}

	return response
}

func handleIndexPage(r *http.Request) string {
	indexPageProps := pages.PrepareIndexPageProps()

	response, err := server.GetAndRenderTemplate(pages.INDEX_PAGE, indexPageProps)
	if err != nil {
		fmt.Println(err)
	}

	return response
}

func main() {
	server.AddHandler(pages.INDEX_PAGE.Url, handleIndexPage)
	server.AddHandler(pages.TOP_PAGE.Url, handleListPage)

	server.AddStatic()

	server.Start()
}
