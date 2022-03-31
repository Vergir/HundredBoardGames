package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"golang.org/x/net/html"
)

type Game struct {
	Title         string
	Year          int
	Geek_rating   float64
	Voters_rating float64
	Voters_count  int
	algo_rating   float64
}

var whitespace_regexp = regexp.MustCompile(`\s`)

const file_name = "games.json"

func newGame(title string, year int, geek_rating float64, voters_rating float64, voters_count int) Game {
	game := Game{title, year, geek_rating, voters_rating, voters_count, 0}

	return game
}

func parseGameNode(game_node *html.Node) Game {
	firstCell := game_node.FirstChild.NextSibling

	titleCell := firstCell.NextSibling.NextSibling.NextSibling.NextSibling
	titleYearLabel := titleCell.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild.NextSibling
	yearNode := titleYearLabel.NextSibling.NextSibling
	title := titleYearLabel.FirstChild.Data
	year := 0
	if yearNode != nil {
		yearNodeText := yearNode.FirstChild.Data
		year, _ = strconv.Atoi(yearNodeText[1 : len(yearNodeText)-1])
	}

	geekRatingCell := titleCell.NextSibling.NextSibling
	geekRatingLabel := whitespace_regexp.ReplaceAllString(geekRatingCell.FirstChild.Data, "")
	geek_rating, _ := strconv.ParseFloat(geekRatingLabel, 64)

	votersRatingCell := geekRatingCell.NextSibling.NextSibling
	votersRatingLabel := whitespace_regexp.ReplaceAllString(votersRatingCell.FirstChild.Data, "")
	voters_rating, _ := strconv.ParseFloat(votersRatingLabel, 64)

	votersCountCell := votersRatingCell.NextSibling.NextSibling
	votersCountLabel := whitespace_regexp.ReplaceAllString(votersCountCell.FirstChild.Data, "")
	voters_count, _ := strconv.Atoi(votersCountLabel)

	return newGame(title, year, geek_rating, voters_rating, voters_count)

}

func parseGamesFromHtml(node *html.Node) []Game {
	var games []Game

	isGameRowNode := node.Type == html.ElementNode && node.Data == "tr" && len(node.Attr) > 0
	if isGameRowNode {
		game := parseGameNode(node)
		games = append(games, game)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		games = append(games, parseGamesFromHtml(child)...)
	}

	return games
}

func write_games(games []Game) error {
	games_json, json_error := json.Marshal(games)
	if json_error != nil {
		return json_error
	}

	write_file_error := os.WriteFile(file_name, games_json, 0600)
	if write_file_error != nil {
		return write_file_error
	}

	return nil
}

func read_games() ([]Game, error) {
	file, open_file_error := os.Open(file_name)
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

func download_games_data(pages_to_read int) error {
	var games []Game
	for page_number := 1; page_number <= pages_to_read; page_number++ {
		fmt.Println("Reading page " + fmt.Sprint(page_number))
		resp, err := http.Get("https://boardgamegeek.com/browse/boardgame/page/" + fmt.Sprint(page_number))
		if err != nil {
			return err
		}

		tree, html_err := html.Parse(resp.Body)
		if html_err != nil {
			return err
		}

		games = append(games, parseGamesFromHtml(tree)...)
	}

	write_error := write_games(games)
	if write_error != nil {
		return write_error
	}

	return nil
}

func main() {
	games, error := read_games()
	if error != nil {
		log.Fatalln(error)
	}

	fmt.Println(len(games))

}
