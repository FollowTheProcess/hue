package hue_test

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/FollowTheProcess/hue"
)

func TestFprint(t *testing.T) {
	tests := []struct {
		name    string    // Name of the test case
		input   string    // Text to style
		want    string    // Expected result including escape sequences
		enabled bool      // Whether hue is enabled
		style   hue.Style // Style under test
	}{
		{
			name:    "basic",
			input:   "hello",
			style:   hue.Green,
			enabled: true,
			want:    "\x1b[32mhello\x1b[0m",
		},
		{
			name:    "many styles",
			input:   "hello",
			style:   hue.Green | hue.BlueBackground | hue.Bold | hue.Underline,
			enabled: true,
			want:    "\x1b[1;4;32;44mhello\x1b[0m",
		},
		{
			name:    "basic disabled",
			input:   "hello",
			style:   hue.Green,
			enabled: false,
			want:    "hello",
		},
		{
			name:    "many styles disabled",
			input:   "hello",
			style:   hue.Green | hue.BlueBackground | hue.Bold | hue.Underline,
			enabled: false,
			want:    "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure the behaviour is explicitly as requested
			hue.Enabled(tt.enabled)
			buf := &bytes.Buffer{}
			tt.style.Fprint(buf, tt.input)

			got := strconv.Quote(buf.String())
			want := strconv.Quote(tt.want)

			if got != want {
				t.Errorf("\nGot:\t%v\nWanted:\t%v\n", got, want)
			}
		})
	}
}

func TestFprintf(t *testing.T) {
	tests := []struct {
		name    string    // Name of the test case
		input   string    // Text to style
		want    string    // Expected result including escape sequences
		args    []any     // Args to Fprintf
		enabled bool      // Whether hue is enabled
		style   hue.Style // Style under test
	}{
		{
			name:    "basic",
			input:   "hello %s",
			args:    []any{"hue"},
			style:   hue.Magenta,
			enabled: true,
			want:    "\x1b[35mhello hue\x1b[0m",
		},
		{
			name:    "many styles",
			input:   "how many styles %s? %d",
			args:    []any{"hue", 4},
			style:   hue.Blue | hue.RedBackground | hue.Italic | hue.Bold,
			enabled: true,
			want:    "\x1b[1;3;34;41mhow many styles hue? 4\x1b[0m",
		},
		{
			name:    "basic disabled",
			input:   "hello %s",
			args:    []any{"hue"},
			style:   hue.Magenta,
			enabled: false,
			want:    "hello hue",
		},
		{
			name:    "many styles disabled",
			input:   "how many styles %s? %d",
			args:    []any{"hue", 4},
			style:   hue.Blue | hue.RedBackground | hue.Italic | hue.Bold,
			enabled: false,
			want:    "how many styles hue? 4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure the behaviour is explicitly as requested
			hue.Enabled(tt.enabled)
			buf := &bytes.Buffer{}
			tt.style.Fprintf(buf, tt.input, tt.args...)

			got := strconv.Quote(buf.String())
			want := strconv.Quote(tt.want)

			if got != want {
				t.Errorf("\nGot:\t%v\nWanted:\t%v\n", got, want)
			}
		})
	}
}

func TestFprintln(t *testing.T) {
	tests := []struct {
		name    string    // Name of the test case
		input   string    // Text to style
		want    string    // Expected result including escape sequences
		enabled bool      // Whether hue is enabled
		style   hue.Style // Style under test
	}{
		{
			name:    "basic",
			input:   "woah!",
			style:   hue.BrightGreen,
			enabled: true,
			want:    "\x1b[92mwoah!\x1b[0m\n",
		},
		{
			name:    "many styles",
			input:   "such wow",
			style:   hue.BrightCyan | hue.Strikethrough | hue.Italic,
			enabled: true,
			want:    "\x1b[3;9;96msuch wow\x1b[0m\n",
		},
		{
			name:    "basic disabled",
			input:   "woah!",
			style:   hue.BrightGreen,
			enabled: false,
			want:    "woah!\n",
		},
		{
			name:    "many styles disabled",
			input:   "such wow",
			style:   hue.BrightCyan | hue.Strikethrough | hue.Italic,
			enabled: false,
			want:    "such wow\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure the behaviour is explicitly as requested
			hue.Enabled(tt.enabled)
			buf := &bytes.Buffer{}
			tt.style.Fprintln(buf, tt.input)

			got := strconv.Quote(buf.String())
			want := strconv.Quote(tt.want)

			if got != want {
				t.Errorf("\nGot:\t%v\nWanted:\t%v\n", got, want)
			}
		})
	}
}

func TestPrint(t *testing.T) {
	tests := []struct {
		name    string    // Name of the test case
		input   string    // Text to style
		want    string    // Expected result including escape sequences
		enabled bool      // Whether hue is enabled
		style   hue.Style // Style under test
	}{
		{
			name:    "basic",
			input:   "hello",
			style:   hue.Red,
			enabled: true,
			want:    "\x1b[31mhello\x1b[0m",
		},
		{
			name:    "many styles",
			input:   "hello",
			style:   hue.Yellow | hue.BlackBackground | hue.Bold | hue.Italic,
			enabled: true,
			want:    "\x1b[1;3;33;40mhello\x1b[0m",
		},
		{
			name:    "basic disabled",
			input:   "hello",
			style:   hue.Red,
			enabled: false,
			want:    "hello",
		},
		{
			name:    "many styles disabled",
			input:   "hello",
			style:   hue.Yellow | hue.BlackBackground | hue.Bold | hue.Italic,
			enabled: false,
			want:    "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure the behaviour is explicitly as requested
			hue.Enabled(tt.enabled)

			stdout := captureOutput(t, func() error {
				_, err := tt.style.Print(tt.input)
				return err
			})

			got := strconv.Quote(stdout)
			want := strconv.Quote(tt.want)

			if got != want {
				t.Errorf("\nGot:\t%v\nWanted:\t%v\n", got, want)
			}
		})
	}
}

func TestPrintf(t *testing.T) {
	tests := []struct {
		name    string    // Name of the test case
		input   string    // Text to style
		want    string    // Expected result including escape sequences
		args    []any     // Args to Fprintf
		enabled bool      // Whether hue is enabled
		style   hue.Style // Style under test
	}{
		{
			name:    "basic",
			input:   "hello %s",
			args:    []any{"hue"},
			style:   hue.BrightYellow,
			enabled: true,
			want:    "\x1b[93mhello hue\x1b[0m",
		},
		{
			name:    "many styles",
			input:   "how many styles %s? %d",
			args:    []any{"hue", 4},
			style:   hue.BrightRed | hue.BrightBlackBackground | hue.Underline | hue.Dim,
			enabled: true,
			want:    "\x1b[2;4;91;100mhow many styles hue? 4\x1b[0m",
		},
		{
			name:    "basic disabled",
			input:   "hello %s",
			args:    []any{"hue"},
			style:   hue.BrightYellow,
			enabled: false,
			want:    "hello hue",
		},
		{
			name:    "many styles disabled",
			input:   "how many styles %s? %d",
			args:    []any{"hue", 4},
			style:   hue.BrightRed | hue.BrightBlackBackground | hue.Underline | hue.Dim,
			enabled: false,
			want:    "how many styles hue? 4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure the behaviour is explicitly as requested
			hue.Enabled(tt.enabled)
			stdout := captureOutput(t, func() error {
				_, err := tt.style.Printf(tt.input, tt.args...)
				return err
			})

			got := strconv.Quote(stdout)
			want := strconv.Quote(tt.want)

			if got != want {
				t.Errorf("\nGot:\t%v\nWanted:\t%v\n", got, want)
			}
		})
	}
}

func TestPrintln(t *testing.T) {
	tests := []struct {
		name    string    // Name of the test case
		input   string    // Text to style
		want    string    // Expected result including escape sequences
		enabled bool      // Whether hue is enabled
		style   hue.Style // Style under test
	}{
		{
			name:    "basic",
			input:   "woah!",
			style:   hue.Italic,
			enabled: true,
			want:    "\x1b[3mwoah!\x1b[0m\n",
		},
		{
			name:    "many styles",
			input:   "such wow",
			style:   hue.BrightGreen | hue.Dim | hue.Underline,
			enabled: true,
			want:    "\x1b[2;4;92msuch wow\x1b[0m\n",
		},
		{
			name:    "basic disabled",
			input:   "woah!",
			style:   hue.BrightGreen,
			enabled: false,
			want:    "woah!\n",
		},
		{
			name:    "many styles disabled",
			input:   "such wow",
			style:   hue.BrightCyan | hue.Strikethrough | hue.Bold,
			enabled: false,
			want:    "such wow\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure the behaviour is explicitly as requested
			hue.Enabled(tt.enabled)
			stdout := captureOutput(t, func() error {
				_, err := tt.style.Println(tt.input)
				return err
			})

			got := strconv.Quote(stdout)
			want := strconv.Quote(tt.want)

			if got != want {
				t.Errorf("\nGot:\t%v\nWanted:\t%v\n", got, want)
			}
		})
	}
}

func TestSprint(t *testing.T) {
	tests := []struct {
		name    string    // Name of the test case
		input   string    // Text to style
		want    string    // Expected result including escape sequences
		enabled bool      // Whether hue is enabled
		style   hue.Style // Style under test
	}{
		{
			name:    "basic",
			input:   "woah!",
			style:   hue.Blue,
			enabled: true,
			want:    "\x1b[34mwoah!\x1b[0m",
		},
		{
			name:    "many styles",
			input:   "such wow",
			style:   hue.BrightMagenta | hue.Underline | hue.GreenBackground,
			enabled: true,
			want:    "\x1b[4;42;95msuch wow\x1b[0m",
		},
		{
			name:    "basic disabled",
			input:   "woah!",
			style:   hue.BrightWhite,
			enabled: false,
			want:    "woah!",
		},
		{
			name:    "many styles",
			input:   "such wow",
			style:   hue.BrightMagenta | hue.Underline | hue.GreenBackground,
			enabled: false,
			want:    "such wow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure the behaviour is explicitly as requested
			hue.Enabled(tt.enabled)

			got := strconv.Quote(tt.style.Sprint(tt.input))
			want := strconv.Quote(tt.want)

			if got != want {
				t.Errorf("\nGot:\t%v\nWanted:\t%v\n", got, want)
			}
		})
	}
}

func TestSprintf(t *testing.T) {
	tests := []struct {
		name    string    // Name of the test case
		input   string    // Text to style
		want    string    // Expected result including escape sequences
		args    []any     // Args to Sprintf
		enabled bool      // Whether hue is enabled
		style   hue.Style // Style under test
	}{
		{
			name:    "basic",
			input:   "hello %s",
			args:    []any{"hue"},
			style:   hue.BrightGreenBackground,
			enabled: true,
			want:    "\x1b[102mhello hue\x1b[0m",
		},
		{
			name:    "many styles",
			input:   "how many styles %s? %d",
			args:    []any{"hue", 4},
			style:   hue.Blue | hue.BrightGreenBackground | hue.Underline | hue.Dim,
			enabled: true,
			want:    "\x1b[2;4;34;102mhow many styles hue? 4\x1b[0m",
		},
		{
			name:    "basic disabled",
			input:   "hello %s",
			args:    []any{"hue"},
			style:   hue.BrightYellow,
			enabled: false,
			want:    "hello hue",
		},
		{
			name:    "many styles disabled",
			input:   "how many styles %s? %d",
			args:    []any{"hue", 4},
			style:   hue.BrightRed | hue.BrightBlackBackground | hue.Underline | hue.Dim,
			enabled: false,
			want:    "how many styles hue? 4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure the behaviour is explicitly as requested
			hue.Enabled(tt.enabled)

			got := strconv.Quote(tt.style.Sprintf(tt.input, tt.args...))
			want := strconv.Quote(tt.want)

			if got != want {
				t.Errorf("\nGot:\t%v\nWanted:\t%v\n", got, want)
			}
		})
	}
}

func TestSprintln(t *testing.T) {
	tests := []struct {
		name    string    // Name of the test case
		input   string    // Text to style
		want    string    // Expected result including escape sequences
		enabled bool      // Whether hue is enabled
		style   hue.Style // Style under test
	}{
		{
			name:    "basic",
			input:   "woah!",
			style:   hue.White,
			enabled: true,
			want:    "\x1b[37mwoah!\x1b[0m\n",
		},
		{
			name:    "many styles",
			input:   "such wow",
			style:   hue.BrightMagenta | hue.Reverse | hue.GreenBackground,
			enabled: true,
			want:    "\x1b[7;42;95msuch wow\x1b[0m\n",
		},
		{
			name:    "basic disabled",
			input:   "woah!",
			style:   hue.BrightWhite,
			enabled: false,
			want:    "woah!\n",
		},
		{
			name:    "many styles",
			input:   "such wow",
			style:   hue.BrightMagenta | hue.Underline | hue.GreenBackground,
			enabled: false,
			want:    "such wow\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure the behaviour is explicitly as requested
			hue.Enabled(tt.enabled)

			got := strconv.Quote(tt.style.Sprintln(tt.input))
			want := strconv.Quote(tt.want)

			if got != want {
				t.Errorf("\nGot:\t%v\nWanted:\t%v\n", got, want)
			}
		})
	}
}

func TestStyleCode(t *testing.T) {
	tests := []struct {
		name  string    // Name of the test case
		want  string    // Expected string
		style hue.Style // The style under test
	}{
		{name: "bold", style: hue.Bold, want: "1"},
		{name: "dim", style: hue.Dim, want: "2"},
		{name: "italic", style: hue.Italic, want: "3"},
		{name: "underline", style: hue.Underline, want: "4"},
		{name: "reverse", style: hue.Reverse, want: "7"},
		{name: "hidden", style: hue.Hidden, want: "8"},
		{name: "strikethrough", style: hue.Strikethrough, want: "9"},
		{name: "red", style: hue.Black, want: "30"},
		{name: "red", style: hue.Red, want: "31"},
		{name: "green", style: hue.Green, want: "32"},
		{name: "yellow", style: hue.Yellow, want: "33"},
		{name: "blue", style: hue.Blue, want: "34"},
		{name: "magenta", style: hue.Magenta, want: "35"},
		{name: "cyan", style: hue.Cyan, want: "36"},
		{name: "white", style: hue.White, want: "37"},
		{name: "bright black", style: hue.BrightBlack, want: "90"},
		{name: "bright red", style: hue.BrightRed, want: "91"},
		{name: "bright green", style: hue.BrightGreen, want: "92"},
		{name: "bright yellow", style: hue.BrightYellow, want: "93"},
		{name: "bright blue", style: hue.BrightBlue, want: "94"},
		{name: "bright magenta", style: hue.BrightMagenta, want: "95"},
		{name: "bright cyan", style: hue.BrightCyan, want: "96"},
		{name: "bright white", style: hue.BrightWhite, want: "97"},
		{name: "black background", style: hue.BlackBackground, want: "40"},
		{name: "red background", style: hue.RedBackground, want: "41"},
		{name: "green background", style: hue.GreenBackground, want: "42"},
		{name: "yellow background", style: hue.YellowBackground, want: "43"},
		{name: "blue background", style: hue.BlueBackground, want: "44"},
		{name: "magenta background", style: hue.MagentaBackground, want: "45"},
		{name: "cyan background", style: hue.CyanBackground, want: "46"},
		{name: "white background", style: hue.WhiteBackground, want: "47"},
		{name: "bright black background", style: hue.BrightBlackBackground, want: "100"},
		{name: "bright red background", style: hue.BrightRedBackground, want: "101"},
		{name: "bright green background", style: hue.BrightGreenBackground, want: "102"},
		{name: "bright yellow background", style: hue.BrightYellowBackground, want: "103"},
		{name: "bright blue background", style: hue.BrightBlueBackground, want: "104"},
		{name: "bright magenta background", style: hue.BrightMagentaBackground, want: "105"},
		{name: "bright cyan background", style: hue.BrightCyanBackground, want: "106"},
		{name: "bright white background", style: hue.BrightWhiteBackground, want: "107"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hue.Enabled(true)
			got, err := tt.style.Code()
			if err != nil {
				t.Fatalf("Code() returned an error: %v", err)
			}

			if got != tt.want {
				t.Errorf("\nGot:\t%v\nWanted:\t%v\n", got, tt.want)
			}
		})
	}
}

func TestStyleError(t *testing.T) {
	tests := []struct {
		name  string    // Name of the test case
		style hue.Style // Style under test
	}{
		{
			name:  "zero",
			style: 0,
		},
		{
			name:  "too high",
			style: 2199023255553, // > maxStyle
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hue.Enabled(true)
			got, err := tt.style.Code()
			if err == nil {
				t.Fatalf("expected an error, would have got %s", got)
			}
		})
	}
}

func TestStyleCodeCombinations(t *testing.T) {
	tests := []struct {
		name  string    // Name of the test case
		want  string    // Expected string
		style hue.Style // The style under test
	}{
		{
			name:  "bold cyan",
			style: hue.Bold | hue.Cyan,
			want:  "1;36",
		},
		{
			name:  "bold white underlined",
			style: hue.Bold | hue.White | hue.Underline,
			want:  "1;4;37",
		},
		{
			name:  "bold white underlined different order",
			style: hue.White | hue.Underline | hue.Bold,
			want:  "1;4;37",
		},
		{
			name:  "multiple colors",
			style: hue.White | hue.Cyan | hue.Red,
			want:  "31;36;37",
		},
		{
			name:  "lots of everything",
			style: hue.Blue | hue.Red | hue.BlackBackground | hue.Italic | hue.Strikethrough | hue.Bold,
			want:  "1;3;9;31;34;40",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hue.Enabled(true)
			got, err := tt.style.Code()
			if err != nil {
				t.Fatalf("Code() returned an error: %v", err)
			}

			if got != tt.want {
				t.Errorf("\nGot:\t%v\nWanted:\t%v\n", got, tt.want)
			}
		})
	}
}

func TestVisual(t *testing.T) {
	// Run with go test -v, simple visual check to see if we're writing
	// the correct colours
	tests := []struct {
		text  string    // What to write
		style hue.Style // Style under test
	}{
		{style: hue.Bold, text: "Bold"},
		{style: hue.Dim, text: "Dim"},
		{style: hue.Italic, text: "Italic"},
		{style: hue.Underline, text: "Underline"},
		{style: hue.Reverse, text: "Reverse"},
		{style: hue.Hidden, text: "Hidden"},
		{style: hue.Strikethrough, text: "Strikethrough"},
		{style: hue.Red, text: "Red"},
		{style: hue.Green, text: "Green"},
		{style: hue.Yellow, text: "Yellow"},
		{style: hue.Blue, text: "Blue"},
		{style: hue.Magenta, text: "Magenta"},
		{style: hue.Cyan, text: "Cyan"},
		{style: hue.White, text: "White"},
		{style: hue.BrightBlack, text: "BrightBlack"},
		{style: hue.BrightRed, text: "BrightRed"},
		{style: hue.BrightGreen, text: "BrightGreen"},
		{style: hue.BrightYellow, text: "BrightYellow"},
		{style: hue.BrightBlue, text: "BrightBlue"},
		{style: hue.BrightMagenta, text: "BrightMagenta"},
		{style: hue.BrightCyan, text: "BrightCyan"},
		{style: hue.BrightWhite, text: "BrightWhite"},
		{style: hue.BlackBackground, text: "BlackBackground"},
		{style: hue.RedBackground, text: "RedBackground"},
		{style: hue.GreenBackground, text: "GreenBackground"},
		{style: hue.YellowBackground, text: "YellowBackground"},
		{style: hue.BlueBackground, text: "BlueBackground"},
		{style: hue.MagentaBackground, text: "MagentaBackground"},
		{style: hue.CyanBackground, text: "CyanBackground"},
		{style: hue.WhiteBackground, text: "WhiteBackground"},
		{style: hue.Black | hue.Bold, text: "Bold Black"},
		{style: hue.Red | hue.Bold, text: "Bold Red"},
		{style: hue.Green | hue.Bold, text: "Bold Green"},
		{style: hue.Yellow | hue.Bold, text: "Bold Yellow"},
		{style: hue.Blue | hue.Bold, text: "Bold Blue"},
		{style: hue.Magenta | hue.Bold, text: "Bold Magenta"},
		{style: hue.Cyan | hue.Bold, text: "Bold Cyan"},
		{style: hue.White | hue.Bold, text: "Bold White"},
		{style: hue.Black | hue.Underline, text: "Underlined Black"},
		{style: hue.Red | hue.Underline, text: "Underlined Red"},
		{style: hue.Green | hue.Underline, text: "Underlined Green"},
		{style: hue.Yellow | hue.Underline, text: "Underlined Yellow"},
		{style: hue.Blue | hue.Underline, text: "Underlined Blue"},
		{style: hue.Magenta | hue.Underline, text: "Underlined Magenta"},
		{style: hue.Cyan | hue.Underline, text: "Underlined Cyan"},
		{style: hue.White | hue.Underline, text: "Underlined White"},
		{style: hue.Black | hue.Italic, text: "Italic Black"},
		{style: hue.Red | hue.Italic, text: "Italic Red"},
		{style: hue.Green | hue.Italic, text: "Italic Green"},
		{style: hue.Yellow | hue.Italic, text: "Italic Yellow"},
		{style: hue.Blue | hue.Italic, text: "Italic Blue"},
		{style: hue.Magenta | hue.Italic, text: "Italic Magenta"},
		{style: hue.Cyan | hue.Italic, text: "Italic Cyan"},
		{style: hue.White | hue.Italic, text: "Italic White"},
	}
	for _, tt := range tests {
		tt.style.Println(tt.text)
	}
}

func BenchmarkStyle(b *testing.B) {
	hue.Enabled(true)
	b.Run("simple", func(b *testing.B) {
		style := hue.Cyan
		for range b.N {
			_, err := style.Code()
			if err != nil {
				b.Fatalf("Code returned an unexpected error: %v", err)
			}
		}
	})

	b.Run("composite", func(b *testing.B) {
		style := hue.Cyan | hue.WhiteBackground | hue.Bold | hue.Strikethrough
		for range b.N {
			_, err := style.Code()
			if err != nil {
				b.Fatalf("Code returned an unexpected error: %v", err)
			}
		}
	})
}

// captureOutput captures and returns data printed to [os.Stdout] and [os.Stderr] by the provided function fn, allowing
// you to test functions that write to those streams and do not have an option to pass in an [io.Writer].
//
// If the provided function returns a non nil error, the test is failed with the error logged as the reason.
//
// If any error occurs capturing stdout or stderr, the test will also be failed with a descriptive log.
//
//	fn := func() error {
//		fmt.Println("hello stdout")
//		return nil
//	}
//
//	stdout, stderr := test.CaptureOutput(t, fn)
//	fmt.Print(stdout) // "hello stdout\n"
//	fmt.Print(stderr) // ""
func captureOutput(tb testing.TB, fn func() error) (stdout string) {
	tb.Helper()

	// Take copies of the original streams
	oldStdout := os.Stdout

	defer func() {
		// Restore everything back to normal
		os.Stdout = oldStdout
	}()

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		tb.Fatalf("CaptureOutput: could not construct an os.Pipe(): %v", err)
	}

	// Set stdout and stderr streams to the pipe writers
	os.Stdout = stdoutWriter

	stdoutCapture := make(chan string)

	var wg sync.WaitGroup
	wg.Add(1)

	// Copy in goroutines to avoid blocking
	go func(wg *sync.WaitGroup) {
		defer func() {
			close(stdoutCapture)
			wg.Done()
		}()
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, stdoutReader); err != nil {
			tb.Fatalf("CaptureOutput: failed to copy from stdout reader: %v", err)
		}
		stdoutCapture <- buf.String()
	}(&wg)

	// Call the test function that produces the output
	if err := fn(); err != nil {
		tb.Fatalf("CaptureOutput: user function returned an error: %v", err)
	}

	// Close the writers
	stdoutWriter.Close()

	capturedStdout := <-stdoutCapture

	wg.Wait()

	return capturedStdout
}
