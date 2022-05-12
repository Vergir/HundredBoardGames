package datamining

import (
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func findGamesIds(reader io.Reader) ([]uint, error) {
	htmlNode, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	gamesIds, err := findGamesIdsInHtml(htmlNode)
	if err != nil {
		return nil, err
	}

	return gamesIds, nil
}

func findGamesIdsInHtml(node *html.Node) ([]uint, error) {
	var gamesIds []uint

	isGameRowNode := node.Type == html.ElementNode && node.Data == "tr" && len(node.Attr) > 0
	if isGameRowNode {
		gameId, err := parseGameNode(node)
		if err != nil {
			return nil, err
		}
		gamesIds = append(gamesIds, gameId)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		foundGamesIds, err := findGamesIdsInHtml(child)
		if err != nil {
			return nil, err
		}
		gamesIds = append(gamesIds, foundGamesIds...)
	}

	return gamesIds, nil
}

func parseGameNode(gameNode *html.Node) (uint, error) {
	firstCell := gameNode.FirstChild.NextSibling

	pictureCell := firstCell.NextSibling.NextSibling

	//follows format "/boardgame/{id}/{handle} i.e /boardgame/174430/gloomhaven"
	bggIdHandleToken := pictureCell.FirstChild.NextSibling.Attr[0].Val
	bggIdHandlePieces := strings.Split(bggIdHandleToken, "/")

	bggId, err := strconv.ParseUint(bggIdHandlePieces[2], 10, 0)
	if err != nil {
		return 0, err
	}

	return uint(bggId), nil
}
