package cmd

import (
	"errors"
	"runtime"

	generated "github.com/opentdf/otdfctl/cmd/generated"
	profilegen "github.com/opentdf/otdfctl/cmd/generated/profile"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/spf13/cobra"
)

var (
	runningInLinux    = runtime.GOOS == "linux"
	runningInTestMode = config.TestMode == "true"
)

func init() {
	// Create commands using generated constructors with handler functions
	profileCmd := generated.NewProfileCommand(handleProfile)
	createCmd := profilegen.NewCreateCommand(handleProfileCreate)
	listCmd := profilegen.NewListCommand(handleProfileList)
	getCmd := profilegen.NewGetCommand(handleProfileGet)
	deleteCmd := profilegen.NewDeleteCommand(handleProfileDelete)
	setDefaultCmd := profilegen.NewSetDefaultCommand(handleProfileSetDefault)
	setEndpointCmd := profilegen.NewSetEndpointCommand(handleProfileSetEndpoint)

	// Set Linux/test mode visibility
	profileCmd.Hidden = runningInLinux && !runningInTestMode

	// Add subcommands
	profileCmd.AddCommand(createCmd)
	profileCmd.AddCommand(listCmd)
	profileCmd.AddCommand(getCmd)
	profileCmd.AddCommand(deleteCmd)
	profileCmd.AddCommand(setDefaultCmd)
	profileCmd.AddCommand(setEndpointCmd)

	// Add to root command
	RootCmd.AddCommand(profileCmd)
}

// handleProfile implements the parent profile command (if called without subcommands)
func handleProfile(cmd *cobra.Command, req *generated.ProfileRequest) error {
	return cmd.Help()
}

// handleProfileCreate implements the business logic for the create command
func handleProfileCreate(cmd *cobra.Command, req *profilegen.CreateRequest) error {
	c := cli.New(cmd, []string{req.Arguments.Profile, req.Arguments.Endpoint})
	InitProfile(c, true)

	profileName := req.Arguments.Profile
	endpoint := req.Arguments.Endpoint

	c.Printf("Creating profile %s... ", profileName)
	if err := profile.AddProfile(profileName, endpoint, req.Flags.TlsNoVerify, req.Flags.SetDefault); err != nil {
		c.Println("failed")
		c.ExitWithError("Failed to create profile", err)
	}
	c.Println("ok")

	return nil
}

// handleProfileList implements the business logic for the list command
func handleProfileList(cmd *cobra.Command, req *profilegen.ListRequest) error {
	c := cli.New(cmd, []string{})
	InitProfile(c, false)

	for _, p := range profile.GetGlobalConfig().ListProfiles() {
		if p == profile.GetGlobalConfig().GetDefaultProfile() {
			c.Printf("* %s\n", p)
			continue
		}
		c.Printf("  %s\n", p)
	}

	return nil
}

// handleProfileGet implements the business logic for the get command
func handleProfileGet(cmd *cobra.Command, req *profilegen.GetRequest) error {
	c := cli.New(cmd, []string{req.Arguments.Profile})
	InitProfile(c, false)

	profileName := req.Arguments.Profile
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
	return nil
}

// handleProfileDelete implements the business logic for the delete command
func handleProfileDelete(cmd *cobra.Command, req *profilegen.DeleteRequest) error {
	c := cli.New(cmd, []string{req.Arguments.Profile})
	InitProfile(c, false)

	profileName := req.Arguments.Profile

	c.Printf("Deleting profile %s... ", profileName)
	if err := profile.DeleteProfile(profileName); err != nil {
		if errors.Is(err, profiles.ErrDeletingDefaultProfile) {
			c.ExitWithWarning("Profile is set as default. Please set another profile as default before deleting.")
		}
		c.ExitWithError("Failed to delete profile", err)
	}
	c.Println("ok")

	return nil
}

// handleProfileSetDefault implements the business logic for the set-default command
func handleProfileSetDefault(cmd *cobra.Command, req *profilegen.SetDefaultRequest) error {
	c := cli.New(cmd, []string{req.Arguments.Profile})
	InitProfile(c, false)

	profileName := req.Arguments.Profile

	c.Printf("Setting profile %s as default... ", profileName)
	if err := profile.SetDefaultProfile(profileName); err != nil {
		c.ExitWithError("Failed to set default profile", err)
	}
	c.Println("ok")

	return nil
}

// handleProfileSetEndpoint implements the business logic for the set-endpoint command
func handleProfileSetEndpoint(cmd *cobra.Command, req *profilegen.SetEndpointRequest) error {
	c := cli.New(cmd, []string{req.Arguments.Profile, req.Arguments.Endpoint})
	InitProfile(c, false)

	profileName := req.Arguments.Profile
	endpoint := req.Arguments.Endpoint

	p, err := profile.GetProfile(profileName)
	if err != nil {
		cli.ExitWithError("Failed to load profile", err)
	}

	c.Printf("Setting endpoint for profile %s... ", profileName)
	if err := p.SetEndpoint(endpoint); err != nil {
		c.ExitWithError("Failed to set endpoint", err)
	}
	c.Println("ok")

	return nil
}