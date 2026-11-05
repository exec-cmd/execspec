package cli

import (
	"fmt"
	"os"

	"github.com/exec-cmd/execspec/src/skill"
	"github.com/exec-cmd/execspec/src/spec"
	"github.com/spf13/cobra"
)

var force bool

var installCmd = &cobra.Command{
	Use:       "install project|global",
	Short:     "Install the execspec agent skill",
	Long:      "Install the execspec agent skill into .claude/skills/ and .agents/skills/\nof the project (git top-level) or of your home directory.",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"project", "global"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var base string
		var err error
		if args[0] == "global" {
			base, err = os.UserHomeDir()
		} else {
			base, err = os.Getwd()
			base = spec.ProjectRoot(base)
		}
		if err != nil {
			return err
		}
		dirs, err := skill.Install(base, force)
		for _, d := range dirs {
			fmt.Println("installed", rel(d))
		}
		return err
	},
}

func init() {
	installCmd.Flags().BoolVar(&force, "force", false, "replace an existing skill directory")

	rootCmd.AddCommand(installCmd)
}
