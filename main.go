package main

import (
	"fmt"
	"hundred-board-games/code/datamining"
	"hundred-board-games/code/games"
	"hundred-board-games/code/server"
	"hundred-board-games/code/server/pages"
	"net/http"
	"os"
)

func handleListPage(r *http.Request) string {
	gamesList, _ := datamining.ReadGamesFromStorage()
	gamesList = games.GetTopGames(gamesList, games.RATING_ID_WILSON, 100)
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
	err := datamining.UpdateStorageFromInternet()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("FINISH")
	os.Exit(0)

	server.AddHandler(pages.INDEX_PAGE.Url, handleIndexPage)
	server.AddHandler(pages.TOP_PAGE.Url, handleListPage)

	server.AddStatic()

	server.Start()
}
