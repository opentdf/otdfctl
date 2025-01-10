package cmd

import (
	"errors"
	"runtime"

	"github.com/opentdf/otdfctl/internal/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/spf13/cobra"
)

var (
	runningInLinux    = runtime.GOOS == "linux"
	runningInTestMode = config.TestMode == "true"
)

var profileCmd = &cobra.Command{
	Use:     "profile",
	Aliases: []string{"p", "profiles"},
	Short:   "Manage profiles (experimental)",
	Hidden:  runningInLinux && !runningInTestMode,
}

var profileCreateCmd = &cobra.Command{
	Use:     "create <profile> <endpoint>",
	Aliases: []string{"add"},
	Short:   "Create a new profile",
	//nolint:mnd // two args
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := cli.New(cmd, args)
		profileMgr, _ := InitProfile(c, true)

		profileName := args[0]
		endpoint := args[1]

		setDefault := c.FlagHelper.GetOptionalBool("set-default")
		tlsNoVerify := c.FlagHelper.GetOptionalBool("tls-no-verify")

		c.Printf("Creating profile %s... ", profileName)
		profile := &profiles.ProfileCLI{
			Name:        profileName,
			Endpoint:    endpoint,
			TlsNoVerify: tlsNoVerify,
			// no credentials yet creating new profile pre-login
		}

		if err := profileMgr.AddProfile(profile, setDefault); err != nil {
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
		profileMgr, currProfile := InitProfile(c, false)

		profiles, err := profileMgr.ListProfiles()
		if err != nil {
			c.ExitWithError("Failed to list profiles", err)
		}

		for _, p := range profiles {
			if p.GetName() == currProfile.GetName() {
				c.Printf("* %s\n", p.GetName())
				continue
			}
			c.Printf("  %s\n", p.GetName())
		}
	},
}

var profileGetCmd = &cobra.Command{
	Use:   "get <profile>",
	Short: "Get a profile value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := cli.New(cmd, args)
		profileMgr, currProfile := InitProfile(c, false)

		profileName := args[0]
		p, err := profileMgr.GetProfile(profileName)
		if err != nil {
			c.ExitWithError("Failed to load profile", err)
		}

		isDefault := "false"
		if p.GetName() == currProfile.GetName() {
			isDefault = "true"
		}

		var authType string
		ac := p.GetAuthCredentials()
		if ac.AuthType == auth.AUTH_TYPE_CLIENT_CREDENTIALS {
			maskedSecret := "********"
			authType = "client-credentials (" + ac.ClientID + ", " + maskedSecret + ")"
		}

		t := cli.NewTabular(
			[]string{"Profile", p.GetName()},
			[]string{"Endpoint", p.GetEndpoint()},
			[]string{"Is default", isDefault},
			[]string{"Auth type", authType},
		)

		c.Print(t.View())
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete <profile>",
	Short: "Delete a profile",
	Run: func(cmd *cobra.Command, args []string) {
		c := cli.New(cmd, args)
		profileMgr, _ := InitProfile(c, false)

		profileName := args[0]

		// TODO: suggest delete-all command to delete all profiles including default

		c.Printf("Deleting profile %s... ", profileName)
		if err := profileMgr.DeleteProfile(profileName); err != nil {
			if errors.Is(err, profiles.ErrDeletingDefaultProfile) {
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
		profileMgr, _ := InitProfile(c, false)

		profileName := args[0]

		c.Printf("Setting profile %s as default... ", profileName)
		if err := profileMgr.SetDefaultProfile(profileName); err != nil {
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
		profileMgr, _ := InitProfile(c, false)

		profileName := args[0]
		endpoint := args[1]

		p, err := profileMgr.GetProfile(profileName)
		if err != nil {
			cli.ExitWithError("Failed to load profile", err)
		}
		if err := p.SetEndpoint(endpoint); err != nil {
			cli.ExitWithError("Failed to set endpoint", err)
		}
		if err := profileMgr.UpdateProfile(p); err != nil {
			c.Println("failed")
			c.ExitWithError("Failed to save profile", err)
		}
		c.Println("ok")
	},
}

func init() {
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
