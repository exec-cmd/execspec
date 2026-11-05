package cli

import (
	"fmt"

	"github.com/exec-cmd/execspec/src/spec"
	"github.com/spf13/cobra"
)

var doneOpt spec.DoneOptions

var approveCmd = &cobra.Command{
	Use:   "approve <num|slug>",
	Short: "Approve an item, and everything it includes (owner only)",
	Long: "Approve an item, and everything it includes, so it can start. Refuses while any of\n" +
		"them has [?: …] markers. The approval records a hash of the content: an edit\n" +
		"afterwards (ticking task boxes aside) needs a new approval. Agents must not run it.",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := specRoot()
		if err != nil {
			return err
		}
		g, err := spec.Approve(root, args[0])
		for _, it := range g {
			fmt.Printf("%s approved  %s\n", it.ID(), it.Title)
		}
		return err
	},
}

var startCmd = &cobra.Command{
	Use:   "start <num|slug>",
	Short: "Set an approved item, and everything it includes, to doing",
	Long:  "Set an item, and everything it includes, to doing. Each must be approved and unchanged since.",
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
	Long: "Close an item, and everything it includes: run their Verify commands, set done,\n" +
		"write the Change line (--change, or the current git HEAD, plus today's date)\n" +
		"and move them to spec/archive/.\n" +
		"Refuses while a task box is unchecked, a [?: …] marker remains, or a Verify command fails.",
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
	doneCmd.Flags().BoolVar(&doneOpt.NoVerify, "no-verify", false, "skip the Verify commands")

	rootCmd.AddCommand(approveCmd, startCmd, doneCmd)
}
