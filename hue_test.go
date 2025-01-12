package hue_test

import (
	"strconv"
	"testing"

	"github.com/FollowTheProcess/hue"
	"github.com/FollowTheProcess/test"
)

func TestStyleString(t *testing.T) {
	tests := []struct {
		name  string    // Name of the test case
		want  string    // Expected string
		style hue.Style // The style under test
	}{
		{
			name:  "above max",
			style: hue.Style(2199023255552),
			want:  "invalid style: Style(2199023255552)",
		},
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
			got := tt.style.String()
			test.Equal(t, got, tt.want)
		})
	}
}

func TestStyleStringCombinations(t *testing.T) {
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
			got := tt.style.String()
			test.Equal(t, got, tt.want)
		})
	}
}

func TestColour(t *testing.T) {
	tests := []struct {
		name    string    // Name of the test case
		text    string    // The message to colour
		want    string    // Expected (raw) output
		enabled bool      // What to set hue.Enabled to
		style   hue.Style // The style to apply
	}{
		{
			name:    "basic",
			text:    "hello",
			style:   hue.Green,
			enabled: true,
			want:    "\x1b[32mhello\x1b[0m",
		},
		{
			name:    "basic",
			text:    "hello",
			style:   hue.Green,
			enabled: false,
			want:    "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hue.Enabled = tt.enabled
			got := strconv.Quote(tt.style.Sprint(tt.text))
			want := strconv.Quote(tt.want)

			test.Equal(t, got, want)
		})
	}
}

func BenchmarkStyle(b *testing.B) {
	b.Run("simple", func(b *testing.B) {
		style := hue.Cyan
		for range b.N {
			_ = style.String()
		}
	})

	b.Run("composite", func(b *testing.B) {
		style := hue.Cyan | hue.WhiteBackground | hue.Bold | hue.Strikethrough
		for range b.N {
			_ = style.String()
		}
	})
}
