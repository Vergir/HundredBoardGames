package pages

import (
	"html/template"
	"hundred-board-games/games"
	"strings"
)

type listPageData struct {
	PageTitle string
	Games     []games.Game
}

func RenderListPage(gamesList []games.Game) (string, error) {
	tmpl, err := template.ParseFiles("templates/list.tmpl")
	if err != nil {
		return "", err
	}

	data := listPageData{
		PageTitle: "TOP GAMES LIST",
		Games:     gamesList,
	}

	var stringBuilder strings.Builder

	err = tmpl.Execute(&stringBuilder, data)
	if err != nil {
		return "", err
	}

	return stringBuilder.String(), nil
}
