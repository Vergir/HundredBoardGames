package games

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"golang.org/x/net/html"
)

var whitespaceRegexp = regexp.MustCompile(`\s`)

func UpdateStorageFromInternet() error {
	const pagesToRead = 100

	var games []Game
	for page_number := 1; page_number <= pagesToRead; page_number++ {
		fmt.Println("Reading page " + fmt.Sprint(page_number))
		resp, err := http.Get("https://boardgamegeek.com/browse/boardgame/page/" + fmt.Sprint(page_number))
		if err != nil {
			return err
		}

		tree, htmlErr := html.Parse(resp.Body)
		if htmlErr != nil {
			return err
		}

		games = append(games, parseGamesFromHtml(tree)...)
	}

	writeError := writeGamesToStorage(games)
	if writeError != nil {
		return writeError
	}

	return nil
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
	geekRatingLabel := whitespaceRegexp.ReplaceAllString(geekRatingCell.FirstChild.Data, "")
	geekRating, _ := strconv.ParseFloat(geekRatingLabel, 64)

	votersRatingCell := geekRatingCell.NextSibling.NextSibling
	votersRatingLabel := whitespaceRegexp.ReplaceAllString(votersRatingCell.FirstChild.Data, "")
	votersRating, _ := strconv.ParseFloat(votersRatingLabel, 64)

	votersCountCell := votersRatingCell.NextSibling.NextSibling
	votersCountLabel := whitespaceRegexp.ReplaceAllString(votersCountCell.FirstChild.Data, "")
	votersCount, _ := strconv.Atoi(votersCountLabel)

	return NewGame(title, year, geekRating, votersRating, uint(votersCount))

}
