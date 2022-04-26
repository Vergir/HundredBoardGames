package games

import (
	"fmt"
	"net/http"
	"time"
)

func UpdateStorageFromInternet() error {
	const pagesToRead = 10

	var games []Game
	for pageNumber := 1; pageNumber <= pagesToRead; pageNumber++ {
		fmt.Printf("Reading page %v\n ", pageNumber)

		url := fmt.Sprintf("https://boardgamegeek.com/browse/boardgame/page/%v", pageNumber)
		gameListReponse, err := http.Get(url)
		if err != nil {
			return err
		}

		gamesIds, err := findGamesIds(gameListReponse.Body)
		if err != nil {
			return err
		}

		for _, gameId := range gamesIds {
			time.Sleep(10 * time.Second)
			url := fmt.Sprintf("https://api.geekdo.com/xmlapi2/thing?id=%v&stats=1", gameId)
			gameDataResponse, err := http.Get(url)
			if err != nil {
				return err
			}

			game, err := parseGame(gameDataResponse.Body)
			if err != nil {
				return err
			}

			game.UpdateAlgoRatings()

			games = append(games, *game)
		}
	}

	writeError := writeGamesToStorage(games)
	if writeError != nil {
		return writeError
	}

	return nil
}
