package cmd

import (
	"runtime"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage profiles (experimental)",
}

var profileCreateCmd = &cobra.Command{
	Use:     "create <profile> <endpoint>",
	Aliases: []string{"add"},
	Short:   "Create a new profile",
	//nolint:mnd // two args
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := cli.New(cmd, args)
		InitProfile(c, true)

		profileName := args[0]
		endpoint := args[1]

		setDefault := c.FlagHelper.GetOptionalBool("set-default")
		tlsNoVerify := c.FlagHelper.GetOptionalBool("tls-no-verify")

		c.Printf("Creating profile %s... ", profileName)
		if err := profile.AddProfile(profileName, endpoint, tlsNoVerify, setDefault); err != nil {
			c.Println("failed")
			c.ExitWithError("Failed to create profile", err)
		}
		c.Println("ok")

		// suggest the user to set up authentication
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List profiles",
	Run: func(cmd *cobra.Command, args []string) {
		c := cli.New(cmd, args)
		InitProfile(c, false)

		for _, p := range profile.GetGlobalConfig().ListProfiles() {
			if p == profile.GetGlobalConfig().GetDefaultProfile() {
				c.Printf("* %s\n", p)
				continue
			}
			c.Printf("  %s\n", p)
		}
	},
}

var profileGetCmd = &cobra.Command{
	Use:   "get <profile>",
	Short: "Get a profile value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := cli.New(cmd, args)
		InitProfile(c, false)

		profileName := args[0]
		p, err := profile.GetProfile(profileName)
		if err != nil {
			c.ExitWithError("Failed to load profile", err)
		}

		isDefault := "false"
		if p.GetProfileName() == profile.GetGlobalConfig().GetDefaultProfile() {
			isDefault = "true"
		}

		var auth string
		ac := p.GetAuthCredentials()
		if ac.AuthType == profiles.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS {
			maskedSecret := "********"
			auth = "client-credentials (" + ac.ClientId + ", " + maskedSecret + ")"
		}

		t := cli.NewTabular(
			[]string{"Profile", p.GetProfileName()},
			[]string{"Endpoint", p.GetEndpoint()},
			[]string{"Is default", isDefault},
			[]string{"Auth type", auth},
		)

		c.Print(t.View())
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete <profile>",
	Short: "Delete a profile",
	Run: func(cmd *cobra.Command, args []string) {
		c := cli.New(cmd, args)
		InitProfile(c, false)

		profileName := args[0]

		// TODO: suggest delete-all command to delete all profiles including default

		c.Printf("Deleting profile %s... ", profileName)
		if err := profile.DeleteProfile(profileName); err != nil {
			if err == profiles.ErrDeletingDefaultProfile {
				c.ExitWithWarning("Profile is set as default. Please set another profile as default before deleting.")
			}
			c.ExitWithError("Failed to delete profile", err)
		}
		c.Println("ok")
	},
}

// TODO add delete-all command

var profileSetDefaultCmd = &cobra.Command{
	Use:   "set-default <profile>",
	Short: "Set a profile as default",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := cli.New(cmd, args)
		InitProfile(c, false)

		profileName := args[0]

		c.Printf("Setting profile %s as default... ", profileName)
		if err := profile.SetDefaultProfile(profileName); err != nil {
			c.ExitWithError("Failed to set default profile", err)
		}
		c.Println("ok")
	},
}

var profileSetEndpointCmd = &cobra.Command{
	Use:   "set-endpoint <profile> <endpoint>",
	Short: "Set a profile value",
	//nolint:mnd // two args
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := cli.New(cmd, args)
		InitProfile(c, false)

		profileName := args[0]
		endpoint := args[1]

		p, err := profile.GetProfile(profileName)
		if err != nil {
			cli.ExitWithError("Failed to load profile", err)
		}

		c.Printf("Setting endpoint for profile %s... ", profileName)
		if err := p.SetEndpoint(endpoint); err != nil {
			c.ExitWithError("Failed to set endpoint", err)
		}
		c.Println("ok")
	},
}

func init() {
	// Profiles are not supported on Linux (unless mocked in test mode)
	if runtime.GOOS == "linux" && config.TestMode != "true" {
		return
	}

	profileCreateCmd.Flags().Bool("set-default", false, "Set the profile as default")
	profileCreateCmd.Flags().Bool("tls-no-verify", false, "Disable TLS verification")

	profileSetEndpointCmd.Flags().Bool("tls-no-verify", false, "Disable TLS verification")

	RootCmd.AddCommand(profileCmd)

	profileCmd.AddCommand(profileCreateCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileGetCmd)
	profileCmd.AddCommand(profileDeleteCmd)
	profileCmd.AddCommand(profileSetDefaultCmd)
	profileCmd.AddCommand(profileSetEndpointCmd)
}
