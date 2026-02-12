package auth

import (
	"fmt"
	"strings"

	"github.com/opentdf/otdfctl/cmd/common"
	"github.com/opentdf/otdfctl/pkg/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/spf13/cobra"
)

var clientCredScopes []string

func clientCredentialsRun(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	cp := common.InitProfile(c)

	var clientID string
	var clientSecret string

	if len(args) > 0 {
		clientID = args[0]
	}
	if len(args) > 1 {
		clientSecret = args[1]
	}

	if clientID == "" {
		clientID = cli.AskForInput("Enter client id: ")
	}
	if clientSecret == "" {
		clientSecret = cli.AskForSecret("Enter client secret: ")
	}

	// Prompt for scopes if not provided via flag
	if len(clientCredScopes) == 0 {
		scopeInput := cli.AskForInput("Enter scopes (comma-separated, or leave empty to skip): ")
		if scopeInput != "" {
			for _, s := range strings.Split(scopeInput, ",") {
				s = strings.TrimSpace(s)
				if s != "" {
					clientCredScopes = append(clientCredScopes, s)
				}
			}
		}
	}

	// Set the client credentials
	err := cp.SetAuthCredentials(profiles.AuthCredentials{
		AuthType:     profiles.AuthTypeClientCredentials,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       clientCredScopes,
	})
	if err != nil {
		c.ExitWithError("Failed to set client credentials", err)
	}

	// Validate the client credentials
	if err := auth.ValidateProfileAuthCredentials(cmd.Context(), cp); err != nil {
		c.ExitWithError("An error occurred during login. Please check your credentials and try again", err)
	}

	c.ExitWithMessage(fmt.Sprintf("Client credentials set for profile [%s]", cp.Name()), cli.ExitCodeSuccess)
}

// newClientCredentialsCmd creates and configures the client-credentials command.
func newClientCredentialsCmd() *cobra.Command {
	doc := man.Docs.GetCommand("auth/client-credentials",
		man.WithRun(clientCredentialsRun),
		man.WithHiddenFlags("with-client-creds", "with-client-creds-file"),
	)
	doc.Flags().StringSliceVar(&clientCredScopes, "scopes", []string{}, "Optional scopes to request during authentication (e.g. openid,email)")
	return &doc.Command
}
