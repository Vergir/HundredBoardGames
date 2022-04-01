package games

import (
	"encoding/json"
	"io"
	"os"
)

const storageFilename = "games.json"

func ReadGamesFromStorage() ([]Game, error) {
	file, open_file_error := os.Open(storageFilename)
	if open_file_error != nil {
		return nil, open_file_error
	}

	defer file.Close()

	read_bytes, read_error := io.ReadAll(file)
	if read_error != nil {
		return nil, read_error
	}

	var games []Game

	unmarshal_error := json.Unmarshal(read_bytes, &games)
	if unmarshal_error != nil {
		return nil, unmarshal_error
	}

	return games, nil
}

func writeGamesToStorage(games []Game) error {
	games_json, json_error := json.Marshal(games)
	if json_error != nil {
		return json_error
	}

	write_file_error := os.WriteFile(storageFilename, games_json, 0600)
	if write_file_error != nil {
		return write_file_error
	}

	return nil
}
