package top

import (
	"fmt"
	"hundred-board-games/code/games"
	"hundred-board-games/code/utils"
	"math"
	"strconv"
	"strings"
)

const titleFontSizePx = 24.6

type extendedGame struct {
	Rank              uint
	Rating            float32
	RankClass         string
	TitleClass        string
	CardPictureUrl    string
	Playtime          string
	Players           string
	ComplexityClasses [3]string
	ComplexityLabel   string
	games.Game
}

type topPageTemplateProps struct {
	Games []extendedGame
}

func formComplexity(gameAvgWeight float64) (complexityClasses [3]string, complexityLabel string) {
	complexityClasses = [3]string{}
	complexityLabel = "Complex"

	if gameAvgWeight <= 4 {
		complexityClasses[2] = "gameMainAttributeSvg--pale"
		complexityLabel = "Moderate"
	}

	if gameAvgWeight <= 2 {
		complexityClasses[1] = "gameMainAttributeSvg--pale"
		complexityLabel = "Simple"
	}

	return complexityClasses, complexityLabel
}

func formGameRankClass(rank int) string {
	rankClass := "gameRank"
	switch {
	case rank > 99:
		rankClass += " gameRank--smaller"
	case rank > 9:
		rankClass += " gameRank--small"
	}

	return rankClass
}

func formTitleFontSizeClass(gamePrimaryTitle string) string {
	titleClass := "gameTitle"
	titleLength := computeStringPixelLength(gamePrimaryTitle, titleFontSizePx)
	switch {
	case titleLength >= 290.0:
		titleClass += " gameTitle--small75"
	case titleLength >= 262.0:
		titleClass += " gameTitle--small80"
	case titleLength >= 240.0:
		titleClass += " gameTitle--small85"
	case titleLength >= 224.0:
		titleClass += " gameTitle--small90"
	case titleLength >= 200.0:
		titleClass += " gameTitle--small95"
	}

	return titleClass
}

func formGamePlaytimeLabel(gameMinPlaytime uint, gameMaxPlaytime uint) string {
	playtime := strconv.FormatUint(uint64(gameMinPlaytime), 10)
	if gameMinPlaytime != gameMaxPlaytime {
		playtime += "-" + strconv.FormatUint(uint64(gameMaxPlaytime), 10)
	}

	return playtime
}

func formGamePlayersLabel(gameMinPlayers uint8, gameMaxPlayers uint8) string {
	players := strconv.FormatUint(uint64(gameMinPlayers), 10)
	if gameMinPlayers != gameMaxPlayers {
		players += "-" + strconv.FormatUint(uint64(gameMaxPlayers), 10)
	}

	return players
}

func formExtendedGame(sorted_index int, game games.Game) extendedGame {
	rank := sorted_index + 1
	roundedRating := math.Round(game.BayesRating*100) / 100
	complexityClasses, complexityLabel := formComplexity(game.AvgWeight)
	rankClass := formGameRankClass(rank)
	titleClass := formTitleFontSizeClass(game.PrimaryTitle)
	cardPictureUrl := "/static/images/covers/200/" + utils.FormComplexFilename(fmt.Sprint(game.GeekId), game.PictureUrl)
	playtime := formGamePlaytimeLabel(game.MinPlaytime, game.MaxPlaytime)
	players := formGamePlayersLabel(game.MinPlayers, game.MaxPlayers)

	extendedGame := extendedGame{
		uint(rank),
		float32(roundedRating),
		rankClass,
		titleClass,
		cardPictureUrl,
		playtime,
		players,
		complexityClasses,
		complexityLabel,
		game,
	}

	return extendedGame
}

func formGamesExtrasQuery(endpointJsPaths []string, gamesIds []uint) ([]string, error) {
	b64gamesIds, err := gamesIdsToBase64(gamesIds)
	if err != nil {
		return nil, err
	}

	newJsPaths := endpointJsPaths

	for i, jsPath := range endpointJsPaths {
		if strings.HasPrefix(jsPath, GAMESEXTRAS_URL) {
			newJsPaths[i] = GAMESEXTRAS_URL + "?games=" + b64gamesIds
		}
	}

	return newJsPaths, nil
}

func buildTopPageProps(gamesList []games.Game) (*topPageTemplateProps, error) {
	extGames := make([]extendedGame, len(gamesList))

	gamesIds := make([]uint, len(gamesList))

	for i, game := range gamesList {
		gamesIds[i] = game.GeekId
		extGames[i] = formExtendedGame(i, game)
	}

	var err error
	ENDPOINT.JsPaths, err = formGamesExtrasQuery(ENDPOINT.JsPaths, gamesIds)
	if err != nil {
		return nil, err
	}

	props := topPageTemplateProps{
		Games: extGames,
	}

	return &props, nil
}
