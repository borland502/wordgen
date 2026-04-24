/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/adrg/xdg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const paletteGridColumns = 3

var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Run a colored demo of the configured palette, Cobra, Viper, and XDG setup",
	Run: func(cmd *cobra.Command, args []string) {
		color.NoColor = !cfg.Output.Color

		header := rgbColorFromHex(cfg.Palette.Purple, color.Bold)
		highlight := rgbColorFromHex(cfg.Palette.Green, color.Bold)
		cobraLabel := rgbColorFromHex(cfg.Palette.Yellow, color.Bold)
		viperLabel := rgbColorFromHex(cfg.Palette.Pink, color.Bold)
		colorLabel := rgbColorFromHex(cfg.Palette.Cyan, color.Bold)
		xdgLabel := rgbColorFromHex(cfg.Palette.Orange, color.Bold)
		paletteLabel := rgbColorFromHex(cfg.Palette.Purple, color.Bold)
		muted := rgbColorFromHex(cfg.Palette.Gray)

		header.Fprintf(cmd.OutOrStdout(), "%s example\n", binaryName)
		muted.Fprintf(cmd.OutOrStdout(), "using palette %s from config values\n\n", paletteName)

		highlight.Fprintf(
			cmd.OutOrStdout(),
			"%s, %s. Welcome to the %s with a %s.\n\n",
			cfg.Example.Greeting,
			cfg.Example.Name,
			cfg.Example.Location,
			cfg.Example.FavoriteTide,
		)

		cobraLabel.Fprintf(cmd.OutOrStdout(), "Cobra")
		fmt.Fprintf(cmd.OutOrStdout(), ": the %q command and flags like --name and --location come from Cobra.\n", cmd.Name())

		viperLabel.Fprintf(cmd.OutOrStdout(), "Viper")
		fmt.Fprintf(cmd.OutOrStdout(), ": values were resolved from %s.\n", configSourceDescription())

		colorLabel.Fprintf(cmd.OutOrStdout(), "fatih/color")
		fmt.Fprintf(cmd.OutOrStdout(), ": RGB output uses the palette values stored in the config file.\n")

		paletteLabel.Fprintf(cmd.OutOrStdout(), "Palette")
		fmt.Fprintf(cmd.OutOrStdout(), ": %s (compact curated subset)\n", paletteName)
		printPaletteGrid(cmd.OutOrStdout(), cfg.Palette.entries(), paletteGridColumns)

		if cfg.Output.ShowPaths {
			fmt.Fprintln(cmd.OutOrStdout())
			xdgLabel.Fprintf(cmd.OutOrStdout(), "XDG")
			fmt.Fprintf(cmd.OutOrStdout(), ": config home %s\n", xdg.ConfigHome)
			muted.Fprintf(cmd.OutOrStdout(), "      searched config dirs: %s\n", strings.Join(xdg.ConfigDirs, ", "))
			muted.Fprintf(cmd.OutOrStdout(), "      default user config: %s\n", defaultUserConfigFilePath())
		}

		fmt.Fprintln(cmd.OutOrStdout())
		muted.Fprintln(cmd.OutOrStdout(), "Try overrides with --name, GSEA_EXAMPLE_NAME, or GSEA_PALETTE_PURPLE.")
	},
}

func init() {
	rootCmd.AddCommand(exampleCmd)

	defaults := defaultConfig()
	exampleCmd.Flags().String("greeting", defaults.Example.Greeting, "example greeting")
	exampleCmd.Flags().String("name", defaults.Example.Name, "name to greet in the example output")
	exampleCmd.Flags().String("location", defaults.Example.Location, "location used by the example output")
	exampleCmd.Flags().String("favorite-tide", defaults.Example.FavoriteTide, "favorite tide used by the example output")
	exampleCmd.Flags().Bool("color", defaults.Output.Color, "enable colored output")
	exampleCmd.Flags().Bool("show-paths", defaults.Output.ShowPaths, "show resolved XDG paths")
}

func parseHexColor(hex string) (int, int, int, bool) {
	hex = strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(hex) != 6 {
		return 0, 0, 0, false
	}

	value, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}

	return int(value>>16) & 0xFF, int(value>>8) & 0xFF, int(value) & 0xFF, true
}

func rgbColorFromHex(hex string, attrs ...color.Attribute) *color.Color {
	r, g, b, ok := parseHexColor(hex)
	if !ok {
		return color.New(attrs...)
	}

	paletteColor := color.RGB(r, g, b)
	if len(attrs) > 0 {
		paletteColor = paletteColor.Add(attrs...)
	}

	return paletteColor
}

func bgRGBColorFromHex(hex string, attrs ...color.Attribute) *color.Color {
	r, g, b, ok := parseHexColor(hex)
	if !ok {
		return color.New(attrs...)
	}

	paletteColor := color.BgRGB(r, g, b)
	if len(attrs) > 0 {
		paletteColor = paletteColor.Add(attrs...)
	}

	return paletteColor
}

func paletteSwatch(hex string) string {
	if color.NoColor {
		return "----"
	}

	return bgRGBColorFromHex(hex).Sprint("    ")
}

func printPaletteGrid(writer io.Writer, entries []paletteEntry, columns int) {
	for index, entry := range entries {
		entryColor := rgbColorFromHex(entry.Hex, color.Bold)
		entryColor.Fprintf(writer, "  %-7s", entry.Key)
		fmt.Fprintf(writer, " %s %s", paletteSwatch(entry.Hex), entry.Hex)

		if (index+1)%columns == 0 || index == len(entries)-1 {
			fmt.Fprintln(writer)
			continue
		}

		fmt.Fprint(writer, "  ")
	}
}
