package cli

import (
	"fmt"

	"github.com/exec-cmd/execspec/src/spec"
	"github.com/spf13/cobra"
)

var doneOpt spec.DoneOptions

var doneCmd = &cobra.Command{
	Use:   "done <num|slug>",
	Short: "Close an item, and everything it includes, and archive it",
	Long: "Close an item, and everything it includes: run their Verify commands, set done,\n" +
		"write the Change line (--change, or the current git HEAD, plus today's date)\n" +
		"and move them to spec/archive/.\n" +
		"Refuses while a task box is unchecked, a [?: …] marker remains, or a Verify command fails.\n" +
		"--decision adds one line (what was chosen and why) to spec/notes/decisions.md.",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := specRoot()
		if err != nil {
			return err
		}
		if doneOpt.NoVerify {
			fmt.Println("skipping Verify commands (--no-verify)")
		}
		doneOpt.Run = func(c string) error {
			fmt.Println("verify:", c)
			return spec.Shell(root)(c)
		}
		g, err := spec.Done(root, args[0], doneOpt)
		for _, it := range g {
			fmt.Printf("%s done   %s -> %s\n", it.ID(), it.Title, rel(it.Path))
		}
		return err
	},
}

func init() {
	doneCmd.Flags().StringVar(&doneOpt.Change, "change", "", "what closed it: a commit, PR or path (default: git HEAD)")
	doneCmd.Flags().StringVar(&doneOpt.Decision, "decision", "", "one line for notes/decisions.md: what was chosen and why")
	doneCmd.Flags().BoolVar(&doneOpt.NoVerify, "no-verify", false, "skip the Verify commands")

	rootCmd.AddCommand(doneCmd)
}
