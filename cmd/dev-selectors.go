package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
)

var selectors []string

func dev_selectorsGen(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	subject := c.Flags.GetRequiredString("subject")
	contextType := c.Flags.GetRequiredString("type")

	var value any
	//nolint:gocritic,nestif // this is more readable than a switch statement
	if contextType == "json" || contextType == "" {
		if err := json.Unmarshal([]byte(subject), &value); err != nil {
			cli.ExitWithError(fmt.Sprintf("Could not unmarshal JSON subject context input: %s", subject), err)
		}
	} else if contextType == "jwt" {
		// get the payload from the decoded JWT
		token, _, err := new(jwt.Parser).ParseUnverified(subject, jwt.MapClaims{})
		if err != nil {
			cli.ExitWithError("Failed to parse JWT token", err)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			value = claims
		} else {
			cli.ExitWithError("Failed to get claims from JWT token", nil)
		}
	} else {
		cli.ExitWithError("Invalid subject context type. Must be of type: [json, jwt]", nil)
	}

	result, err := handlers.ProcessSubjectContext(value, "", []*policy.SubjectProperty{})
	if err != nil {
		cli.ExitWithError("Failed to process subject context keys and values", err)
	}

	rows := [][]string{}
	for _, r := range result {
		rows = append(rows, []string{r.GetExternalSelectorValue(), r.GetExternalValue()})
	}

	t := cli.NewTabular(rows...)
	cli.PrintSuccessTable(cmd, "", t)
}

func dev_selectorsTest(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	subject := c.Flags.GetRequiredString("subject")
	contextType := c.Flags.GetRequiredString("type")
	selectors = c.Flags.GetStringSlice("selectors", selectors, cli.FlagsStringSliceOptions{Min: 1})

	var value any
	//nolint:gocritic,nestif // this is more readable than a switch statement
	if contextType == "json" || contextType == "" {
		if err := json.Unmarshal([]byte(subject), &value); err != nil {
			cli.ExitWithError(fmt.Sprintf("Could not unmarshal JSON subject context input: %s", subject), err)
		}
	} else if contextType == "jwt" {
		token, _, err := new(jwt.Parser).ParseUnverified(subject, jwt.MapClaims{})
		if err != nil {
			cli.ExitWithError("Failed to parse JWT token", err)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			value = claims
		} else {
			cli.ExitWithError("Failed to get claims from JWT token", nil)
		}
	} else {
		cli.ExitWithError("Invalid subject context type. Must be of type: [json, jwt]", nil)
	}

	result, err := handlers.TestSubjectContext(value, selectors)
	if err != nil {
		cli.ExitWithError("Failed to process subject context keys and values", err)
	}

	rows := [][]string{}
	for _, r := range result {
		rows = append(rows, []string{r.GetExternalSelectorValue(), r.GetExternalValue()})
	}

	t := cli.NewTabular(rows...)
	cli.PrintSuccessTable(cmd, "", t)
}

func init() {
	genCmd := man.Docs.GetCommand("dev/selectors/gen",
		man.WithRun(dev_selectorsGen),
	)
	genCmd.Flags().StringP(
		genCmd.GetDocFlag("subject").Name,
		genCmd.GetDocFlag("subject").Shorthand,
		genCmd.GetDocFlag("subject").Default,
		genCmd.GetDocFlag("subject").Description,
	)
	genCmd.Flags().StringP(
		genCmd.GetDocFlag("type").Name,
		genCmd.GetDocFlag("type").Shorthand,
		genCmd.GetDocFlag("type").Default,
		genCmd.GetDocFlag("type").Description,
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
	testCmd.Flags().StringP(
		testCmd.GetDocFlag("type").Name,
		testCmd.GetDocFlag("type").Shorthand,
		testCmd.GetDocFlag("type").Default,
		testCmd.GetDocFlag("type").Description,
	)
	testCmd.Flags().StringArrayP(
		testCmd.GetDocFlag("selector").Name,
		testCmd.GetDocFlag("selector").Shorthand,
		[]string{},
		testCmd.GetDocFlag("selector").Description,
	)

	// TODO: put back dev selectors command once the flattening lib is provided by platform
	// issue: https://github.com/opentdf/otdfctl/issues/125

	// doc := man.Docs.GetCommand("dev/selectors",
	// 	man.WithSubcommands(genCmd, testCmd),
	// )

	// dev_selectorsCmd = &doc.Command
	// devCmd.AddCommand(dev_selectorsCmd)
}
