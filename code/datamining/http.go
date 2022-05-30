package datamining

import (
	"fmt"
	"hundred-board-games/code/games"
	"hundred-board-games/code/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UpdateStorageFromInternet() error {
	const pagesToRead = 10

	gamesIds, err := queryGamesIds(pagesToRead)
	if err != nil {
		return err
	}

	for _, gameId := range gamesIds {
		fmt.Print("Start game #", gameId, ". ")
		err = downloadGameData(gameId)
		if err != nil {
			return err
		}
		fmt.Println("Finish game #", gameId)
	}

	return nil
}

func queryGamesIds(pagesToRead uint) ([]uint, error) {
	const gamesOnPage = 100

	gamesIds := make([]uint, 0, pagesToRead*gamesOnPage)
	for pageNumber := 1; pageNumber <= int(pagesToRead); pageNumber++ {
		fmt.Printf("Reading page %v\n ", pageNumber)

		url := fmt.Sprintf("https://boardgamegeek.com/browse/boardgame/page/%v", pageNumber)
		gameListReponse, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		pageGamesIds, err := findGamesIds(gameListReponse.Body)
		if err != nil {
			return nil, err
		}

		gamesIds = append(gamesIds, pageGamesIds...)
	}

	return gamesIds, nil
}

func downloadGameData(gameId uint) error {
	game, err := queryGameFields(gameId)
	if err != nil {
		return nil
	}

	err = downloadGameCover(game.PictureUrl, gameId)
	if err != nil {
		return nil
	}

	err = downloadGameImages(gameId)
	if err != nil {
		return nil
	}

	err = saveGameAsJson(*game)
	if err != nil {
		return err
	}

	return nil
}

func queryGameFields(gameId uint) (*games.Game, error) {
	url := fmt.Sprintf("https://api.geekdo.com/xmlapi2/thing?id=%v&stats=1", gameId)
	gameDataResponse, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	game, err := parseGame(gameDataResponse.Body)
	if err != nil {
		return nil, err
	}

	game.UpdateAlgoRatings()

	return game, nil
}

func downloadGameCover(pictureUrl string, gameId uint) error {
	fileName := utils.FormFullFilename(int(gameId), pictureUrl)
	filePath := filepath.Join("static", "images", "covers", fileName)

	if _, err := os.Stat(filePath); err == nil {
		return nil //already have picture
	}

	response, err := http.Get(pictureUrl)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
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

func downloadGameImages(gameBggId uint) error {
	url := fmt.Sprintf("https://api.geekdo.com/api/images?gallery=game&pageid=1&showcount=50&size=thumb&sort=hot&objectid=%v", gameBggId)
	gameDataResponse, err := http.Get(url)
	if err != nil {
		return err
	}

	imagesUrls, err := parseImagesUrlsJson(gameDataResponse.Body)
	if err != nil {
		return err
	}

	for imageIndex, imageUrl := range imagesUrls {
		imageFileResponse, err := http.Get(imageUrl)
		if err != nil {
			return nil
		}

		bytesBuffer, err := utils.ReaderToBytes(imageFileResponse.Body)
		if err != nil {
			return nil
		}

		filename := utils.FormFullFilename(imageIndex, imageUrl)
		imagePath := filepath.Join(games.FormGameImagesPath(gameBggId), filename)

		err = os.MkdirAll(filepath.Dir(imagePath), os.ModePerm)
		if err != nil {
			return err
		}
		err = os.WriteFile(imagePath, bytesBuffer, 0600)
		if err != nil {
			return err
		}
	}

	return nil
}
