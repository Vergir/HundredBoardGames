package datamining

import (
	"bytes"
	"errors"
	"fmt"
	"hundred-board-games/code/games"
	"hundred-board-games/code/utils"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UpdateStorageFromInternet(pagesToRead uint) error {
	gamesIds, err := queryGamesIds(pagesToRead)
	if err != nil {
		return err
	}

	for gameIndex, gameId := range gamesIds {
		fmt.Printf("Start game %v. #%v. ", gameIndex, gameId)
		err = downloadGameData(gameId)
		if err != nil {
			return err
		}
		fmt.Printf("Finish game #%v\n", gameId)
		time.Sleep(15 * time.Second)
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
		return err
	}

	err = saveGameAsJson(*game)
	if err != nil {
		return err
	}

	return nil

	err = downloadGameCover(game.PictureUrl, gameId)
	if err != nil {
		return err
	}

	err = downloadGameImages(gameId)
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
	fileName := utils.FormFullFilename(fmt.Sprint(gameId), pictureUrl)
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
	urlsMap := make(map[string]string, 2)
	urlsMap["a"] = fmt.Sprintf("https://api.geekdo.com/api/images?gallery=game&pageid=1&showcount=50&size=thumb&sort=hot&tag=play&objectid=%v", gameBggId)
	urlsMap["b"] = fmt.Sprintf("https://api.geekdo.com/api/images?gallery=game&pageid=1&showcount=50&size=thumb&sort=hot&objectid=%v", gameBggId)

	for prefix, url := range urlsMap {
		err := downloadGameImagesByUrl(gameBggId, url, prefix)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadGameImagesByUrl(gameBggId uint, url string, imageFilenamePrefix string) error {
	imagesFolderPath := games.FormGameImagesPath(gameBggId)
	err := os.MkdirAll(imagesFolderPath, os.ModePerm)
	if err != nil {
		return err
	}

	folderFilenames, err := utils.ListFolderFiles(imagesFolderPath)
	if err != nil {
		return err
	}

	gameDataResponse, err := http.Get(url)
	if err != nil {
		return err
	}

	imagesUrls, err := parseImagesUrlsJson(gameDataResponse.Body)
	if err != nil {
		return err
	}

	for _, imageUrl := range imagesUrls {
		imageId, imageExtension, err := parseImageUrl(imageUrl)
		if err != nil {
			return err
		}

		if utils.AnyStringHasSubstring(folderFilenames, imageId) {
			continue
		}

		imageFileResponse, err := http.Get(imageUrl)
		if err != nil {
			return err
		}

		imageBytes, err := utils.ReaderToBytes(imageFileResponse.Body)
		if err != nil {
			return err
		}

		imageFilenamePrefix, err = updateImagePrefix(imageBytes, imageExtension, imageFilenamePrefix)
		if err != nil {
			return err
		}

		filename := utils.FormFullFilename(imageFilenamePrefix+imageId, imageUrl)
		imagePath := filepath.Join(imagesFolderPath, filename)

		err = os.WriteFile(imagePath, imageBytes, 0600)
		if err != nil {
			return err
		}
	}

	return nil
}

// expects url to end like ".../pic{id}.{ext}"
func parseImageUrl(url string) (imageId string, imageExtension string, err error) {
	idStartIndex := strings.LastIndex(url, "pic")
	idEndIndex := strings.LastIndexByte(url, '.')

	if idStartIndex == -1 || idEndIndex == -1 {
		return "", "", errors.New("can't extract game id from url")
	}

	imageId = url[idStartIndex+3 : idEndIndex]

	imageExtension = url[idEndIndex+1:]

	return imageId, imageExtension, nil
}

func updateImagePrefix(imageBytes []byte, imageExtension string, oldPrefix string) (string, error) {
	const prefixLast = "z"

	imageBytesReader := bytes.NewReader(imageBytes)

	var imageConfig image.Config
	var err error
	switch imageExtension {
	case "jpg", "jpeg":
		imageConfig, err = jpeg.DecodeConfig(imageBytesReader)
	case "png":
		imageConfig, err = png.DecodeConfig(imageBytesReader)
	default:
		err = errors.New("can't decode image due to unknown extension")
	}

	if err != nil {
		fmt.Println("Error decoding image. " + err.Error())

		return prefixLast, nil
	}

	newPrefix := oldPrefix

	imageRatio := float64(imageConfig.Width) / float64(imageConfig.Height)
	if imageRatio >= 2.0 || imageRatio <= 0.5 {
		newPrefix = prefixLast
	}

	return newPrefix, nil
}
