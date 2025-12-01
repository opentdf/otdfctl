package cmd

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	osprofiles "github.com/jrschumacher/go-osprofiles"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	runningInLinux       = runtime.GOOS == "linux"
	runningInTestMode    = config.TestMode == "true"
	errCleaningUpKeyring = errors.New("error occurred when cleaning up keyring")
)

func newProfilerFromCLI(c *cli.Cli) *osprofiles.Profiler {
	driverType := getDriverTypeFromUser(c)
	profiler, err := profiles.NewProfiler(string(driverType))
	if err != nil {
		cli.ExitWithError("Error creating profiler", err)
	}

	return profiler
}

func getDriverTypeFromUser(c *cli.Cli) profiles.ProfileDriver {
	driverTypeStr := string(profiles.PROFILE_DRIVER_DEFAULT)
	store := c.FlagHelper.GetOptionalString("store")
	if len(store) > 0 {
		driverTypeStr = store
	}

	driverType, err := profiles.ToProfileDriver(driverTypeStr)
	if err != nil {
		cli.ExitWithError("Error converting store type", err)
	}

	return driverType
}

var profileCmd = &cobra.Command{
	Use:     "profile",
	Aliases: []string{"profiles", "prof"},
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
		profileName := args[0]
		endpoint := args[1]

		setDefault := c.FlagHelper.GetOptionalBool("set-default")
		tlsNoVerify := c.FlagHelper.GetOptionalBool("tls-no-verify")
		driverType := getDriverTypeFromUser(c)

		c.Printf("Creating profile %s...", profileName)
		_, err := profiles.NewOtdfctlProfileStore(driverType, profileName, endpoint, tlsNoVerify, setDefault)
		if err != nil {
			c.Println("failed")
			c.ExitWithError("Failed to create profile", err)
		}
		c.Printf("ok")
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List profiles",
	Run: func(cmd *cobra.Command, args []string) {
		c := cli.New(cmd, args)
		driverType := getDriverTypeFromUser(c)
		profiler := newProfilerFromCLI(c)

		globalCfg := osprofiles.GetGlobalConfig(profiler)
		defaultProfile := globalCfg.GetDefaultProfile()

		c.Printf("Listing profiles from %s\n", string(driverType))
		for _, p := range osprofiles.ListProfiles(profiler) {
			if p == defaultProfile {
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
		profileName := args[0]

		driverType := getDriverTypeFromUser(c)
		profileStore, err := profiles.LoadOtdfctlProfileStore(driverType, profileName)
		if err != nil {
			cli.ExitWithError(fmt.Sprintf("Error loading profile store for profile %s", profileName), err)
		}

		isDefault := "false"
		if profileStore.IsDefault() {
			isDefault = "true"
		}

		var auth string
		ac := profileStore.GetAuthCredentials()
		if ac.AuthType == profiles.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS {
			maskedSecret := "********"
			auth = "client-credentials (" + ac.ClientId + ", " + maskedSecret + ")"
		}

		t := cli.NewTabular(
			[]string{"Profile", profileStore.Name()},
			[]string{"Endpoint", profileStore.GetEndpoint()},
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
		profileName := args[0]

		// TODO: suggest delete-all command to delete all profiles including default
		driverType := getDriverTypeFromUser(c)
		profiler := newProfilerFromCLI(c)

		c.Printf("Deleting profile %s, from %s...", profileName, driverType)
		if err := osprofiles.DeleteProfile[*profiles.ProfileConfig](profiler, profileName); err != nil {
			if strings.Contains(err.Error(), "cannot delete the default profile") {
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
		profileName := args[0]
		profiler := newProfilerFromCLI(c)

		c.Printf("Setting profile %s as default...", profileName)
		if err := osprofiles.SetDefaultProfile(profiler, profileName); err != nil {
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
		profileName := args[0]
		endpoint := args[1]
		profiler := newProfilerFromCLI(c)

		store, err := osprofiles.GetProfile[*profiles.ProfileConfig](profiler, profileName)
		if err != nil {
			cli.ExitWithError("Failed to load profile", err)
		}

		p, ok := store.Profile.(*profiles.ProfileConfig)
		if !ok || p == nil {
			cli.ExitWithError("Failed to load profile", errors.New("invalid profile configuration"))
		}

		u, err := utils.NormalizeEndpoint(endpoint)
		if err != nil {
			c.ExitWithError("Failed to set endpoint", err)
		}

		c.Printf("Setting endpoint for profile %s... ", profileName)
		p.Endpoint = u.String()
		if err := store.Save(); err != nil {
			c.ExitWithError("Failed to set endpoint", err)
		}
		c.Println("ok")
	},
}

func migrateKeyringProfilesToFilesystem(c *cli.Cli) {
	keyringProfiler, err := profiles.NewProfiler(string(profiles.PROFILE_DRIVER_KEYRING))
	if err != nil {
		c.ExitWithError("Failed to initialize keyring profile store", err)
	}

	filesystemProfiler := newProfilerFromCLI(c)

	profilesInKeyring := osprofiles.ListProfiles(keyringProfiler)
	if len(profilesInKeyring) == 0 {
		c.Println("No profiles found in keyring store to migrate.")
		return
	}

	defaultKeyringProfile := osprofiles.GetGlobalConfig(keyringProfiler).GetDefaultProfile()

	c.Printf("Migrating %d profiles from keyring to filesystem...\n", len(profilesInKeyring))

	for _, profileName := range profilesInKeyring {
		store, err := osprofiles.GetProfile[*profiles.ProfileConfig](keyringProfiler, profileName)
		if err != nil {
			c.ExitWithError("Failed to load profile from keyring", err)
		}

		p, ok := store.Profile.(*profiles.ProfileConfig)
		if !ok || p == nil {
			c.ExitWithError("Failed to load profile from keyring", errors.New("invalid profile configuration"))
		}

		setDefault := profileName == defaultKeyringProfile

		if err := filesystemProfiler.AddProfile(p, setDefault); err != nil {
			c.ExitWithError("Failed to migrate profile", err)
		}

		c.Printf("Migrated profile %s, set to default: %t\n", profileName, setDefault)
	}

	c.Printf("Removing %d profiles from the keyring\n", len(profilesInKeyring))
	if err = keyringProfiler.Cleanup(false); err != nil {
		cli.ExitWithError(errCleaningUpKeyring.Error(), err)
	}

	c.Println("Migration complete.")
}

var profileMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate all profiles from keyring to filesystem",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := cli.New(cmd, args)
		migrateKeyringProfilesToFilesystem(c)
	},
}

var profileKeyringCleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove all profiles and configuration from the keyring store. Use when migration fails.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := cli.New(cmd, args)

		force := c.Flags.GetOptionalBool("force")
		cli.ConfirmAction(cli.ActionDelete, "all profiles and configuration stored in the keyring", config.AppName, force)

		keyringProfiler, err := osprofiles.New(config.AppName, osprofiles.WithKeyringStore())
		if err != nil {
			c.ExitWithError("Failed to initialize keyring profile store", err)
		}

		c.Println("Cleaning up keyring profile store...")
		if err := keyringProfiler.Cleanup(false); err != nil {
			cli.ExitWithError(errCleaningUpKeyring.Error(), err)
		}
		c.Println("Keyring profile store cleanup complete.")
	},
}

func init() {
	profileCreateCmd.Flags().Bool("set-default", false, "Set the profile as default")
	profileCreateCmd.Flags().Bool("tls-no-verify", false, "Disable TLS verification")
	profileCreateCmd.Flags().String("store", "filesystem", "Profile store to use: filesystem or keyring")

	profileListCmd.Flags().String("store", "filesystem", "Profile store to use: filesystem or keyring")
	profileGetCmd.Flags().String("store", "filesystem", "Profile store to use: filesystem or keyring")
	profileDeleteCmd.Flags().String("store", "filesystem", "Profile store to use: filesystem or keyring")

	profileSetEndpointCmd.Flags().Bool("tls-no-verify", false, "Disable TLS verification")

	RootCmd.AddCommand(profileCmd)

	profileCmd.AddCommand(profileCreateCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileGetCmd)
	profileCmd.AddCommand(profileDeleteCmd)
	profileCmd.AddCommand(profileSetDefaultCmd)
	profileCmd.AddCommand(profileSetEndpointCmd)
	profileCmd.AddCommand(profileMigrateCmd)
	profileCmd.AddCommand(profileKeyringCleanupCmd)

	profileKeyringCleanupCmd.Flags().Bool("force", false, "Skip confirmation prompt")
}
