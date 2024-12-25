<div align="center">
  <h3 align="center">Gotex</h3>
  <p align="center">
        A cross-platform TUI Go Test Explorer for Go.
  </p>
</div>

![Gotex Screenshot][img/screenshot-gotex.png]

# gotex

A Terminal User Interface (TUI) for discovering and executing Go tests with ease. Gotex provides an intuitive interface to navigate, filter and run your Go test suite from the terminal.

## Features

- Interactive terminal interface for browsing test packages, tests and test cases
- Fuzzy find tests or navigate the tree
- Allows for results to be parsed through other tools
- Vim like keybindings
- Rerun last test

## Installation

```bash
go install github.com/sgrumley/gotex@latest
```

## Usage

Navigate to your Go project directory and run:

```bash
gotex
```

### Keyboard Controls

- `j/k or ↑/↓`: Navigate through tests
- `h/l`: Expand or close child tests
- `r`: Run selected test
- `R`: Rerun last test
- `A`: Run all tests
- `/`: Search tests
- `s`: Sync project
- `c`: Config display
- `C`: Debug Console
- `ctrl-u/ ctrl-d`: Scroll text
- `q`: Quit


## Configuration

Gotex can be configured using a provided yaml file:

```yaml
# Default configuration
json: false
timeout:
short: false
verbose: true
failfast: false
cover: true
# pipeto: tparse

```

## Requirements

- Go 1.16 or higher
- Terminal with 256 color support

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by similar TUI tools like [gotests](https://github.com/cweill/gotests)
- Built with [tview](https://github.com/rivo/tview)

