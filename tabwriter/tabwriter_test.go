// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tabwriter_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"

	"go.followtheprocess.codes/hue"
	"go.followtheprocess.codes/hue/tabwriter"
)

type buffer struct {
	a []byte
}

func (b *buffer) init(n int) { b.a = make([]byte, 0, n) }

func (b *buffer) clear() { b.a = b.a[0:0] }

func (b *buffer) Write(buf []byte) (written int, err error) {
	n := len(b.a)
	m := len(buf)

	if n+m > cap(b.a) {
		panic("buffer.Write: buffer too small")
	}

	b.a = b.a[0 : n+m]
	for i := range m {
		b.a[n+i] = buf[i]
	}

	return len(buf), nil
}

func (b *buffer) String() string { return string(b.a) }

func write(t *testing.T, testname string, w *tabwriter.Writer, src string) {
	t.Helper()

	written, err := io.WriteString(w, src)
	if err != nil {
		t.Errorf("--- test: %s\n--- src:\n%q\n--- write error: %v\n", testname, src, err)
	}

	if written != len(src) {
		t.Errorf(
			"--- test: %s\n--- src:\n%q\n--- written = %d, len(src) = %d\n",
			testname,
			src,
			written,
			len(src),
		)
	}
}

func verify( //nolint: revive // argument-limit
	t *testing.T,
	testname string,
	w *tabwriter.Writer,
	b *buffer,
	src, expected string,
) {
	t.Helper()

	err := w.Flush()
	if err != nil {
		t.Errorf("--- test: %s\n--- src:\n%q\n--- flush error: %v\n", testname, src, err)
	}

	res := b.String()
	if res != expected {
		t.Errorf(
			"--- test: %s\n--- src:\n%q\n--- found:\n%q\n--- expected:\n%q\n",
			testname,
			src,
			res,
			expected,
		)
	}
}

func check( //nolint: revive // argument-limit
	t *testing.T,
	testname string,
	minwidth, tabwidth, padding int,
	padchar byte,
	flags uint,
	src, expected string,
) {
	t.Helper()

	var b buffer

	b.init(1000)

	var w tabwriter.Writer

	w.Init(&b, minwidth, tabwidth, padding, padchar, flags)

	// write all at once
	title := testname + " (written all at once)"

	b.clear()
	write(t, title, &w, src)
	verify(t, title, &w, &b, src, expected)

	// write byte-by-byte
	title = testname + " (written byte-by-byte)"

	b.clear()

	for i := range len(src) {
		write(t, title, &w, src[i:i+1])
	}

	verify(t, title, &w, &b, src, expected)

	// write using Fibonacci slice sizes
	title = testname + " (written in fibonacci slices)"

	b.clear()

	for i, d := 0, 0; i < len(src); {
		write(t, title, &w, src[i:i+d])

		i, d = i+d, d+1
		if i+d > len(src) {
			d = len(src) - i
		}
	}

	verify(t, title, &w, &b, src, expected)
}

var tests = []struct {
	testname                    string
	src, expected               string
	minwidth, tabwidth, padding int
	flags                       uint
	padchar                     byte
}{
	{
		testname: "1a",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "",
		expected: "",
	},

	{
		testname: "1a debug",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: tabwriter.Debug,
		src:      "",
		expected: "",
	},

	{
		testname: "1b esc stripped",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: tabwriter.StripEscape,
		src:      "\xff\xff",
		expected: "",
	},

	{
		testname: "1b esc",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "\xff\xff",
		expected: "\xff\xff",
	},

	{
		testname: "1c esc stripped",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: tabwriter.StripEscape,
		src:      "\xff\t\xff",
		expected: "\t",
	},

	{
		testname: "1c esc",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "\xff\t\xff",
		expected: "\xff\t\xff",
	},

	{
		testname: "1d esc stripped",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: tabwriter.StripEscape,
		src:      "\xff\"foo\t\n\tbar\"\xff",
		expected: "\"foo\t\n\tbar\"",
	},

	{
		testname: "1d esc",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "\xff\"foo\t\n\tbar\"\xff",
		expected: "\xff\"foo\t\n\tbar\"\xff",
	},

	{
		testname: "1e esc stripped",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: tabwriter.StripEscape,
		src:      "abc\xff\tdef", // unterminated escape
		expected: "abc\tdef",
	},

	{
		testname: "1e esc",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "abc\xff\tdef", // unterminated escape
		expected: "abc\xff\tdef",
	},

	{
		testname: "1f esc ansi",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "abc\x1b[\tdef", // unterminated ANSI escape sequence
		expected: "abc\x1b[\tdef",
	},

	{
		testname: "2",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "\n\n\n",
		expected: "\n\n\n",
	},

	{
		testname: "3",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "a\nb\nc",
		expected: "a\nb\nc",
	},

	{
		testname: "3 colours",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "a\nb\n\x1b[93;41mc\x1b[0m",
		expected: "a\nb\n\x1b[93;41mc\x1b[0m",
	},

	{
		testname: "4a",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "\t", // '\t' terminates an empty cell on last line - nothing to print
		expected: "",
	},

	{
		testname: "4b",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: tabwriter.AlignRight,
		src:      "\t", // '\t' terminates an empty cell on last line - nothing to print
		expected: "",
	},

	{
		testname: "5",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "*\t*",
		expected: "*.......*",
	},

	{
		testname: "5b",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "*\t*\n",
		expected: "*.......*\n",
	},

	{
		testname: "5c",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "*\t*\t",
		expected: "*.......*",
	},

	{
		testname: "5c debug",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: tabwriter.Debug,
		src:      "*\t*\t",
		expected: "*.......|*",
	},

	{
		testname: "5d",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: tabwriter.AlignRight,
		src:      "*\t*\t",
		expected: ".......**",
	},

	{
		testname: "6",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "\t\n",
		expected: "........\n",
	},

	{
		testname: "7a",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "a) foo",
		expected: "a) foo",
	},

	{
		testname: "7b",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: ' ', flags: 0,
		src:      "b) foo\tbar",
		expected: "b) foo  bar",
	},

	{
		testname: "7b colours",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: ' ', flags: 0,
		src:      "b) \x1b[93;41mfoo\x1b[0m\tbar",
		expected: "b) \x1b[93;41mfoo\x1b[0m  bar",
	},

	{
		testname: "7c",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "c) foo\tbar\t",
		expected: "c) foo..bar",
	},

	{
		testname: "7c colours",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "c) \x1b[93;41mfoo\x1b[0m\tbar\t",
		expected: "c) \x1b[93;41mfoo\x1b[0m..bar",
	},

	{
		testname: "7d",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "d) foo\tbar\n",
		expected: "d) foo..bar\n",
	},

	{
		testname: "7d colours",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "d) \x1b[93;41mfoo\x1b[0m\tbar\n",
		expected: "d) \x1b[93;41mfoo\x1b[0m..bar\n",
	},

	{
		testname: "7e",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "e) foo\tbar\t\n",
		expected: "e) foo..bar.....\n",
	},

	{
		testname: "7e colours",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src:      "e) \x1b[93;41mfoo\x1b[0m\tbar\t\n",
		expected: "e) \x1b[93;41mfoo\x1b[0m..bar.....\n",
	},

	{
		testname: "7f",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: tabwriter.FilterHTML,
		src:      "f) f&lt;o\t<b>bar</b>\t\n",
		expected: "f) f&lt;o..<b>bar</b>.....\n",
	},

	{
		testname: "7g",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: tabwriter.FilterHTML,
		src:      "g) f&lt;o\t<b>bar</b>\t non-terminated entity &amp",
		expected: "g) f&lt;o..<b>bar</b>..... non-terminated entity &amp",
	},

	{
		testname: "7g debug",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: tabwriter.FilterHTML | tabwriter.Debug,
		src:      "g) f&lt;o\t<b>bar</b>\t non-terminated entity &amp",
		expected: "g) f&lt;o..|<b>bar</b>.....| non-terminated entity &amp",
	},

	{
		testname: "8",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '*', flags: 0,
		src:      "Hello, world!\n",
		expected: "Hello, world!\n",
	},

	{
		testname: "9a",
		minwidth: 1, tabwidth: 0, padding: 0, padchar: '.', flags: 0,
		src: "1\t2\t3\t4\n" +
			"11\t222\t3333\t44444\n",

		expected: "1.2..3...4\n" +
			"11222333344444\n",
	},

	{
		testname: "9b",
		minwidth: 1, tabwidth: 0, padding: 0, padchar: '.', flags: tabwriter.FilterHTML,
		src: "1\t2<!---\f--->\t3\t4\n" + // \f inside HTML is ignored
			"11\t222\t3333\t44444\n",

		expected: "1.2<!---\f--->..3...4\n" +
			"11222333344444\n",
	},

	{
		testname: "9c",
		minwidth: 1, tabwidth: 0, padding: 0, padchar: '.', flags: 0,
		src: "1\t2\t3\t4\f" + // \f causes a newline and flush
			"11\t222\t3333\t44444\n",

		expected: "1234\n" +
			"11222333344444\n",
	},

	{
		testname: "9c debug",
		minwidth: 1, tabwidth: 0, padding: 0, padchar: '.', flags: tabwriter.Debug,
		src: "1\t2\t3\t4\f" + // \f causes a newline and flush
			"11\t222\t3333\t44444\n",

		expected: "1|2|3|4\n" +
			"---\n" +
			"11|222|3333|44444\n",
	},

	{
		testname: "10a",
		minwidth: 5, tabwidth: 0, padding: 0, padchar: '.', flags: 0,
		src:      "1\t2\t3\t4\n",
		expected: "1....2....3....4\n",
	},

	{
		testname: "10a colours",
		minwidth: 5, tabwidth: 0, padding: 0, padchar: '.', flags: 0,
		src:      "1\t2\t\x1b[93;41m3\x1b[0m\t4\n",
		expected: "1....2....\x1b[93;41m3\x1b[0m....4\n",
	},

	{
		testname: "10b",
		minwidth: 5, tabwidth: 0, padding: 0, padchar: '.', flags: 0,
		src:      "1\t2\t3\t4\t\n",
		expected: "1....2....3....4....\n",
	},

	{
		testname: "10b colours",
		minwidth: 5, tabwidth: 0, padding: 0, padchar: '.', flags: 0,
		src:      "1\t2\t\x1b[93;41m3\x1b[0m\t4\t\n",
		expected: "1....2....\x1b[93;41m3\x1b[0m....4....\n",
	},

	{
		testname: "11",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '.', flags: 0,
		src: "本\tb\tc\n" +
			"aa\t\u672c\u672c\u672c\tcccc\tddddd\n" +
			"aaa\tbbbb\n",

		expected: "本.......b.......c\n" +
			"aa......本本本.....cccc....ddddd\n" +
			"aaa.....bbbb\n",
	},

	{
		testname: "12a",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: ' ', flags: tabwriter.AlignRight,
		src: "a\tè\tc\t\n" +
			"aa\tèèè\tcccc\tddddd\t\n" +
			"aaa\tèèèè\t\n",

		expected: "       a       è       c\n" +
			"      aa     èèè    cccc   ddddd\n" +
			"     aaa    èèèè\n",
	},

	{
		testname: "12b",
		minwidth: 2, tabwidth: 0, padding: 0, padchar: ' ', flags: 0,
		src: "a\tb\tc\n" +
			"aa\tbbb\tcccc\n" +
			"aaa\tbbbb\n",

		expected: "a  b  c\n" +
			"aa bbbcccc\n" +
			"aaabbbb\n",
	},

	{
		testname: "12c",
		minwidth: 8, tabwidth: 0, padding: 1, padchar: '_', flags: 0,
		src: "a\tb\tc\n" +
			"aa\tbbb\tcccc\n" +
			"aaa\tbbbb\n",

		expected: "a_______b_______c\n" +
			"aa______bbb_____cccc\n" +
			"aaa_____bbbb\n",
	},

	{
		testname: "13a",
		minwidth: 4, tabwidth: 0, padding: 1, padchar: '-', flags: 0,
		src: "4444\t日本語\t22\t1\t333\n" +
			"999999999\t22\n" +
			"7\t22\n" +
			"\t\t\t88888888\n" +
			"\n" +
			"666666\t666666\t666666\t4444\n" +
			"1\t1\t999999999\t0000000000\n",

		expected: "4444------日本語-22--1---333\n" +
			"999999999-22\n" +
			"7---------22\n" +
			"------------------88888888\n" +
			"\n" +
			"666666-666666-666666----4444\n" +
			"1------1------999999999-0000000000\n",
	},

	{
		testname: "13b",
		minwidth: 4, tabwidth: 0, padding: 3, padchar: '.', flags: 0,
		src: "4444\t333\t22\t1\t333\n" +
			"999999999\t22\n" +
			"7\t22\n" +
			"\t\t\t88888888\n" +
			"\n" +
			"666666\t666666\t666666\t4444\n" +
			"1\t1\t999999999\t0000000000\n",

		expected: "4444........333...22...1...333\n" +
			"999999999...22\n" +
			"7...........22\n" +
			"....................88888888\n" +
			"\n" +
			"666666...666666...666666......4444\n" +
			"1........1........999999999...0000000000\n",
	},

	{
		testname: "13c",
		minwidth: 8, tabwidth: 8, padding: 1, padchar: '\t', flags: tabwriter.FilterHTML,
		src: "4444\t333\t22\t1\t333\n" +
			"999999999\t22\n" +
			"7\t22\n" +
			"\t\t\t88888888\n" +
			"\n" +
			"666666\t666666\t666666\t4444\n" +
			"1\t1\t<font color=red attr=日本語>999999999</font>\t0000000000\n",

		expected: "4444\t\t333\t22\t1\t333\n" +
			"999999999\t22\n" +
			"7\t\t22\n" +
			"\t\t\t\t88888888\n" +
			"\n" +
			"666666\t666666\t666666\t\t4444\n" +
			"1\t1\t<font color=red attr=日本語>999999999</font>\t0000000000\n",
	},

	{
		testname: "14",
		minwidth: 1, tabwidth: 0, padding: 2, padchar: ' ', flags: tabwriter.AlignRight,
		src: ".0\t.3\t2.4\t-5.1\t\n" +
			"23.0\t12345678.9\t2.4\t-989.4\t\n" +
			"5.1\t12.0\t2.4\t-7.0\t\n" +
			".0\t0.0\t332.0\t8908.0\t\n" +
			".0\t-.3\t456.4\t22.1\t\n" +
			".0\t1.2\t44.4\t-13.3\t\t",

		expected: "    .0          .3    2.4    -5.1\n" +
			"  23.0  12345678.9    2.4  -989.4\n" +
			"   5.1        12.0    2.4    -7.0\n" +
			"    .0         0.0  332.0  8908.0\n" +
			"    .0         -.3  456.4    22.1\n" +
			"    .0         1.2   44.4   -13.3",
	},

	{
		testname: "14 debug",
		minwidth: 1, tabwidth: 0, padding: 2, padchar: ' ', flags: tabwriter.AlignRight | tabwriter.Debug,
		src: ".0\t.3\t2.4\t-5.1\t\n" +
			"23.0\t12345678.9\t2.4\t-989.4\t\n" +
			"5.1\t12.0\t2.4\t-7.0\t\n" +
			".0\t0.0\t332.0\t8908.0\t\n" +
			".0\t-.3\t456.4\t22.1\t\n" +
			".0\t1.2\t44.4\t-13.3\t\t",

		expected: "    .0|          .3|    2.4|    -5.1|\n" +
			"  23.0|  12345678.9|    2.4|  -989.4|\n" +
			"   5.1|        12.0|    2.4|    -7.0|\n" +
			"    .0|         0.0|  332.0|  8908.0|\n" +
			"    .0|         -.3|  456.4|    22.1|\n" +
			"    .0|         1.2|   44.4|   -13.3|",
	},

	{
		testname: "15a",
		minwidth: 4, tabwidth: 0, padding: 0, padchar: '.', flags: 0,
		src:      "a\t\tb",
		expected: "a.......b",
	},

	{
		testname: "15a colours",
		minwidth: 4, tabwidth: 0, padding: 0, padchar: '.', flags: 0,
		src:      "\x1b[93;41ma\x1b[0m\t\tb",
		expected: "\x1b[93;41ma\x1b[0m.......b",
	},

	{
		testname: "15b",
		minwidth: 4, tabwidth: 0, padding: 0, padchar: '.', flags: tabwriter.DiscardEmptyColumns,
		src:      "a\t\tb", // htabs - do not discard column
		expected: "a.......b",
	},

	{
		testname: "15c",
		minwidth: 4, tabwidth: 0, padding: 0, padchar: '.', flags: tabwriter.DiscardEmptyColumns,
		src:      "a\v\vb",
		expected: "a...b",
	},

	{
		testname: "15d",
		minwidth: 4, tabwidth: 0, padding: 0, padchar: '.', flags: tabwriter.AlignRight | tabwriter.DiscardEmptyColumns,
		src:      "a\v\vb",
		expected: "...ab",
	},

	{
		testname: "16a",
		minwidth: 100, tabwidth: 100, padding: 0, padchar: '\t', flags: 0,
		src: "a\tb\t\td\n" +
			"a\tb\t\td\te\n" +
			"a\n" +
			"a\tb\tc\td\n" +
			"a\tb\tc\td\te\n",

		expected: "a\tb\t\td\n" +
			"a\tb\t\td\te\n" +
			"a\n" +
			"a\tb\tc\td\n" +
			"a\tb\tc\td\te\n",
	},

	{
		testname: "16b",
		minwidth: 100, tabwidth: 100, padding: 0, padchar: '\t', flags: tabwriter.DiscardEmptyColumns,
		src: "a\vb\v\vd\n" +
			"a\vb\v\vd\ve\n" +
			"a\n" +
			"a\vb\vc\vd\n" +
			"a\vb\vc\vd\ve\n",

		expected: "a\tb\td\n" +
			"a\tb\td\te\n" +
			"a\n" +
			"a\tb\tc\td\n" +
			"a\tb\tc\td\te\n",
	},

	{
		testname: "16b debug",
		minwidth: 100, tabwidth: 100, padding: 0, padchar: '\t', flags: tabwriter.DiscardEmptyColumns | tabwriter.Debug,
		src: "a\vb\v\vd\n" +
			"a\vb\v\vd\ve\n" +
			"a\n" +
			"a\vb\vc\vd\n" +
			"a\vb\vc\vd\ve\n",

		expected: "a\t|b\t||d\n" +
			"a\t|b\t||d\t|e\n" +
			"a\n" +
			"a\t|b\t|c\t|d\n" +
			"a\t|b\t|c\t|d\t|e\n",
	},

	{
		testname: "16c",
		minwidth: 100, tabwidth: 100, padding: 0, padchar: '\t', flags: tabwriter.DiscardEmptyColumns,
		src: "a\tb\t\td\n" + // hard tabs - do not discard column
			"a\tb\t\td\te\n" +
			"a\n" +
			"a\tb\tc\td\n" +
			"a\tb\tc\td\te\n",

		expected: "a\tb\t\td\n" +
			"a\tb\t\td\te\n" +
			"a\n" +
			"a\tb\tc\td\n" +
			"a\tb\tc\td\te\n",
	},

	{
		testname: "16c debug",
		minwidth: 100, tabwidth: 100, padding: 0, padchar: '\t', flags: tabwriter.DiscardEmptyColumns | tabwriter.Debug,
		src: "a\tb\t\td\n" + // hard tabs - do not discard column
			"a\tb\t\td\te\n" +
			"a\n" +
			"a\tb\tc\td\n" +
			"a\tb\tc\td\te\n",

		expected: "a\t|b\t|\t|d\n" +
			"a\t|b\t|\t|d\t|e\n" +
			"a\n" +
			"a\t|b\t|c\t|d\n" +
			"a\t|b\t|c\t|d\t|e\n",
	},
}

func Test(t *testing.T) {
	for _, e := range tests {
		check(
			t,
			e.testname,
			e.minwidth,
			e.tabwidth,
			e.padding,
			e.padchar,
			e.flags,
			e.src,
			e.expected,
		)
	}
}

type panicWriter struct{}

func (panicWriter) Write([]byte) (int, error) {
	panic("cannot write")
}

func wantPanicString(t *testing.T, want string) {
	t.Helper()

	if e := recover(); e != nil { //nolint: revive // This is deferred
		got, ok := e.(string)

		switch {
		case !ok:
			t.Errorf("got %v (%T), want panic string", e, e)
		case got != want:
			t.Errorf("wrong panic message: got %q, want %q", got, want)
		}
	}
}

func TestPanicDuringFlush(t *testing.T) {
	defer wantPanicString(t, "tabwriter: panic during Flush (cannot write)")

	var p panicWriter

	w := new(tabwriter.Writer)
	w.Init(p, 0, 0, 5, ' ', 0)
	io.WriteString(w, "a") //nolint: errcheck
	w.Flush()
	t.Errorf("failed to panic during Flush")
}

func TestPanicDuringWrite(t *testing.T) {
	defer wantPanicString(t, "tabwriter: panic during Write (cannot write)")

	var p panicWriter

	w := new(tabwriter.Writer)
	w.Init(p, 0, 0, 5, ' ', 0)
	// the second \n triggers a call to w.Write and thus a panic
	io.WriteString(w, "a\n\n") //nolint: errcheck
	t.Errorf("failed to panic during Write")
}

func TestVisual(t *testing.T) {
	hue.Enabled(true) // go test buffers output so autodetection disabled colour

	writer := tabwriter.NewWriter(os.Stdout, 1, 8, 2, ' ', 0)
	defer writer.Flush()

	green := hue.Green
	cyan := hue.Cyan
	boldRed := hue.Red | hue.Bold
	strikeThroughYellow := hue.Yellow | hue.Strikethrough
	magenta := hue.Magenta
	blueUnderline := hue.Blue | hue.Underline
	brightGreen := hue.BrightGreen
	italicMagenta := hue.Magenta | hue.Italic
	grey := hue.BrightBlack
	brightBlue := hue.BrightBlue

	fmt.Fprintf(
		writer,
		"%s\t%s\t%s\t%s\t\n",
		green.Sprint("Green"),
		cyan.Sprint("Cyan"),
		boldRed.Sprint("BoldRed"),
		strikeThroughYellow.Sprint("Strikethrough Yellow"),
	)
	fmt.Fprintf(
		writer,
		"%s\t%s\t%s\t%s\t\n",
		cyan.Sprint("Look"),
		strikeThroughYellow.Sprint("Colours"),
		green.Sprint("In"),
		boldRed.Sprint("Tables!"),
	)
	fmt.Fprintf(
		writer,
		"%s\t%s\t%s\t%s\t\n",
		magenta.Sprint("All"),
		blueUnderline.Sprint("Properly"),
		green.Sprint("Lined"),
		strikeThroughYellow.Sprint("Up"),
	)
	fmt.Fprintf(
		writer,
		"%s\t%s\t%s\t%s\t\n",
		brightGreen.Sprint("How"),
		italicMagenta.Sprint("Cool"),
		grey.Sprint("Is"),
		brightBlue.Sprint("That!"),
	)
}

func BenchmarkTable(b *testing.B) {
	for _, w := range [...]int{1, 10, 100} {
		// Build a line with w cells.
		line := bytes.Repeat([]byte("a\t"), w)
		line = append(line, '\n')

		for _, h := range [...]int{10, 1000, 100000} {
			b.Run(fmt.Sprintf("%dx%d", w, h), func(b *testing.B) {
				b.Run("new", func(b *testing.B) {
					b.ReportAllocs()

					for b.Loop() {
						w := tabwriter.NewWriter(
							io.Discard,
							4,
							4,
							1,
							' ',
							0,
						) // no particular reason for these settings
						// Write the line h times.
						for range h {
							w.Write(line) //nolint: errcheck
						}

						w.Flush()
					}
				})

				b.Run("reuse", func(b *testing.B) {
					b.ReportAllocs()

					w := tabwriter.NewWriter(
						io.Discard,
						4,
						4,
						1,
						' ',
						0,
					) // no particular reason for these settings

					for range b.N {
						// Write the line h times.
						for range h {
							w.Write(line) //nolint: errcheck
						}

						w.Flush()
					}
				})
			})
		}
	}
}

func BenchmarkPyramid(b *testing.B) {
	for _, x := range [...]int{10, 100, 1000} {
		// Build a line with x cells.
		line := bytes.Repeat([]byte("a\t"), x)
		b.Run(strconv.Itoa(x), func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				w := tabwriter.NewWriter(
					io.Discard,
					4,
					4,
					1,
					' ',
					0,
				) // no particular reason for these settings
				// Write increasing prefixes of that line.
				for j := range x {
					w.Write(line[:j*2])   //nolint: errcheck
					w.Write([]byte{'\n'}) //nolint: errcheck
				}

				w.Flush()
			}
		})
	}
}

func BenchmarkRagged(b *testing.B) {
	var lines [8][]byte
	for i, w := range [8]int{6, 2, 9, 5, 5, 7, 3, 8} {
		// Build a line with w cells.
		lines[i] = bytes.Repeat([]byte("a\t"), w)
	}

	for _, h := range [...]int{10, 100, 1000} {
		b.Run(strconv.Itoa(h), func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				w := tabwriter.NewWriter(
					io.Discard,
					4,
					4,
					1,
					' ',
					0,
				) // no particular reason for these settings
				// Write the lines in turn h times.
				for j := range h {
					w.Write(lines[j%len(lines)]) //nolint: errcheck
					w.Write([]byte{'\n'})        //nolint: errcheck
				}

				w.Flush()
			}
		})
	}
}

const codeSnippet = `
some command

foo	# aligned
barbaz	# comments

but
mostly
single
cell
lines
`

func BenchmarkCode(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		w := tabwriter.NewWriter(
			io.Discard,
			4,
			4,
			1,
			' ',
			0,
		) // no particular reason for these settings
		// The code is small, so it's reasonable for the tabwriter user
		// to write it all at once, or buffer the writes.
		w.Write([]byte(codeSnippet)) //nolint: errcheck
		w.Flush()
	}
}
