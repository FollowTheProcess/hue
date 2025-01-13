// Package hue is a simple, modern colour/style package for CLI applications in Go.
package hue

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"

	"golang.org/x/term"
)

const (
	escape = "\x1b["       // escape is the ANSI escape start sequence.
	reset  = escape + "0m" // reset is the universal style reset sequence.
)

// numStyles is a constant representing the expected number of styles that would cover most usage, i.e. only
// pathological combination of styles would exceed this number. This is a useful optimisation as to print
// the escape codes for a combination style we must append to a slice, this number sets the capacity for this
// slice to prevent reallocations.
//
// It is set to 6 to cover what I expect to be most common usage:
//   - A foreground colour
//   - A background colour
//   - Modifiers like bold, italic, underline, strikethrough etc.
//
// It is highly unlikely that people will use more than 6 styles in combination, and if they do, a tiny
// performance penalty is the only downside vs the normal case of < 6.
const numStyles = 6

// disabled is whether this package is turned off and output should be uncoloured.
//
// It defaults to automatic detection.
var disabled atomic.Bool

func init() {
	disabled.Store(autoDisabled())
}

// Enable defines whether the output from this package is colourised. It defaults to automatic
// detection based on a number of attributes:
//   - The value of $NO_COLOR and/or $FORCE_COLOR
//   - The value of $TERM (xterm enables colour)
//   - Whether [os.Stdout] is pointing to a terminal
//
// This means that hue should do a reasonable job of auto-detecting when to colourise output
// and should not write escape sequences when piping between processes or when writing to files etc.
//
// This function may be called to bypass the above detection and force colourisation in all cases. It
// may be called safely from concurrent goroutines.
func Enable() {
	disabled.Store(false)
}

// Disable defines whether the output from this package is colourised. It defaults to automatic
// detection based on a number of attributes:
//   - The value of $NO_COLOR and/or $FORCE_COLOR
//   - The value of $TERM (xterm enables colour)
//   - Whether [os.Stdout] is pointing to a terminal
//
// This means that hue should do a reasonable job of auto-detecting when to colourise output
// and should not write escape sequences when piping between processes or when writing to files etc.
//
// This function may be called to bypass the above detection and disable colourisation in all cases. It
// may be called safely from concurrent goroutines.
func Disable() {
	disabled.Store(true)
}

// Style is a terminal style. It can be a mix of colours and other attributes
// describing the entire appearance of a piece of text.
type Style uint

const (
	Bold                    Style = 1 << iota // Set bold text mode, some terminals may use bright colour variants instead of bold
	Dim                                       // Set dim/faint text mode, not all terminals support this mode
	Italic                                    // Set italic text mode
	Underline                                 // Set underline text mode
	BlinkSlow                                 // Set blink mode with a slow repeat rate
	BlinkFast                                 // Set blink mode with a fast repeat rate
	Reverse                                   // Set reverse/inverse mode, this swaps foreground and background style configuration
	Hidden                                    // Set hidden mode, this hides all text
	Strikethrough                             // Set strikethrough mode
	Black                                     // Black foreground text
	Red                                       // Red foreground text
	Green                                     // Green foreground text
	Yellow                                    // Yellow foreground text
	Blue                                      // Blue foreground text
	Magenta                                   // Magenta foreground text
	Cyan                                      // Cyan foreground text
	White                                     // White foreground text
	BlackBackground                           // Black background
	RedBackground                             // Red background
	GreenBackground                           // Green background
	YellowBackground                          // Yellow background
	BlueBackground                            // Blue background
	MagentaBackground                         // Magenta background
	CyanBackground                            // Cyan background
	WhiteBackground                           // White background
	BrightBlack                               // Bright (high intensity) black foreground text, this means grey on most terminals
	BrightRed                                 // Bright (high intensity) red foreground text
	BrightGreen                               // Bright (high intensity) green foreground text
	BrightYellow                              // Bright (high intensity) yellow foreground text
	BrightBlue                                // Bright (high intensity) blue foreground text
	BrightMagenta                             // Bright (high intensity) magenta foreground text
	BrightCyan                                // Bright (high intensity) cyan foreground text
	BrightWhite                               // Bright (high intensity) white foreground text
	BrightBlackBackground                     // Bright (high intensity) black background, this means grey on most terminals
	BrightRedBackground                       // Bright (high intensity) red background
	BrightGreenBackground                     // Bright (high intensity) green background
	BrightYellowBackground                    // Bright (high intensity) yellow background
	BrightBlueBackground                      // Bright (high intensity) blue background
	BrightMagentaBackground                   // Bright (high intensity) magenta background
	BrightCyanBackground                      // Bright (high intensity) cyan background
	BrightWhiteBackground                     // Bright (high intensity) white background

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

// Code returns the ANSI escape code for the given style, minus the escape
// characters '\x1b[' and 'm' which mark the start and end of the ANSI sequence; respectively.
//
// Callers rarely need this code and should use one of the print style methods instead
// but it is occasionally useful for debugging.
//
// Code returns an error if the style is invalid.
func (s Style) Code() (string, error) {
	if s >= maxStyle || s == 0 {
		return "", fmt.Errorf("invalid style: Style(%d)", s)
	}
	if str, ok := styleStrings[s]; ok {
		return str, nil
	}

	styles := make([]string, 0, numStyles)
	for style := Bold; style <= BrightWhiteBackground; style <<= 1 {
		// If the given style has this style bit set, add its code to the string
		if s&style != 0 {
			code, err := style.Code()
			if err != nil {
				return "", err
			}
			styles = append(styles, code)
		}
	}

	return strings.Join(styles, ";"), nil
}

// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (s Style) Fprint(w io.Writer, a ...any) (n int, err error) {
	text := s.wrap(fmt.Sprint(a...))
	return fmt.Fprint(w, text)
}

// Fprintf formats according to a format specifier and writes to w. It returns
// the number of bytes written and any write error.
func (s Style) Fprintf(w io.Writer, format string, a ...any) (n int, err error) {
	text := s.wrap(fmt.Sprintf(format, a...))
	return fmt.Fprint(w, text)
}

// Fprintln formats using the default format for its operands and writes to w. Spaces are always
// added between operands and a newline is appended. It returns the number of bytes written
// and any write error encountered.
func (s Style) Fprintln(w io.Writer, a ...any) (n int, err error) {
	// Important to add the newline at the very end so wrap the raw text
	// then do Fprintln
	text := s.wrap(fmt.Sprint(a...))
	return fmt.Fprintln(w, text)
}

// Print formats using the default formats for its operands and writes to [os.Stdout]. Spaces are
// added between operands when neither is a string. It returns the number of bytes written and
// any write error encountered.
func (s Style) Print(a ...any) (n int, err error) {
	return s.Fprint(os.Stdout, a...)
}

// Printf formats according to a format specifier and writes to [os.Stdout]. It returns
// the number of bytes written and any write error encountered.
func (s Style) Printf(format string, a ...any) (n int, err error) {
	return s.Fprintf(os.Stdout, format, a...)
}

// Println formats using the default formats for its operands and writes to [os.Stdout]. Spaces are always
// added between operands and a newline is appended. It returns the number of bytes written
// and any write error encountered.
func (s Style) Println(a ...any) (n int, err error) {
	return s.Fprintln(os.Stdout, a...)
}

// Sprint formats using the default formats for its operands and returns the resulting stylised string. Spaces
// are added between operands when neither is a string.
func (s Style) Sprint(a ...any) string {
	return s.wrap(fmt.Sprint(a...))
}

// Sprintf formats according to a format specifier and returns the resulting stylised string.
func (s Style) Sprintf(format string, a ...any) string {
	return s.wrap(fmt.Sprintf(format, a...))
}

// Sprintln formats using the default formats for its operands and returns the resulting string. Spaces are always
// added between operands and a newline is appended.
func (s Style) Sprintln(a ...any) string {
	// Important to add the newline at the very end so wrap the raw text
	// then do Sprintln
	text := s.wrap(fmt.Sprint(a...))
	return fmt.Sprintln(text)
}

// wrap wraps text with the styles escape and reset sequences.
func (s Style) wrap(text string) string {
	if disabled.Load() {
		return text
	}

	code, err := s.Code()
	if err != nil {
		return text
	}
	return escape + code + "m" + text + reset
}

// autoDisabled performs checks to auto detect whether or not this package should be disabled
// based on it's environment.
func autoDisabled() bool {
	// Note: did some digging to see how to avoid potentially 3 different syscalls to get env vars
	// went down a bit of a rabbit hole. It turns out that under the hood, os.Getenv is guarded by a sync.Once
	// so only on the first call to Getenv are we actually making a syscall, all future calls just use the
	// cached copy so no need to do anything clever in user code!

	// $FORCE_COLOR overrides everything
	if os.Getenv("FORCE_COLOR") != "" {
		return false // as in, should not be disabled
	}

	// $NO_COLOR is next
	if os.Getenv("NO_COLOR") != "" {
		return true // no color means we should be disabled
	}

	// If the $TERM env var looks like xtermXXX then it's
	// probably safe e.g. xterm-256-color, xterm-ghostty etc.
	if strings.HasPrefix(os.Getenv("TERM"), "xterm") {
		return false
	}

	// Finally check if stdout's file descriptor is a terminal (best effort)
	if term.IsTerminal(int(os.Stdout.Fd())) {
		return false
	}

	// Can't detect otherwise so be safe and disable colour
	return true
}
