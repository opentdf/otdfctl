package cmd

import (
	"fmt"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var selectors []string

func dev_selectorsGen(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)

	subject := c.Flags.GetRequiredString("subject")

	flattened, err := handlers.FlattenSubjectContext(subject)
	if err != nil {
		cli.ExitWithError("Failed to parse subject context keys and values", err)
	}

	rows := [][]string{}
	for _, item := range flattened {
		rows = append(rows, []string{item.Key, fmt.Sprintf("%v", item.Value)})
	}

	t := cli.NewTabular(rows...)
	cli.PrintSuccessTable(cmd, "", t)
}

func dev_selectorsTest(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)

	subject := c.Flags.GetRequiredString("subject")
	selectors = c.Flags.GetStringSlice("selector", selectors, cli.FlagsStringSliceOptions{Min: 1})

	flattened, err := handlers.FlattenSubjectContext(subject)
	if err != nil {
		cli.ExitWithError("Failed to process subject context keys and values", err)
	}

	rows := [][]string{}
	for _, item := range flattened {
		for _, selector := range selectors {
			if selector == item.Key {
				rows = append(rows, []string{item.Key, fmt.Sprintf("%v", item.Value)})
			}
		}
	}

	t := cli.NewTabular(rows...)
	cli.PrintSuccessTable(cmd, "", t)
}

func init() {
	genCmd := man.Docs.GetCommand("dev/selectors/generate",
		man.WithRun(dev_selectorsGen),
	)
	genCmd.Flags().StringP(
		genCmd.GetDocFlag("subject").Name,
		genCmd.GetDocFlag("subject").Shorthand,
		genCmd.GetDocFlag("subject").Default,
		genCmd.GetDocFlag("subject").Description,
	)

	testCmd := man.Docs.GetCommand("dev/selectors/test",
		man.WithRun(dev_selectorsTest),
	)
	testCmd.Flags().StringP(
		testCmd.GetDocFlag("subject").Name,
		testCmd.GetDocFlag("subject").Shorthand,
		testCmd.GetDocFlag("subject").Default,
		testCmd.GetDocFlag("subject").Description,
	)
	testCmd.Flags().StringSliceVarP(
		&selectors,
		testCmd.GetDocFlag("selector").Name,
		testCmd.GetDocFlag("selector").Shorthand,
		[]string{},
		testCmd.GetDocFlag("selector").Description,
	)

	doc := man.Docs.GetCommand("dev/selectors",
		man.WithSubcommands(genCmd, testCmd),
	)

	dev_selectorsCmd := &doc.Command
	devCmd.AddCommand(dev_selectorsCmd)
}
