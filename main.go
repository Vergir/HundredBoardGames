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
	data := pages.PrepareListPageData(gamesList)

	response, err := server.GetAndRenderTemplate(pages.LIST_PAGE.TemplateName, data)
	if err != nil {
		fmt.Println(err)
	}

	return response
}

func handleIndexPage(r *http.Request) string {
	data := pages.PrepareIndexPageData()

	response, err := server.GetAndRenderTemplate(pages.INDEX_PAGE.TemplateName, data)
	if err != nil {
		fmt.Println(err)
	}

	return response
}

func main() {
	server.AddHandler(pages.INDEX_PAGE.Url, handleIndexPage)
	server.AddHandler(pages.LIST_PAGE.Url, handleListPage)

	server.AddStatic()

	server.Start()
}
