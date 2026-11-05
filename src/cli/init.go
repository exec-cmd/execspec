package cli

import (
	"fmt"
	"os"

	"github.com/exec-cmd/execspec/src/spec"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create spec/ and keep it out of git",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		created, err := spec.Init(spec.ProjectRoot(wd))
		for _, p := range created {
			fmt.Println("created", rel(p))
		}
		if err == nil && len(created) == 0 {
			fmt.Println("spec/ is already set up")
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
