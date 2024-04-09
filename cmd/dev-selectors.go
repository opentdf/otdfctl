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

	dev_selectorsCmd = &cobra.Command{
		Use:   "selectors",
		Short: "Validate selectors",
		Long:  `Generation and validation of different selectors with jq selector syntax.`,
	}

	dev_selectorsGenCmd = &cobra.Command{
		Use:   "gen",
		Short: "Generate a set of selector expressions for keys and values of a Subject Context",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			subject := flagHelper.GetRequiredString("subject")
			contextType := flagHelper.GetRequiredString("type")

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

			result, err := handlers.ProcessSubjectContext(value, "", []*policy.SubjectProperty{})
			if err != nil {
				cli.ExitWithError("Failed to process subject context keys and values", err)
			}

			rows := [][]string{}
			for _, r := range result {
				rows = append(rows, []string{r.ExternalField, r.ExternalValue})
			}

			t := cli.NewTabular().Rows(rows...)
			cli.PrintSuccessTable(cmd, "", t)
		},
	}

	dev_selectorsTestCmd = &cobra.Command{
		Use:   "test",
		Short: "Test resolution of a set of selector expressions for keys and values of a Subject Context.",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			subject := flagHelper.GetRequiredString("subject")
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

			result, err := handlers.TestSubjectContext(value, selectors)
			if err != nil {
				cli.ExitWithError("Failed to process subject context keys and values", err)
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
	devCmd.AddCommand(dev_selectorsCmd)

	dev_selectorsCmd.AddCommand(dev_selectorsGenCmd)
	dev_selectorsGenCmd.Flags().StringP("subject", "s", "", "A Subject Context string (JSON or JWT, default JSON)")
	dev_selectorsGenCmd.Flags().StringP("type", "t", "json", "The type of the Subject Context: [json, jwt]")

	dev_selectorsCmd.AddCommand(dev_selectorsTestCmd)
	dev_selectorsTestCmd.Flags().StringP("subject", "s", "", "A Subject Context string (JSON or JWT, default JSON)")
	dev_selectorsTestCmd.Flags().StringP("type", "t", "json", "The type of the Subject Context: [json, jwt]")
	dev_selectorsTestCmd.Flags().StringArrayVarP(&selectors, "selector", "x", []string{}, "Individual selectors to test against the Subject Context (i.e. .key, .example[1].group)")
}
