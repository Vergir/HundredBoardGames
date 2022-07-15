package top

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
