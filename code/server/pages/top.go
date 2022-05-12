package pages

import (
	"fmt"
	"html/template"
	"hundred-board-games/code/games"
	"hundred-board-games/code/utils"
	"math"
	"strconv"
)

const COMPLEXITY_SVG_HTML = `<svg class="gameMainAttributeSvg" xmlns="http://www.w3.org/2000/svg" fill="none" stroke="#000" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" aria-labelledby="extensionIconTitle" color="#000" viewBox="0 0 24 24"><path d="M9 4a2 2 0 1 1 4 0v2h5v5h2a2 2 0 1 1 0 4h-2v5h-5v-2a2 2 0 1 0-4 0v2H4v-5h2a2 2 0 1 0 0-4H4V6h5V4Z"/></svg>`
const COMPLEXITY_SVG_PALE_HTML = `<svg class="gameMainAttributeSvg gameMainAttributeSvg--pale" xmlns="http://www.w3.org/2000/svg" fill="none" stroke="#000" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" aria-labelledby="extensionIconTitle" color="#000" viewBox="0 0 24 24"><path d="M9 4a2 2 0 1 1 4 0v2h5v5h2a2 2 0 1 1 0 4h-2v5h-5v-2a2 2 0 1 0-4 0v2H4v-5h2a2 2 0 1 0 0-4H4V6h5V4Z"/></svg>`

var TOP_PAGE = newPage("top", "list", "Топ", "top")

type topPageTemplateProps struct {
	Games []extendedGame
}

type extendedGame struct {
	Rank          uint
	Rating        float32
	RankClass     string
	TitleClass    string
	PictureHbbUrl string
	Playtime      string
	Players       string
	Complexity    template.HTML
	games.Game
}

func PrepareTopPageProps(gamesList []games.Game) topPageTemplateProps {
	games := make([]extendedGame, len(gamesList))

	for i, game := range gamesList {
		rank := i + 1
		roundedRating := math.Round(game.BayesRating*100) / 100
		var complexity string
		switch {
		case game.AvgWeight <= 2:
			complexity = COMPLEXITY_SVG_HTML + COMPLEXITY_SVG_PALE_HTML + COMPLEXITY_SVG_PALE_HTML + " Simple"
		case game.AvgWeight <= 4:
			complexity = COMPLEXITY_SVG_HTML + COMPLEXITY_SVG_HTML + COMPLEXITY_SVG_PALE_HTML + " Moderate"
		default:
			complexity = COMPLEXITY_SVG_HTML + COMPLEXITY_SVG_HTML + COMPLEXITY_SVG_HTML + " Complex"
		}

		rankClass := "gameRank"
		switch {
		case rank > 99:
			rankClass += " gameRank--smaller"
		case rank > 9:
			rankClass += " gameRank--small"
		}

		titleClass := "gameTitle"
		switch {
		case len(game.PrimaryTitle) > 40:
			titleClass += " gameTitle--smaller"
		case len(game.PrimaryTitle) > 30:
			titleClass += " gameTitle--small"
		}

		pictureHbbFilename := utils.FormFullFilename(int(game.GeekId), game.PictureUrl)
		pictureHbbUrl := fmt.Sprint("/static/images/covers/", pictureHbbFilename)

		playtime := strconv.FormatUint(uint64(game.MinPlaytime), 10)
		if game.MinPlaytime != game.MaxPlaytime {
			playtime += "-" + strconv.FormatUint(uint64(game.MaxPlaytime), 10)
		}

		players := strconv.FormatUint(uint64(game.MinPlayers), 10)
		if game.MinPlayers != game.MaxPlayers {
			players += "-" + strconv.FormatUint(uint64(game.MaxPlayers), 10)
		}

		extendedGame := extendedGame{
			uint(rank),
			float32(roundedRating),
			rankClass,
			titleClass,
			pictureHbbUrl,
			playtime,
			players,
			template.HTML(complexity),
			game,
		}
		games[i] = extendedGame
	}

	props := topPageTemplateProps{
		Games: games,
	}

	return props
}
