package cmd

import (
	"runtime"

	"github.com/opentdf/otdfctl/pkg/cli"
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
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// ensure profile is initialized
		InitProfile(cmd, true)

		profileName := args[0]
		endpoint := args[1]

		fh := cli.NewFlagHelper(cmd)
		setDefault := fh.GetOptionalBool("set-default")
		tlsNoVerify := fh.GetOptionalBool("tls-no-verify")

		print := cli.NewPrinter(true)
		print.Printf("Creating profile %s... ", profileName)
		if err := profile.AddProfile(profileName, endpoint, tlsNoVerify, setDefault); err != nil {
			print.Println("failed")
			cli.ExitWithError("Failed to create profile", err)
		}
		print.Println("ok")

		// suggest the user to set up authentication
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List profiles",
	Run: func(cmd *cobra.Command, args []string) {
		// ensure profile is initialized
		InitProfile(cmd, false)

		print := cli.NewPrinter(true)
		for _, p := range profile.GetGlobalConfig().ListProfiles() {
			if p == profile.GetGlobalConfig().GetDefaultProfile() {
				print.Printf("* %s\n", p)
				continue
			}
			print.Printf("  %s\n", p)
		}
	},
}

var profileGetCmd = &cobra.Command{
	Use:   "get <profile>",
	Short: "Get a profile value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// ensure profile is initialized
		InitProfile(cmd, false)

		profileName := args[0]
		p, err := profile.GetProfile(profileName)
		if err != nil {
			cli.ExitWithError("Failed to load profile", err)
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

		print := cli.NewPrinter(true)
		print.Print(t.View())
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete <profile>",
	Short: "Delete a profile",
	Run: func(cmd *cobra.Command, args []string) {
		// ensure profile is initialized
		InitProfile(cmd, false)

		profileName := args[0]

		// TODO check if the profile is the default and prevent
		// suggest delete-all command to delete all profiles including default

		print := cli.NewPrinter(true)
		print.Printf("Deleting profile %s... ", profileName)
		if err := profile.DeleteProfile(profileName); err != nil {
			cli.ExitWithError("Failed to delete profile", err)
		}
		print.Println("ok")
	},
}

// TODO add delete-all command

var profileSetDefaultCmd = &cobra.Command{
	Use:   "set-default <profile>",
	Short: "Set a profile as default",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// ensure profile is initialized
		InitProfile(cmd, false)

		profileName := args[0]

		print := cli.NewPrinter(true)
		print.Printf("Setting profile %s as default... ", profileName)
		if err := profile.SetDefaultProfile(profileName); err != nil {
			cli.ExitWithError("Failed to set default profile", err)
		}
		print.Println("ok")
	},
}

var profileSetEndpointCmd = &cobra.Command{
	Use:   "set-endpoint <profile> <endpoint>",
	Short: "Set a profile value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// ensure profile is initialized
		InitProfile(cmd, false)

		profileName := args[0]
		endpoint := args[1]

		p, err := profile.GetProfile(profileName)
		if err != nil {
			cli.ExitWithError("Failed to load profile", err)
		}

		print := cli.NewPrinter(true)
		print.Printf("Setting endpoint for profile %s... ", profileName)
		if err := p.SetEndpoint(endpoint); err != nil {
			cli.ExitWithError("Failed to set endpoint", err)
		}
		print.Println("ok")
	},
}

func init() {
	// Profiles are not supported on Linux
	if runtime.GOOS == "linux" {
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
