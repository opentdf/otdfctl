package auth

import (
	"fmt"

	"github.com/opentdf/otdfctl/cmd/common"
	"github.com/opentdf/otdfctl/pkg/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/spf13/cobra"
)

func printAccessTokenRun(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	cp := common.InitProfile(c)

	ac := cp.GetAuthCredentials()
	switch ac.AuthType {
	case profiles.AuthTypeClientCredentials:
	case profiles.AuthTypeAccessToken:
	default:
		c.ExitWithError("Invalid auth type", nil)
	}
	tok, err := auth.GetTokenWithProfile(cmd.Context(), cp)
	if err != nil {
		cli.ExitWithError("Failed to get token", err)
	}

	c.ExitWithStyled(fmt.Sprintf("Access Token: %s\n", tok.AccessToken), cli.ExitCodeSuccess)
	c.ExitWithJSON(tok, cli.ExitCodeSuccess)
}

// newPrintAccessTokenCmd creates and configures the print-access-token command.
func newPrintAccessTokenCmd() *cobra.Command {
	doc := man.Docs.GetCommand("auth/print-access-token", man.WithRun(printAccessTokenRun))
	return &doc.Command
}
