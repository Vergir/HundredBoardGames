package datamining

import (
	"encoding/json"
	"errors"
	"fmt"
	"hundred-board-games/code/games"
	"hundred-board-games/code/utils"
	"io"
	"os"
	"path/filepath"
	"time"
)

type imagesPageData struct {
	Images []struct {
		Imageid      string      `json:"imageid"`
		ImageurlLg   string      `json:"imageurl_lg"`
		Name         interface{} `json:"name"`
		Caption      string      `json:"caption"`
		Numrecommend string      `json:"numrecommend"`
		Numcomments  string      `json:"numcomments"`
		User         struct {
			Username   string `json:"username"`
			Avatar     string `json:"avatar"`
			Avatarfile string `json:"avatarfile"`
		} `json:"user"`
		Imageurl string `json:"imageurl"`
		Href     string `json:"href"`
	} `json:"images"`
	Config struct {
		Sorttypes []struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"sorttypes"`
		Numitems    int `json:"numitems"`
		Endpage     int `json:"endpage"`
		Datefilters []struct {
			Value string `json:"value"`
			Name  string `json:"name"`
		} `json:"datefilters"`
		Licensefilters []struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"licensefilters"`
		Filters []struct {
			Name     string `json:"name"`
			Listname string `json:"listname"`
			Type     string `json:"type"`
		} `json:"filters"`
	} `json:"config"`
	Operations []struct {
		Key     string `json:"key"`
		Label   string `json:"label"`
		Options []struct {
			Value string `json:"value"`
			Label string `json:"label"`
		} `json:"options"`
		Default string `json:"default"`
	} `json:"operations"`
	Pagination struct {
		PerPage int `json:"perPage"`
		Pageid  int `json:"pageid"`
		Total   int `json:"total"`
	} `json:"pagination"`
}

const GAMES_FOLDER_NAME = "games"

func parseImagesUrlsJson(reader io.Reader) ([]string, error) {
	var pd imagesPageData
	readerBytes, err := utils.ReaderToBytes(reader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(readerBytes, &pd)
	if err != nil && errors.Is(err, &json.UnmarshalTypeError{}) {
		return nil, err
	}

	urls := make([]string, len(pd.Images))
	for i, entry := range pd.Images {
		urls[i] = entry.ImageurlLg
	}

	return urls, nil
}

func ReadGamesFromStorage() ([]games.Game, error) {
	filesNames, err := utils.ListFolderFiles(GAMES_FOLDER_NAME)
	if err != nil {
		return nil, err
	}

	filesNamesChan := make(chan string, len(filesNames))
	parsedGamesChan := make(chan games.Game, len(filesNames))

	const workersCount = 8
	for workerIndex := 0; workerIndex < workersCount; workerIndex++ {
		go readGameFromJson(filesNamesChan, parsedGamesChan)
	}

	for _, filename := range filesNames {
		filesNamesChan <- filename
	}
	close(filesNamesChan)

	parsedGames := make([]games.Game, len(filesNames))
	for gamesProcessed := 0; gamesProcessed < len(filesNames); gamesProcessed++ {
		select {
		case parsedGames[gamesProcessed] = <-parsedGamesChan:
		case <-time.After(1 * time.Second):
			return nil, fmt.Errorf("waiting for channels timeout. %v games overall, %v done", len(filesNames), gamesProcessed)
		}
	}

	return parsedGames, nil
}

func readGameFromJson(filesNames <-chan string, parsedGames chan<- games.Game) {
	for fileName, isOpen := <-filesNames; isOpen; fileName, isOpen = <-filesNames {
		if fileName == "" {
			break
		}
		file, err := os.Open(formGameJsonPathByFilename(fileName))
		if err != nil {
			break
		}

		bytes, err := io.ReadAll(file)
		if err != nil {
			break
		}

		var game games.Game

		err = json.Unmarshal(bytes, &game)
		if err != nil {
			break
		}

		err = file.Close()
		if err != nil {
			break
		}

		parsedGames <- game
	}
}

func saveGameAsJson(game games.Game) error {
	game_json, err := json.Marshal(game)
	if err != nil {
		return err
	}

	err = os.WriteFile(formGameJsonPathById(game.GeekId), game_json, 0600)
	if err != nil {
		return err
	}

	return nil
}

func formGameJsonPathByFilename(filename string) string {
	return filepath.Join(GAMES_FOLDER_NAME, filename)
}

func formGameJsonPathById(gameId uint) string {
	return formGameJsonPathByFilename(fmt.Sprint(gameId, ".json"))
}
