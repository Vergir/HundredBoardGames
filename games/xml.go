package games

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type items struct {
	XMLName    xml.Name `xml:"items"`
	Text       string   `xml:",chardata"`
	Termsofuse string   `xml:"termsofuse,attr"`
	Item       struct {
		Text      string `xml:",chardata"`
		Type      string `xml:"type,attr"`
		ID        uint   `xml:"id,attr"`
		Thumbnail string `xml:"thumbnail"`
		Image     string `xml:"image"`
		Name      []struct {
			Text      string `xml:",chardata"`
			Type      string `xml:"type,attr"`
			Sortindex string `xml:"sortindex,attr"`
			Value     string `xml:"value,attr"`
		} `xml:"name"`
		Description   string `xml:"description"`
		Yearpublished struct {
			Text  string `xml:",chardata"`
			Value int16  `xml:"value,attr"`
		} `xml:"yearpublished"`
		Minplayers struct {
			Text  string `xml:",chardata"`
			Value uint   `xml:"value,attr"`
		} `xml:"minplayers"`
		Maxplayers struct {
			Text  string `xml:",chardata"`
			Value uint   `xml:"value,attr"`
		} `xml:"maxplayers"`
		Poll        []poll `xml:"poll"`
		Playingtime struct {
			Text  string `xml:",chardata"`
			Value uint   `xml:"value,attr"`
		} `xml:"playingtime"`
		Minplaytime struct {
			Text  string `xml:",chardata"`
			Value uint   `xml:"value,attr"`
		} `xml:"minplaytime"`
		Maxplaytime struct {
			Text  string `xml:",chardata"`
			Value uint   `xml:"value,attr"`
		} `xml:"maxplaytime"`
		Minage struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"minage"`
		Link []struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			ID    uint   `xml:"id,attr"`
			Value string `xml:"value,attr"`
		} `xml:"link"`
		Statistics struct {
			Text    string `xml:",chardata"`
			Page    uint   `xml:"page,attr"`
			Ratings struct {
				Text       string `xml:",chardata"`
				Usersrated struct {
					Text  string `xml:",chardata"`
					Value uint   `xml:"value,attr"`
				} `xml:"usersrated"`
				Average struct {
					Text  string  `xml:",chardata"`
					Value float64 `xml:"value,attr"`
				} `xml:"average"`
				Bayesaverage struct {
					Text  string  `xml:",chardata"`
					Value float64 `xml:"value,attr"`
				} `xml:"bayesaverage"`
				Ranks struct {
					Text string `xml:",chardata"`
					Rank []struct {
						Text         string  `xml:",chardata"`
						Type         string  `xml:"type,attr"`
						ID           uint    `xml:"id,attr"`
						Name         string  `xml:"name,attr"`
						Friendlyname string  `xml:"friendlyname,attr"`
						Value        uint    `xml:"value,attr"`
						Bayesaverage float64 `xml:"bayesaverage,attr"`
					} `xml:"rank"`
				} `xml:"ranks"`
				Stddev struct {
					Text  string  `xml:",chardata"`
					Value float64 `xml:"value,attr"`
				} `xml:"stddev"`
				Median struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value,attr"`
				} `xml:"median"`
				Owned struct {
					Text  string `xml:",chardata"`
					Value uint   `xml:"value,attr"`
				} `xml:"owned"`
				Trading struct {
					Text  string `xml:",chardata"`
					Value uint   `xml:"value,attr"`
				} `xml:"trading"`
				Wanting struct {
					Text  string `xml:",chardata"`
					Value uint   `xml:"value,attr"`
				} `xml:"wanting"`
				Wishing struct {
					Text  string `xml:",chardata"`
					Value uint   `xml:"value,attr"`
				} `xml:"wishing"`
				Numcomments struct {
					Text  string `xml:",chardata"`
					Value uint   `xml:"value,attr"`
				} `xml:"numcomments"`
				Numweights struct {
					Text  string `xml:",chardata"`
					Value uint   `xml:"value,attr"`
				} `xml:"numweights"`
				Averageweight struct {
					Text  string  `xml:",chardata"`
					Value float64 `xml:"value,attr"`
				} `xml:"averageweight"`
			} `xml:"ratings"`
		} `xml:"statistics"`
	} `xml:"item"`
}

type poll struct {
	Text       string `xml:",chardata"`
	Name       string `xml:"name,attr"`
	Title      string `xml:"title,attr"`
	Totalvotes uint   `xml:"totalvotes,attr"`
	Results    []struct {
		Text       string `xml:",chardata"`
		Numplayers string `xml:"numplayers,attr"`
		Result     []struct {
			Text     string `xml:",chardata"`
			Value    string `xml:"value,attr"`
			Numvotes uint   `xml:"numvotes,attr"`
			Level    uint   `xml:"level,attr"`
		} `xml:"result"`
	} `xml:"results"`
}

func parseGame(reader io.Reader) (*Game, error) {
	streamBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var items items
	err = xml.Unmarshal(streamBytes, &items)
	if err != nil {
		test := string(streamBytes)
		fmt.Println(test)
		return nil, err
	}

	game, err := parseItemsIntoGame(items)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func parseItemsIntoGame(items items) (*Game, error) {
	item := items.Item

	game := Game{
		GeekId:         item.ID,
		Year:           item.Yearpublished.Value,
		Description:    item.Description,
		PictureUrl:     item.Image,
		MinPlayers:     uint8(item.Minplayers.Value),
		MaxPlayers:     uint8(item.Maxplayers.Value),
		MinPlaytime:    item.Minplaytime.Value,
		MaxPlaytime:    item.Maxplaytime.Value,
		AvgRating:      item.Statistics.Ratings.Average.Value,
		BayesRating:    item.Statistics.Ratings.Bayesaverage.Value,
		RatingNumVotes: item.Statistics.Ratings.Usersrated.Value,
		AvgWeight:      item.Statistics.Ratings.Averageweight.Value,
		WeightNumVotes: item.Statistics.Ratings.Numweights.Value,
		Counters: gameCounters{
			Owned:    item.Statistics.Ratings.Owned.Value,
			Trading:  item.Statistics.Ratings.Trading.Value,
			Wanting:  item.Statistics.Ratings.Wanting.Value,
			Wishing:  item.Statistics.Ratings.Wishing.Value,
			Comments: item.Statistics.Ratings.Numcomments.Value,
		},
	}

	ageContainsLetters := len(item.Minage.Value) > 2
	if ageContainsLetters {
		item.Minage.Value = item.Minage.Value[:2]
	}
	minAge, err := strconv.ParseUint(item.Minage.Value, 10, 8)
	if err != nil {
		return nil, err
	}
	game.MinAge = uint8(minAge)

	for _, title := range item.Name {
		if title.Type == "primary" {
			game.PrimaryTitle = title.Value
		} else {
			game.Titles = append(game.Titles, title.Value)
		}
	}

	for _, communityPoll := range item.Poll {
		var err error
		switch communityPoll.Name {
		case "suggested_numplayers":
			err = fillCommunityNumPlayers(&game, communityPoll)
		case "suggested_playerage":
			err = fillCommunityMinAge(&game, communityPoll)
		case "language_dependence":
			fillLanguageDependance(&game, communityPoll)
		}
		if err != nil {
			return nil, err
		}
	}

	for _, link := range item.Link {
		game.Tags = append(game.Tags, tag{Type: link.Type, Id: link.ID})
	}

	return &game, nil
}

func fillCommunityNumPlayers(game *Game, xmlNumPlayersPolls poll) error {
	for _, numPlayersOption := range xmlNumPlayersPolls.Results {
		communityNumPlayersEntry := numPlayersPoll{}

		addOneNumPlayer := false
		if strings.HasSuffix(numPlayersOption.Numplayers, "+") {
			numPlayersOption.Numplayers = numPlayersOption.Numplayers[:len(numPlayersOption.Numplayers)-1]
			addOneNumPlayer = true
		}
		numPlayers, err := strconv.ParseUint(numPlayersOption.Numplayers, 10, 8)
		if err != nil {
			return err
		}
		if addOneNumPlayer {
			numPlayers += 1
		}

		communityNumPlayersEntry.NumPlayers = uint8(numPlayers)

		for _, numPlayersResults := range numPlayersOption.Result {
			switch numPlayersResults.Value {
			case "Best":
				communityNumPlayersEntry.VotedBest = numPlayersResults.Numvotes
			case "Recommended":
				communityNumPlayersEntry.VotedRecommended = numPlayersResults.Numvotes
			case "Not Recommended":
				communityNumPlayersEntry.VotedNotRecommended = numPlayersResults.Numvotes
			}
		}

		game.CommunityNumPlayers = append(game.CommunityNumPlayers, communityNumPlayersEntry)
	}

	return nil
}

func fillCommunityMinAge(game *Game, xmlMinAgePoll poll) error {
	game.CommunityMinAge = make([]minAgePoll, 0)
	if len(xmlMinAgePoll.Results) == 0 {
		return nil
	}
	for _, minAgePollResult := range xmlMinAgePoll.Results[0].Result {
		ageContainsLetters := len(minAgePollResult.Value) > 2
		if ageContainsLetters {
			minAgePollResult.Value = minAgePollResult.Value[:2]
		}

		playerAge, err := strconv.ParseUint(minAgePollResult.Value, 10, 8)
		if err != nil {
			return err
		}

		communityMinAgeEntry := minAgePoll{
			MinAge:   uint8(playerAge),
			NumVotes: minAgePollResult.Numvotes,
		}

		game.CommunityMinAge = append(game.CommunityMinAge, communityMinAgeEntry)
	}

	return nil
}

func fillLanguageDependance(game *Game, xmlLanguageDependancePoll poll) {
	for _, langDepPollResult := range xmlLanguageDependancePoll.Results[0].Result {
		languageDependanceEntry := langDepPoll{
			Level:    uint8(langDepPollResult.Level),
			NumVotes: langDepPollResult.Numvotes,
		}

		game.LanguageDependence = append(game.LanguageDependence, languageDependanceEntry)
	}
}
