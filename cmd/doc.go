// Package cmd wires the gsea command tree, configuration loading, and starter
// commands.
//
// The current shape of the package is intentionally demo-oriented. The
// example command exists to exercise Cobra, Viper, XDG, and terminal color
// output in one place. If the CLI takes on a concrete product objective, treat
// this package as the integration layer and replace the sample pieces instead
// of layering more demo behavior on top of them.
//
// Typical follow-on changes are:
//
//   - remove exampleCmd and its flag bindings once real subcommands exist
//   - rename or delete ExampleConfig and PaletteConfig when production commands
//     no longer need sample greeting or palette data
//   - update defaultConfig, renderConfig, and init-config so generated config
//     only contains sections consumed by the real command set
//
// Execute remains the package entrypoint used by main.
package cmd
