package cmd

import (
	"fmt"
	"os"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/spf13/cobra"
)

var (
	policy_certificatesCmd = &cobra.Command{
		Use:   "certificates",
		Short: "Manage certificates for namespaces",
		Long:  "Assign or remove root certificates from attribute namespaces. Use 'otdfctl policy attributes namespaces get' to view certificates.",
	}

	policy_certificatesConvertCmd = &cobra.Command{
		Use:   "convert-pem",
		Short: "Convert x5c certificate to PEM format",
		Long:  "Convert an x5c format certificate (base64-encoded DER) to PEM-encoded certificate",
		Run:   policy_convertPEMToX5C,
	}

	policy_certificatesAssignCmd = &cobra.Command{
		Use:   "assign",
		Short: "Assign a certificate to a namespace",
		Long:  "Assign a root certificate to an attribute namespace for establishing trust",
		Run:   policy_assignCertificateToNamespace,
	}

	policy_certificatesRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove a certificate from a namespace",
		Long:  "Remove a root certificate from an attribute namespace",
		Run:   policy_removeCertificateFromNamespace,
	}
)

func init() {
	// Convert x5c to PEM
	policy_certificatesConvertCmd.Flags().StringP(
		"file",
		"f",
		"",
		"Path to x5c certificate file",
	)
	if err := policy_certificatesConvertCmd.MarkFlagRequired("file"); err != nil {
		panic(err)
	}
	policy_certificatesConvertCmd.Flags().BoolP(
		"output-x5c",
		"x",
		false,
		"Output as x5c format (for PEM to x5c conversion)",
	)

	// Assign certificate to namespace
	policy_certificatesAssignCmd.Flags().StringP(
		"namespace",
		"n",
		"",
		"Namespace ID or FQN",
	)
	if err := policy_certificatesAssignCmd.MarkFlagRequired("namespace"); err != nil {
		panic(err)
	}
	policy_certificatesAssignCmd.Flags().StringP(
		"file",
		"f",
		"",
		"Path to PEM certificate file",
	)
	if err := policy_certificatesAssignCmd.MarkFlagRequired("file"); err != nil {
		panic(err)
	}
	injectLabelFlags(policy_certificatesAssignCmd, false)

	// Remove certificate from namespace
	policy_certificatesRemoveCmd.Flags().StringP(
		"namespace",
		"n",
		"",
		"Namespace ID or FQN",
	)
	if err := policy_certificatesRemoveCmd.MarkFlagRequired("namespace"); err != nil {
		panic(err)
	}
	policy_certificatesRemoveCmd.Flags().StringP(
		"certificate-id",
		"c",
		"",
		"Certificate ID",
	)
	if err := policy_certificatesRemoveCmd.MarkFlagRequired("certificate-id"); err != nil {
		panic(err)
	}

	// Add subcommands
	policy_certificatesCmd.AddCommand(
		policy_certificatesAssignCmd,
		policy_certificatesRemoveCmd,
		policy_certificatesConvertCmd,
	)

	// Register with policy command
	policyCmd.AddCommand(policy_certificatesCmd)
}

func policy_convertPEMToX5C(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)

	filePath := c.Flags.GetRequiredString("file")
	outputX5C := cmd.Flags().Lookup("output-x5c").Value.String() == "true"

	data, err := os.ReadFile(filePath)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to read file (%s)", filePath), err)
	}

	if outputX5C {
		// Convert PEM to x5c
		x5c, err := handlers.ConvertPEMToX5C(data)
		if err != nil {
			cli.ExitWithError("Failed to convert PEM to x5c", err)
		}
		fmt.Println(x5c)
	} else {
		// Convert x5c to PEM (default behavior)
		x5c := string(data)
		pemData, err := handlers.ConvertX5CToPEM(x5c)
		if err != nil {
			cli.ExitWithError("Failed to convert x5c to PEM", err)
		}
		fmt.Println(string(pemData))
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func policy_assignCertificateToNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	namespace := c.Flags.GetRequiredString("namespace")
	filePath := c.Flags.GetRequiredString("file")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	// Read PEM file
	data, err := os.ReadFile(filePath)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to read certificate file (%s)", filePath), err)
	}

	pem := string(data)

	// Get metadata from labels
	metadata := getMetadataMutable(metadataLabels)
	var labels map[string]string
	if metadata != nil {
		labels = metadata.Labels
	}

	// Assign certificate to namespace
	resp, err := h.AssignCertificateToNamespace(cmd.Context(), namespace, pem, labels)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to assign certificate to namespace (%s)", namespace), err)
	}

	// Prepare and display the result
	// Get namespace identifier (either ID or FQN depending on what was provided)
	namespaceIdentifier := resp.GetNamespaceCertificate().GetNamespace().GetId()
	if namespaceIdentifier == "" {
		namespaceIdentifier = resp.GetNamespaceCertificate().GetNamespace().GetFqn()
	}
	rows := [][]string{
		{"Namespace ID", namespaceIdentifier},
		{"Certificate ID", resp.GetNamespaceCertificate().GetCertificateId()},
		{"Certificate (preview)", truncateString(resp.GetCertificate().GetPem(), 80)},
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, namespace, t, resp)
}

func policy_removeCertificateFromNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	namespace := c.Flags.GetRequiredString("namespace")
	certID := c.Flags.GetRequiredID("certificate-id")

	err := h.RemoveCertificateFromNamespace(cmd.Context(), namespace, certID)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to remove certificate (%s) from namespace (%s)", certID, namespace), err)
	}

	// Prepare and display the result
	rows := [][]string{
		{"Removed", "true"},
		{"Namespace", namespace},
		{"Certificate ID", certID},
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, namespace, t, nil)
}
