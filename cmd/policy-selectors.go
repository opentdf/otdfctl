package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
)

var (
	selectors []string

	policy_selectorsCmd = &cobra.Command{
		Use:   "selectors",
		Short: "Validate selectors",
		Long:  `Test validation of different selectors with different syntax.`,
	}

	policy_selectorsGetCmd = &cobra.Command{
		Use:   "find",
		Short: "Find a set of selected keys and values for a Subject Context",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			subject := flagHelper.GetRequiredString("subject")
			contextType := flagHelper.GetRequiredString("type")
			syntax := flagHelper.GetRequiredString("syntax")

			var value any
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

			result := []*policy.SubjectProperty{}

			if syntax == "jq" {
				r, err := handlers.ProcessSubjectContext(value, "", []*policy.SubjectProperty{})
				if err != nil {
					cli.ExitWithError("Failed to process subject context keys and values", err)
				}
				result = r
			} else {
				cli.ExitWithError("Invalid syntax. Must be of type: [jq]", nil)
			}

			rows := [][]string{}
			for _, r := range result {
				rows = append(rows, []string{r.ExternalField, r.ExternalValue})
			}

			t := cli.NewTabular().Rows(rows...)
			// TODO: we need to escape the JSON so don't use handlesuccess here
			cli.PrintSuccessTable(cmd, "", t)
		},
	}

	policy_selectorsTestCmd = &cobra.Command{
		Use:   "test",
		Short: "Test a set of selectors on a Subject Context",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			subject := flagHelper.GetRequiredString("subject")
			syntax := flagHelper.GetRequiredString("syntax")
			contextType := flagHelper.GetRequiredString("type")
			selectors := flagHelper.GetStringSlice("selectors", selectors, cli.FlagHelperStringSliceOptions{Min: 1})

			var value any
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

			result := []*policy.SubjectProperty{}

			if syntax == "jq" {
				r, err := handlers.TestSubjectContext(value, "jq", selectors)
				if err != nil {
					cli.ExitWithError("Failed to process subject context keys and values", err)
				}
				result = r
			} else {
				cli.ExitWithError("Invalid syntax. Must be of type: [jq]", nil)
			}

			rows := [][]string{}
			for _, r := range result {
				rows = append(rows, []string{r.ExternalField, r.ExternalValue})
			}

			t := cli.NewTabular().Rows(rows...)
			cli.PrintSuccessTable(cmd, "", t)
		},
	}
)

func init() {
	policyCmd.AddCommand(policy_selectorsCmd)

	policy_selectorsCmd.AddCommand(policy_selectorsGetCmd)
	policy_selectorsGetCmd.Flags().StringP("subject", "s", "", "A Subject Context string (JSON or JWT, default JSON)")
	policy_selectorsGetCmd.Flags().StringP("type", "t", "json", "The type of the Subject Context: [json, jwt]")
	policy_selectorsGetCmd.Flags().StringP("syntax", "y", "", "Syntax: [jq]")

	policy_selectorsCmd.AddCommand(policy_selectorsTestCmd)
	policy_selectorsTestCmd.Flags().StringP("subject", "s", "", "A Subject Context string (JSON or JWT, default JSON)")
	policy_selectorsTestCmd.Flags().StringP("type", "t", "json", "The type of the Subject Context: [json, jwt]")
	policy_selectorsTestCmd.Flags().StringP("syntax", "y", "", "Syntax: [jq]")
	policy_selectorsTestCmd.Flags().StringArrayVarP(&selectors, "selector", "x", []string{}, "Individual selectors to test against the Subject Context (i.e. .key, .example[1].group)")
}
