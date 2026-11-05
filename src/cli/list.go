package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/exec-cmd/execspec/src/spec"
	"github.com/spf13/cobra"
)

var listAll bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List numbered items and their status",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := specRoot()
		if err != nil {
			return err
		}
		items, err := spec.Scan(root)
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NUM\tTYPE\tSTATUS\t[?]\tTASKS\tTITLE")
		for _, it := range items {
			if it.Archived && !listAll {
				continue
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d/%d\t%s\n", it.ID(), it.Type.Name, it.Status, len(it.Markers), it.Done, it.Total, it.Title)
		}
		return w.Flush()
	},
}

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "include archived items")

	rootCmd.AddCommand(listCmd)
}
