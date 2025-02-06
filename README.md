# hue

[![License](https://img.shields.io/github/license/FollowTheProcess/hue)](https://github.com/FollowTheProcess/hue)
[![Go Reference](https://pkg.go.dev/badge/github.com/FollowTheProcess/hue.svg)](https://pkg.go.dev/github.com/FollowTheProcess/hue)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/hue)](https://goreportcard.com/report/github.com/FollowTheProcess/hue)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/hue?logo=github&sort=semver)](https://github.com/FollowTheProcess/hue)
[![CI](https://github.com/FollowTheProcess/hue/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/hue/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/FollowTheProcess/hue/branch/main/graph/badge.svg)](https://codecov.io/gh/FollowTheProcess/hue)

A simple, modern colour/style package for CLI applications in Go

![demo](https://github.com/FollowTheProcess/hue/raw/main/docs/img/demo.gif)

## Project Description

The dominant package in this space for Go is [fatih/color] which I've used before and is very good! However, I want to see if I can make something that improves on it. Specifically I want to try and address the following:

- Alignment/width of colourised text is maintained for [text/tabwriter]
  - Sort of... I cheated and shipped a fork of tabwriter with the right modifications to support ANSI colours (`hue/tabwriter`)
- Support both `$NO_COLOR` and `$FORCE_COLOR`
- Smaller public interface
- Make it so simple you don't even have to think about it
- As low performance overhead as possible

Like most libraries that do this sort of thing, hue uses [ANSI Escape Codes] to instruct the terminal emulator to render particular colours. See [here](https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797) for a helpful breakdown of how these codes work.

> [!IMPORTANT]
> Windows support is best effort, I don't own or use any windows devices so it's not a super high priority for me. If Windows support is important to you, you should use [fatih/color]

## Installation

```shell
go get github.com/FollowTheProcess/hue@latest
```

## Quickstart

Colours and styles in `hue` are implemented as a bitmask and are therefore compile time constants! This means you can do this...

```go
package main

import "github.com/FollowTheProcess/hue"

const (
    success = hue.Green | hue.Bold
    failure = hue.Red | hue.Underline
)

func main() {
    success.Println("It worked!")
    failure.Println("Not really")
}
```

> [!TIP]
> Most functions from the `fmt` package are implemented for hue styles including `Sprintf`, `Fprintln` etc.

### Performance

`hue` has been designed such that each new style is not a new allocated struct, plus the use of bitmasks to encode style leads to some nice performance benefits!

This benchmark measures returning a bold cyan string, one from [fatih/color] and one from `hue`:

```plaintext
❯ go test ./... -run None -benchmem -bench BenchmarkColour
goos: darwin
goarch: arm64
pkg: github.com/FollowTheProcess/hue
cpu: Apple M1 Pro
BenchmarkColour/hue-8            7497607               138.7 ns/op            80 B/op          3 allocs/op
BenchmarkColour/color-8          2893501               415.2 ns/op           248 B/op         12 allocs/op
PASS
ok      github.com/FollowTheProcess/hue 3.044s
```

- Nearly 70% faster
- Nearly 70% less bytes copied
- 9 fewer heap allocations!

> [!NOTE]
> This benchmark used to be in here and run as part of CI, but I found that `fatih/color` would show up as an indirect dependency
> in code that imported `hue` so I got rid of it, you'll just have to trust me I guess 😂

### Tabwriter

A common issue with ANSI colour libraries in Go are that they don't play well with [text/tabwriter]. This is because tabwriter includes the ANSI escape sequence in the cell width calculations, throwing off where it aligns the columns.

This has always annoyed me so when making my own ANSI colour library I had to tackle the tabwriter issue!

Enter `hue/tabwriter` a drop in replacement for [text/tabwriter] that doesn't care about ANSI colours ✨

![tabwriter](https://github.com/FollowTheProcess/hue/raw/main/docs/img/tabwriter.gif)

> [!NOTE]
> The actual change is incredibly simple, just teaching [text/tabwriter] to ignore ANSI codes when it sees them so compatibility
> should be seamless

### Credits

This package was created with [copier] and the [FollowTheProcess/go_copier] project template.

[copier]: https://copier.readthedocs.io/en/stable/
[FollowTheProcess/go_copier]: https://github.com/FollowTheProcess/go_copier
[fatih/color]: https://github.com/fatih/color
[text/tabwriter]: https://pkg.go.dev/text/tabwriter
[ANSI Escape Codes]: https://en.wikipedia.org/wiki/ANSI_escape_code
