// Package hue is a placeholder for something cool.
package hue

import (
	"fmt"
	"strings"
)

// escape is the ANSI escape character.
// const escape = "\x1b"

// Style is a terminal style. It can be a mix of colours and other attributes
// describing the entire appearance of a piece of text.
type Style uint

const (
	Bold Style = 1 << iota
	Dim
	Italic
	Underline
	BlinkSlow
	BlinkFast
	Reverse
	Hidden
	Strikethrough
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BlackBackground
	RedBackground
	GreenBackground
	YellowBackground
	BlueBackground
	MagentaBackground
	CyanBackground
	WhiteBackground
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
	BrightBlackBackground
	BrightRedBackground
	BrightGreenBackground
	BrightYellowBackground
	BrightBlueBackground
	BrightMagentaBackground
	BrightCyanBackground
	BrightWhiteBackground

	maxStyle
)

// styleStrings is a map of the style to it's escape sequence digit.
var styleStrings = map[Style]string{ //nolint: exhaustive // We don't need maxStyle in here it's just a marker
	Bold:                    "1",
	Dim:                     "2",
	Italic:                  "3",
	Underline:               "4",
	BlinkSlow:               "5",
	BlinkFast:               "6",
	Reverse:                 "7",
	Hidden:                  "8",
	Strikethrough:           "9",
	Black:                   "30",
	Red:                     "31",
	Green:                   "32",
	Yellow:                  "33",
	Blue:                    "34",
	Magenta:                 "35",
	Cyan:                    "36",
	White:                   "37",
	BlackBackground:         "40",
	RedBackground:           "41",
	GreenBackground:         "42",
	YellowBackground:        "43",
	BlueBackground:          "44",
	MagentaBackground:       "45",
	CyanBackground:          "46",
	WhiteBackground:         "47",
	BrightBlack:             "90",
	BrightRed:               "91",
	BrightGreen:             "92",
	BrightYellow:            "93",
	BrightBlue:              "94",
	BrightMagenta:           "95",
	BrightCyan:              "96",
	BrightWhite:             "97",
	BrightBlackBackground:   "100",
	BrightRedBackground:     "101",
	BrightGreenBackground:   "102",
	BrightYellowBackground:  "103",
	BrightBlueBackground:    "104",
	BrightMagentaBackground: "105",
	BrightCyanBackground:    "106",
	BrightWhiteBackground:   "107",
}

// String implements [fmt.Stringer] for Style.
func (s Style) String() string {
	if s >= maxStyle {
		return fmt.Sprintf("invalid style: Style(%d)", s)
	}
	if str, ok := styleStrings[s]; ok {
		return str
	}

	// TODO(@FollowTheProcess): The below case allocates, see if we can eliminate the allocation
	// TODO(@FollowTheProcess): Width padding so that it aligns with text/tabwriter properly

	// Combinations
	var styles []string
	for style := Bold; style <= BrightWhiteBackground; style <<= 1 {
		// If the given style has this style bit set, add it to the string
		if s&style != 0 {
			styles = append(styles, style.String())
		}
	}

	return strings.Join(styles, ";") + "m"
}
