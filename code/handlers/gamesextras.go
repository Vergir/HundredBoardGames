package handlers

import (
	"encoding/json"
	"html/template"
	"hundred-board-games/code/datamining"
	"hundred-board-games/code/games"
	"hundred-board-games/code/pages"
	"hundred-board-games/code/templates"
	"hundred-board-games/code/utils"
	"net/http"
	"strings"
)

type gameExtra struct {
	Description  string   `json:"d"`
	PicturesUrls []string `json:"p"`
}

type gamesExtras map[uint]gameExtra

type templateProps struct {
	GamesExtrasJson template.HTML //not really html but is the only way to include unescaped json
}

func HandleGamesExtrasQuery(r *http.Request, headers http.Header) (string, error) {
	base64gamesIds := r.URL.Query().Get("games")
	gamesIds, err := pages.Base64ToGamesIds(base64gamesIds)
	if err != nil {
		return "", err
	}

	gamesList, _ := datamining.ReadGamesFromStorage()

	ge := make(gamesExtras)

	for _, game := range gamesList {
		for _, gameId := range gamesIds {
			if gameId != game.GeekId {
				continue
			}

			filesNames, err := utils.ListFolderFiles(games.FormGameImagesPath(game.GeekId))
			if err != nil {
				return "", err
			}

			ge[gameId] = gameExtra{
				Description:  strings.Split(game.Description, "&#10;")[0],
				PicturesUrls: filesNames,
			}
		}
	}

	responseBytes, err := json.Marshal(ge)
	if err != nil {
		return "", err
	}

	response, err := templates.RenderCustom("gamesExtras", templateProps{GamesExtrasJson: template.HTML(responseBytes)})
	if err != nil {
		return "", err
	}

	headers.Set("Content-Type", "text/javascript")

	return response, nil
}
