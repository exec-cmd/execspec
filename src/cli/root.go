package cli

import (
	"os"
	"path/filepath"

	"github.com/exec-cmd/execspec/src/spec"
	"github.com/spf13/cobra"
)

// version is set at build time: -ldflags "-X github.com/exec-cmd/execspec/src/cli.version=v0.3.0".
var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "xspec",
	Short: "Spec-driven development in a private spec/ directory",
	Long: "xspec keeps intent in a private, git-ignored spec/ directory: numbered changes\n" +
		"and updates that people and AI agents share as plain Markdown.",
	SilenceUsage: true,
	Version:      version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// specRoot finds the directory holding spec/, starting from cwd.
func specRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return spec.FindRoot(wd)
}

// rel shortens a path to be relative to cwd when that is shorter.
func rel(p string) string {
	wd, err := os.Getwd()
	if err != nil {
		return p
	}
	if r, err := filepath.Rel(wd, p); err == nil && len(r) < len(p) {
		return r
	}
	return p
}
