package cli

import (
	"fmt"

	"github.com/exec-cmd/execspec/src/spec"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Find problems in spec/; exits 1 when there are any",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := specRoot()
		if err != nil {
			return err
		}
		problems, err := spec.Check(root)
		if err != nil {
			return err
		}
		for _, p := range problems {
			fmt.Println(p)
		}
		if len(problems) > 0 {
			return fmt.Errorf("%d problem(s)", len(problems))
		}
		fmt.Println("spec/ is fine")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
