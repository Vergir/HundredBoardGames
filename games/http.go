package games

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

		for gameNo, gameId := range gamesIds {
			time.Sleep(5 * time.Second)
			url := fmt.Sprintf("https://api.geekdo.com/xmlapi2/thing?id=%v&stats=1", gameId)
			gameDataResponse, err := http.Get(url)
			if err != nil {
				return err
			}

			game, err := parseGame(gameDataResponse.Body)
			if err != nil {
				return err
			}

			err = downloadGamePicture(game.PictureUrl, game.GeekId)
			if err != nil {
				return err
			}

			game.UpdateAlgoRatings()

			games = append(games, *game)
			fmt.Printf("Done: Page %v, Game %v (%v | %v)\n", pageNumber, gameNo, game.GeekId, game.PrimaryTitle)
		}
	}

	writeError := writeGamesToStorage(games)
	if writeError != nil {
		return writeError
	}

	return nil
}

func downloadGamePicture(pictureUrl string, gameId uint) error {
	filepath := fmt.Sprintf("static/images/covers/%v.jpg", gameId)

	if _, err := os.Stat(filepath); err == nil {
		return nil //already have picture
	}

	response, err := http.Get(pictureUrl)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
