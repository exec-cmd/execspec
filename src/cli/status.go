package cli

import (
	"fmt"

	"github.com/exec-cmd/execspec/src/spec"
	"github.com/spf13/cobra"
)

var change string

var startCmd = &cobra.Command{
	Use:   "start <num|slug>",
	Short: "Set an item, and everything it includes, to doing",
	Long:  "Set an item, and everything it includes, to doing. Refuses while any of them has [?: …] markers.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := specRoot()
		if err != nil {
			return err
		}
		g, err := spec.Start(root, args[0])
		for _, it := range g {
			fmt.Printf("%s doing  %s\n", it.ID(), it.Title)
		}
		return err
	},
}

var doneCmd = &cobra.Command{
	Use:   "done <num|slug>",
	Short: "Close an item, and everything it includes, and archive it",
	Long: "Close an item, and everything it includes: set done, write the Change line\n" +
		"(--change, or the current git HEAD) and move it to spec/archive/.\n" +
		"Refuses while a task box is unchecked or a [?: …] marker remains.",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := specRoot()
		if err != nil {
			return err
		}
		g, err := spec.Done(root, args[0], change)
		for _, it := range g {
			fmt.Printf("%s done   %s -> %s\n", it.ID(), it.Title, rel(it.Path))
		}
		return err
	},
}

func init() {
	doneCmd.Flags().StringVar(&change, "change", "", "what closed it: a commit, PR or path (default: git HEAD)")

	rootCmd.AddCommand(startCmd, doneCmd)
}
