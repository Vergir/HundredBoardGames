package pages

import (
	"hundred-board-games/games"
)

var LIST_PAGE = newPage("list", "top100", "list")

type listPageData struct {
	PageTitle string
	Games     []games.Game
}

func PrepareListPageData(gamesList []games.Game) any {
	data := listPageData{
		PageTitle: "TO GAMES LIST",
		Games:     gamesList,
	}

	return data
}
