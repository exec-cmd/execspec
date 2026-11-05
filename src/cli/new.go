package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/exec-cmd/execspec/src/spec"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var newOpt spec.NewOptions

var newCmd = &cobra.Command{
	Use:   "new <type> <title>",
	Short: "Create a spec item from its template",
	Long: "Create a spec item from its template.\n\n" +
		"Types: feature, issue, update, fix, plan (numbered, with status), note, analysis.\n" +
		"An update or fix run in a terminal without --with asks which items to include.",
	Example: "  xspec new feature \"Login form\"\n" +
		"  xspec new update v0.2 --with 7,9\n" +
		"  xspec new fix typos --bug \"typo in --help\" --bug \"wrong exit code\"",
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := spec.TypeByName(args[0])
		if err != nil {
			return err
		}
		root, err := specRoot()
		if err != nil {
			return err
		}
		if (t.Name == "update" || t.Name == "fix") && !cmd.Flags().Changed("with") && isTTY() {
			if newOpt.With, err = pick(root, t.Name); err != nil {
				return err
			}
		}
		path, err := spec.New(root, t, strings.Join(args[1:], " "), newOpt)
		if err != nil {
			return err
		}
		fmt.Println(rel(path))
		return nil
	},
}

func isTTY() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// pick shows the items an update or fix may include and reads "1 3 4".
func pick(root, typ string) ([]int, error) {
	items, err := spec.Scan(root)
	if err != nil {
		return nil, err
	}
	cands := spec.Includable(items, typ)
	if len(cands) == 0 {
		return nil, nil
	}
	for i, it := range cands {
		fmt.Printf("%3d) %s %-7s %-5s %s\n", i+1, it.ID(), it.Type.Name, it.Status, it.Title)
	}
	fmt.Print("Include (numbers from the list, space-separated, empty for none): ")
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	var nums []int
	for _, f := range strings.Fields(strings.ReplaceAll(line, ",", " ")) {
		i, err := strconv.Atoi(f)
		if err != nil || i < 1 || i > len(cands) {
			return nil, fmt.Errorf("%q is not a number from the list", f)
		}
		nums = append(nums, cands[i-1].Num)
	}
	return nums, nil
}

func init() {
	newCmd.Flags().IntSliceVar(&newOpt.With, "with", nil, "update/fix: item numbers to include, e.g. 7,9")
	newCmd.Flags().StringArrayVar(&newOpt.Bugs, "bug", nil, "fix: a small bug, repeatable")
	newCmd.Flags().BoolVar(&newOpt.Dir, "dir", false, "fix: make a folder with design.md")

	rootCmd.AddCommand(newCmd)
}
