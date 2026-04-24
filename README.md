# go-sea

`go-sea` is a small example CLI that builds to the `gsea` binary and demonstrates how to combine these Go libraries in one project:

- [`spf13/cobra`](https://github.com/spf13/cobra) for commands and flags
- [`spf13/viper`](https://github.com/spf13/viper) for config files, environment variables, and flag binding
- [`adrg/xdg`](https://github.com/adrg/xdg) for XDG-aware config discovery and creation
- [`fatih/color`](https://github.com/fatih/color) for colored terminal output

The main example command prints a short, styled terminal message, renders a compact preview of a curated Monokai palette subset, and reports where its configuration came from.

## Requirements

- Go `1.25.8`
- `direnv`
- `task`

If `direnv` is enabled for the repo, `.envrc` will verify the Go version from `go.mod` and install local copies of Go and Task under `.tools/` when needed.

## Quick Start

```bash
direnv allow .
task example
```

Run the example directly:

```bash
go run . example
go run . example --config ./gsea.example.toml --color=false
```

Create a starter config in the default XDG location:

```bash
go run . init-config
```

## Configuration

The default user config path is:

- Linux: `$XDG_CONFIG_HOME/gsea/config.toml` or `~/.config/gsea/config.toml`
- macOS: `~/Library/Application Support/gsea/config.toml`
- Windows: `%LOCALAPPDATA%\gsea\config.toml`

You can also point the CLI at an explicit config file:

```bash
go run . example --config ./gsea.example.toml
```

Example config structure:

```toml
[example]
greeting = "Ahoy"
name = "Harbor Team"
location = "Atlantic Ocean"
favorite_tide = "spring tide"

[output]
color = true
show_paths = true

[palette]
black = "#222222"
gray = "#69676c"
white = "#f7f1ff"
pink = "#FC618D"
orange = "#fd9353"
yellow = "#FCE566"
green = "#7BD88F"
cyan = "#5AD4E6"
purple = "#948ae3"
```

The palette values come from a reduced Monokai Spectrumish set with one black, one gray, one white, and six accent colors:

```text
https://github.com/borland502/nix-config/blob/main/home-manager/config/colors/monokai.base24.yaml
```

Environment variables override config values. Examples:

```bash
GSEA_EXAMPLE_NAME="Dock Crew" go run . example
GSEA_OUTPUT_SHOW_PATHS=false go run . example
GSEA_PALETTE_PURPLE="#A178FF" go run . example
```

## Commands

- `gsea example`: run the example program
- `gsea init-config`: write a starter config file

## Development

Common workflows are in `Taskfile.yml`:

```bash
task fmt
task check
task build
task example -- --config ./gsea.example.toml --color=false
```

Local install and removal tasks are also available:

```bash
task deploy
task undeploy
```

These tasks are intended for POSIX shells on Linux and macOS.

Both tasks accept overrides for the install and config locations when you need a safe or custom target:

```bash
task deploy INSTALL_DIR=./tmp/bin CONFIG_PATH=./tmp/gsea/config.toml
task undeploy INSTALL_DIR=./tmp/bin CONFIG_PATH=./tmp/gsea/config.toml
```

## Workflow Credit

The GitHub Actions workflows in [`.github/workflows`](https://github.com/borland502/go-sea/tree/main/.github/workflows) are adapted from the workflow layout used by [`charmbracelet/gum`](https://github.com/charmbracelet/gum/tree/main/.github/workflows), simplified for this repository.

## License

This project is licensed under the MIT License. See `LICENSE`.
