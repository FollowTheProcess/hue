# hue

[![License](https://img.shields.io/github/license/FollowTheProcess/hue)](https://github.com/FollowTheProcess/hue)
[![Go Reference](https://pkg.go.dev/badge/github.com/FollowTheProcess/hue.svg)](https://pkg.go.dev/github.com/FollowTheProcess/hue)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/hue)](https://goreportcard.com/report/github.com/FollowTheProcess/hue)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/hue?logo=github&sort=semver)](https://github.com/FollowTheProcess/hue)
[![CI](https://github.com/FollowTheProcess/hue/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/hue/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/FollowTheProcess/hue/branch/main/graph/badge.svg)](https://codecov.io/gh/FollowTheProcess/hue)

A simple, modern colour/style package for CLI applications in Go

> [!WARNING]
> **hue is in early development and is not yet ready for use**

![caution](./img/caution.png)

## Project Description

The dominant package in this space for Go is [fatih/color] which I've used before and is very good! However, I want to see if I can make something that improves on it. Specifically I want to try and address the following:

- Alignment/width of colourised text is maintained for [text/tabwriter]
- Support both `$NO_COLOR` and `$FORCE_COLOR`
- Smaller public interface, more simple
- Zero allocations (may not be possible)

Like most libraries that do this sort of thing, hue uses [ANSI Escape Codes] to instruct the terminal emulator to render particular colours. See [here](https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797) for a helpful breakdown of how these codes work.

> [!NOTE]
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

### Credits

This package was created with [copier] and the [FollowTheProcess/go_copier] project template.

[copier]: https://copier.readthedocs.io/en/stable/
[FollowTheProcess/go_copier]: https://github.com/FollowTheProcess/go_copier
[fatih/color]: https://github.com/fatih/color
[text/tabwriter]: https://pkg.go.dev/text/tabwriter
[ANSI Escape Codes]: https://en.wikipedia.org/wiki/ANSI_escape_code
