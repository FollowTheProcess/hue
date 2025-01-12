package hue_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/FollowTheProcess/hue"
	"github.com/FollowTheProcess/test"
)

func TestFprint(t *testing.T) {
	tests := []struct {
		name  string    // Name of the test case
		input string    // Text to style
		want  string    // Expected result including escape sequences
		style hue.Style // Style under test
	}{
		{
			name:  "basic",
			input: "hello",
			style: hue.Green,
			want:  "\x1b[32mhello\x1b[0m",
		},
		{
			name:  "many styles",
			input: "hello",
			style: hue.Green | hue.BlueBackground | hue.Bold | hue.Underline,
			want:  "\x1b[1;4;32;44mhello\x1b[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			tt.style.Fprint(buf, tt.input)

			got := strconv.Quote(buf.String())
			want := strconv.Quote(tt.want)

			test.Equal(t, got, want)
		})
	}
}

func TestFprintf(t *testing.T) {
	tests := []struct {
		name  string    // Name of the test case
		input string    // Text to style
		want  string    // Expected result including escape sequences
		args  []any     // Args to Fprintf
		style hue.Style // Style under test
	}{
		{
			name:  "basic",
			input: "hello %s",
			args:  []any{"hue"},
			style: hue.Magenta,
			want:  "\x1b[35mhello hue\x1b[0m",
		},
		{
			name:  "many styles",
			input: "how many styles %s? %d",
			args:  []any{"hue", 4},
			style: hue.Blue | hue.RedBackground | hue.Italic | hue.Bold,
			want:  "\x1b[1;3;34;41mhow many styles hue? 4\x1b[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			tt.style.Fprintf(buf, tt.input, tt.args...)

			got := strconv.Quote(buf.String())
			want := strconv.Quote(tt.want)

			test.Equal(t, got, want)
		})
	}
}

func TestFprintln(t *testing.T) {
	tests := []struct {
		name  string    // Name of the test case
		input string    // Text to style
		want  string    // Expected result including escape sequences
		style hue.Style // Style under test
	}{
		{
			name:  "basic",
			input: "woah!",
			style: hue.BrightGreen,
			want:  "\x1b[92mwoah!\n\x1b[0m",
		},
		{
			name:  "many styles",
			input: "such wow",
			style: hue.BrightCyan | hue.Strikethrough | hue.BlinkSlow,
			want:  "\x1b[5;9;96msuch wow\n\x1b[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			tt.style.Fprintln(buf, tt.input)

			got := strconv.Quote(buf.String())
			want := strconv.Quote(tt.want)

			test.Equal(t, got, want)
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
		{name: "blink slow", style: hue.BlinkSlow, want: "5"},
		{name: "blink fast", style: hue.BlinkFast, want: "6"},
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
			got, err := tt.style.Code()
			test.Ok(t, err)
			test.Equal(t, got, tt.want)
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
			got, err := tt.style.Code()
			test.Err(t, err, test.Context("would have got %s", got))
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
			style: hue.Blue | hue.Red | hue.BlackBackground | hue.BlinkFast | hue.Strikethrough | hue.Bold,
			want:  "1;6;9;31;34;40",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.style.Code()
			test.Ok(t, err)
			test.Equal(t, got, tt.want)
		})
	}
}

func BenchmarkStyle(b *testing.B) {
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
