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

var profileInitCmd = &cobra.Command{
	Use:   "init <profile-name> <endpoint>",
	Short: "Initialize profile",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]
		endpoint := args[1]

		print := cli.NewPrinter(true)
		print.Println("Initializing profile...")
		if err := profile.AddProfile(profileName, endpoint, true); err != nil {
			cli.ExitWithError("Failed to initialize profile", err)
		}
		print.Println("ok")
	},
}

var profileCreateCmd = &cobra.Command{
	Use:   "create <profile> <endpoint>",
	Short: "Create a new profile",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]
		endpoint := args[1]

		fh := cli.NewFlagHelper(cmd)
		setDefault := fh.GetOptionalBool("set-default")

		print := cli.NewPrinter(true)
		print.Printf("Creating profile %s... ", profileName)
		if err := profile.AddProfile(profileName, endpoint, setDefault); err != nil {
			print.Println("failed")
			cli.ExitWithError("Failed to create profile", err)
		}
		print.Println("ok")
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List profiles",
	Run: func(cmd *cobra.Command, args []string) {
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
		profileName := args[0]

		print := cli.NewPrinter(true)
		print.Printf("Deleting profile %s... ", profileName)
		if err := profile.DeleteProfile(profileName); err != nil {
			cli.ExitWithError("Failed to delete profile", err)
		}
		print.Println("ok")
	},
}

var profileSetDefaultCmd = &cobra.Command{
	Use:   "set-default <profile>",
	Short: "Set a profile as default",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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
	if runtime.GOOS == "linux" {
		return
	}
	var err error
	profile, err = profiles.New()
	if err != nil {
		cli.ExitWithError("Failed to initialize profile store", err)
	}

	RootCmd.AddCommand(profileCmd)

	if profile.GetGlobalConfig().GetDefaultProfile() == "" {
		profileCmd.AddCommand(profileInitCmd)
	} else {
		profileCmd.AddCommand(profileCreateCmd)
		profileCmd.AddCommand(profileListCmd)
		profileCmd.AddCommand(profileGetCmd)
		profileCmd.AddCommand(profileDeleteCmd)
		profileCmd.AddCommand(profileSetDefaultCmd)
		profileCmd.AddCommand(profileSetEndpointCmd)
	}
}
