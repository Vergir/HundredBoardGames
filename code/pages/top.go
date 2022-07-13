package pages

import (
	"bytes"
	"compress/lzw"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hundred-board-games/code/games"
	"hundred-board-games/code/server/paths"
	"hundred-board-games/code/utils"
	"io"
	"math"
	"strconv"
	"strings"
)

var TOP_PAGE = newPage(
	"top", "list", "top", paths.PAGE_TOP,
	[]string{utils.StaticJs("common.js"), utils.StaticJs("lib/lazysizes.min.js"), paths.REQUEST_GAMES_EXTRAS, utils.StaticJs("list.js")},
	[]string{"top.css"},
)

type topPageTemplateProps struct {
	Games []extendedGame
}

type extendedGame struct {
	Rank              uint
	Rating            float32
	RankClass         string
	TitleClass        string
	CardPictureUrl    string
	Playtime          string
	Players           string
	ComplexityClasses [3]string
	ComplexityLabel   string
	games.Game
}

func PrepareTopPageProps(gamesList []games.Game) (*topPageTemplateProps, error) {
	extGames := make([]extendedGame, len(gamesList))

	gamesIds := make([]uint, len(gamesList))

	for i, game := range gamesList {
		rank := i + 1
		roundedRating := math.Round(game.BayesRating*100) / 100
		gamesIds[i] = game.GeekId
		complexityClasses := [3]string{}
		complexityLabel := "Complex"
		if game.AvgWeight <= 4 {
			complexityClasses[2] = "gameMainAttributeSvg--pale"
			complexityLabel = "Moderate"
		}
		if game.AvgWeight <= 2 {
			complexityClasses[1] = "gameMainAttributeSvg--pale"
			complexityLabel = "Simple"
		}

		rankClass := "gameRank"
		switch {
		case rank > 99:
			rankClass += " gameRank--smaller"
		case rank > 9:
			rankClass += " gameRank--small"
		}

		titleClass := "gameTitle"
		const titleFontsize = 24.6
		titleLength := computeStringPixelLength(game.PrimaryTitle, titleFontsize)
		switch {
		case titleLength >= 290.0:
			titleClass += " gameTitle--small75"
		case titleLength >= 262.0:
			titleClass += " gameTitle--small80"
		case titleLength >= 240.0:
			titleClass += " gameTitle--small85"
		case titleLength >= 224.0:
			titleClass += " gameTitle--small90"
		case titleLength >= 200.0:
			titleClass += " gameTitle--small95"
		}

		cardPictureFilename := utils.FormComplexFilename(fmt.Sprint(game.GeekId), game.PictureUrl)
		cardPictureUrl := fmt.Sprint("/static/images/covers/200/", cardPictureFilename)

		playtime := strconv.FormatUint(uint64(game.MinPlaytime), 10)
		if game.MinPlaytime != game.MaxPlaytime {
			playtime += "-" + strconv.FormatUint(uint64(game.MaxPlaytime), 10)
		}

		players := strconv.FormatUint(uint64(game.MinPlayers), 10)
		if game.MinPlayers != game.MaxPlayers {
			players += "-" + strconv.FormatUint(uint64(game.MaxPlayers), 10)
		}

		extendedGame := extendedGame{
			uint(rank),
			float32(roundedRating),
			rankClass,
			titleClass,
			cardPictureUrl,
			playtime,
			players,
			complexityClasses,
			complexityLabel,
			game,
		}
		extGames[i] = extendedGame
	}

	props := topPageTemplateProps{
		Games: extGames,
	}

	b64gamesIds, err := gamesIdsToBase64(gamesIds)
	if err != nil {
		return nil, err
	}

	for i, jsPath := range TOP_PAGE.JsPaths {
		if strings.HasPrefix(jsPath, paths.REQUEST_GAMES_EXTRAS) {
			TOP_PAGE.JsPaths[i] = paths.REQUEST_GAMES_EXTRAS + "?games=" + b64gamesIds
			break
		}
	}

	return &props, nil
}

func gamesIdsToBase64(gamesIds []uint) (string, error) {
	var lzwBytes bytes.Buffer

	lzwWriter := lzw.NewWriter(&lzwBytes, lzw.MSB, 8)
	bytesBuffer := make([]byte, 8)
	for _, gameId := range gamesIds {
		binary.BigEndian.PutUint64(bytesBuffer, uint64(gameId))
		lzwWriter.Write(bytesBuffer)
	}
	lzwWriter.Close()

	b64gamesIds := base64.URLEncoding.EncodeToString(lzwBytes.Bytes())

	return b64gamesIds, nil
}

func Base64ToGamesIds(base64string string) ([]uint, error) {
	compressedGamesIdsBytes, err := base64.URLEncoding.DecodeString(base64string)
	if err != nil {
		return nil, err
	}

	lzwReader := lzw.NewReader(bytes.NewBuffer(compressedGamesIdsBytes), lzw.MSB, 8)
	decompressedGamesIdsBytes, err := io.ReadAll(lzwReader)
	if err != nil {
		return nil, err
	}
	lzwReader.Close()

	gamesIdsUint64 := make([]uint64, len(decompressedGamesIdsBytes)/8) //8 = length of uint64

	err = binary.Read(bytes.NewReader(decompressedGamesIdsBytes), binary.BigEndian, gamesIdsUint64)
	if err != nil {
		return nil, err
	}

	gameIds := make([]uint, len(gamesIdsUint64))
	for i, gameIdUint64 := range gamesIdsUint64 {
		gameIds[i] = uint(gameIdUint64)
	}

	return gameIds, nil
}

func computeStringPixelLength(s string, fontSize float64) float64 {
	var pixelsSum uint
	for _, stringSymbol := range s {
		symbolPixels := symbolsWidths[stringSymbol]
		if symbolPixels == 0 {
			symbolPixels = symbolsWidths['0'] //when encounter unknown symbol take width of Zero as a middleground
		}
		pixelsSum += symbolPixels
	}

	divisionFactor := 1000.0 / fontSize

	stringPixelLength := float64(pixelsSum) / divisionFactor

	return stringPixelLength
}

/**
Symbol => Pixel width (1000px font size)

This map was computed for Pragati Narrow using following code

symbols := " !#$&'(),-./0123456789:?ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	for _, symbol := range symbols {
		symbolcmdArg := fmt.Sprint("'", string(symbol), "'")
		stdout, err := exec.Command("magick", "-debug", "annotate", "xc:", "-family", "Pragati Narrow", "-pointsize", "1000", "-annotate", "0", symbolcmdArg, "null:").CombinedOutput()
		if err != nil {
			fmt.Println("REEEEEEEEEE")
			fmt.Print(err)
		}
		output := string(stdout)
		width_index := strings.Index(output, "width: ")
		end_index := strings.Index(output[width_index:], ";")
		widthStr := output[width_index+len("width: ") : width_index+end_index]
		width, _ := strconv.Atoi(widthStr)
		width -= 290 //just ImageMagick things
		fmt.Printf("'%v': %v,\n", string(symbol), width)
	}
*/
var symbolsWidths = map[rune]uint{
	' ':  210,
	'!':  210,
	'#':  419,
	'$':  419,
	'&':  503,
	'\'': 145,
	'(':  251,
	')':  251,
	',':  210,
	'-':  251,
	'.':  210,
	'/':  210,
	'0':  419,
	'1':  419,
	'2':  419,
	'3':  419,
	'4':  419,
	'5':  419,
	'6':  419,
	'7':  419,
	'8':  419,
	'9':  419,
	':':  210,
	'?':  419,
	'A':  503,
	'B':  503,
	'C':  544,
	'D':  544,
	'E':  503,
	'F':  461,
	'G':  587,
	'H':  544,
	'I':  210,
	'J':  376,
	'K':  503,
	'L':  419,
	'M':  628,
	'N':  544,
	'O':  587,
	'P':  503,
	'Q':  587,
	'R':  544,
	'S':  503,
	'T':  461,
	'U':  544,
	'V':  503,
	'W':  711,
	'X':  503,
	'Y':  503,
	'Z':  461,
	'a':  419,
	'b':  419,
	'c':  376,
	'd':  419,
	'e':  419,
	'f':  210,
	'g':  419,
	'h':  419,
	'i':  167,
	'j':  167,
	'k':  376,
	'l':  167,
	'm':  628,
	'n':  419,
	'o':  419,
	'p':  419,
	'q':  419,
	'r':  251,
	's':  376,
	't':  210,
	'u':  419,
	'v':  376,
	'w':  544,
	'x':  376,
	'y':  376,
	'z':  376,
}
