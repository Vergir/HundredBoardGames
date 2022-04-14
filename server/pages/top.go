package pages

import (
	"hundred-board-games/games"
)

var TOP_PAGE = newPage("top", "list", "Топ", "top")

type topPageTemplateProps struct {
	PageTitle string
	Games     []games.Game
}

func (props *topPageTemplateProps) SetPageTitle(title string) {
	props.PageTitle = title
}

func (props *topPageTemplateProps) GetFinalTemplateProps() any {
	return *props
}

func PrepareTopPageProps(gamesList []games.Game) PageProps {
	props := topPageTemplateProps{
		Games: gamesList,
	}

	return &props
}
