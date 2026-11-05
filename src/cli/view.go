package cli

import (
	"fmt"

	"github.com/exec-cmd/execspec/src/spec"
	"github.com/spf13/cobra"
)

// scanned runs fn on the items of the spec/ above cwd.
func scanned(fn func(root string, items []*spec.Item) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		root, err := specRoot()
		if err != nil {
			return err
		}
		items, err := spec.Scan(root)
		if err != nil {
			return err
		}
		return fn(root, items)
	}
}

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Print what to do next",
	Args:  cobra.NoArgs,
	RunE: scanned(func(root string, items []*spec.Item) error {
		fmt.Println(spec.Next(items))
		return nil
	}),
}

var questionsCmd = &cobra.Command{
	Use:   "questions",
	Short: "List every open [?: …] question",
	Args:  cobra.NoArgs,
	RunE: scanned(func(root string, items []*spec.Item) error {
		fmt.Print(spec.Questions(root, items))
		return nil
	}),
}

var roadmapCmd = &cobra.Command{
	Use:   "roadmap",
	Short: "Rewrite spec/ROADMAP.md and print it",
	Args:  cobra.NoArgs,
	RunE: scanned(func(root string, items []*spec.Item) error {
		fmt.Print(spec.Roadmap(items))
		return spec.WriteRoadmap(root)
	}),
}

var contextCmd = &cobra.Command{
	Use:   "context <num|slug>",
	Short: "Print an item with everything an agent needs to work on it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := specRoot()
		if err != nil {
			return err
		}
		out, err := spec.Context(root, args[0])
		fmt.Print(out)
		return err
	},
}

func init() {
	rootCmd.AddCommand(nextCmd, questionsCmd, roadmapCmd, contextCmd)
}
