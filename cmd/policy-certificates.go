package cmd

import (
	"fmt"
	"os"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/spf13/cobra"
)

var (
	policy_certificatesCmd = &cobra.Command{
		Use:   "certificates",
		Short: "Manage certificates for namespaces",
		Long:  "Commands to create, list, and manage certificates for attribute namespaces",
	}

	policy_certificatesListCmd = &cobra.Command{
		Use:   "list",
		Short: "List certificates for a namespace",
		Long:  "List all root certificates assigned to a namespace",
		Run:   policy_listNamespaceCertificates,
	}

	policy_certificatesGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a certificate by ID",
		Long:  "Retrieve details of a specific certificate",
		Run:   policy_getCertificate,
	}

	policy_certificatesShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show namespace with certificates",
		Long:  "Display a namespace and all its associated root certificates",
		Run:   policy_showNamespaceWithCertificates,
	}

	policy_certificatesConvertCmd = &cobra.Command{
		Use:   "convert-pem",
		Short: "Convert PEM certificate to x5c format",
		Long:  "Convert a PEM-encoded certificate file to x5c format (base64-encoded DER)",
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
	// List certificates for namespace
	policy_certificatesListCmd.Flags().StringP(
		"namespace-id",
		"n",
		"",
		"Namespace ID or FQN to list certificates for",
	)
	if err := policy_certificatesListCmd.MarkFlagRequired("namespace-id"); err != nil {
		panic(err)
	}

	// Get certificate
	policy_certificatesGetCmd.Flags().StringP(
		"id",
		"i",
		"",
		"Certificate ID",
	)
	if err := policy_certificatesGetCmd.MarkFlagRequired("id"); err != nil {
		panic(err)
	}

	// Show namespace with certificates
	policy_certificatesShowCmd.Flags().StringP(
		"namespace-id",
		"n",
		"",
		"Namespace ID or FQN",
	)
	if err := policy_certificatesShowCmd.MarkFlagRequired("namespace-id"); err != nil {
		panic(err)
	}

	// Convert PEM to x5c
	policy_certificatesConvertCmd.Flags().StringP(
		"file",
		"f",
		"",
		"Path to PEM certificate file",
	)
	if err := policy_certificatesConvertCmd.MarkFlagRequired("file"); err != nil {
		panic(err)
	}
	policy_certificatesConvertCmd.Flags().BoolP(
		"output-pem",
		"p",
		false,
		"Output as PEM format (for x5c to PEM conversion)",
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
		policy_certificatesListCmd,
		policy_certificatesGetCmd,
		policy_certificatesShowCmd,
		policy_certificatesConvertCmd,
		policy_certificatesAssignCmd,
		policy_certificatesRemoveCmd,
	)

	// Register with policy command
	policyCmd.AddCommand(policy_certificatesCmd)
}

func policy_listNamespaceCertificates(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	namespaceID := c.Flags.GetRequiredString("namespace-id")

	certs, err := h.ListNamespaceCertificates(cmd.Context(), namespaceID)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to list certificates for namespace (%s)", namespaceID), err)
	}

	if len(certs) == 0 {
		fmt.Println("No certificates found for this namespace")
		return
	}

	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewFlexColumn("x5c_preview", "Certificate (Preview)", cli.FlexColumnWidthFour),
	)

	rows := []table.Row{}
	for _, cert := range certs {
		preview := cert.GetX5C()
		if len(preview) > 50 {
			preview = preview[:50] + "..."
		}
		rows = append(rows,
			table.NewRow(table.RowData{
				"id":          cert.GetId(),
				"x5c_preview": preview,
			}),
		)
	}
	t = t.WithRows(rows)
	HandleSuccess(cmd, "", t, certs)
}

func policy_getCertificate(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredString("id")

	cert, err := h.GetCertificate(cmd.Context(), id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get certificate (%s)", id), err)
	}

	rows := [][]string{
		{"Id", cert.GetId()},
		{"X5C (first 100 chars)", truncateString(cert.GetX5C(), 100)},
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, cert.GetId(), t, cert)
}

func policy_showNamespaceWithCertificates(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	namespaceID := c.Flags.GetRequiredString("namespace-id")

	ns, err := h.GetNamespaceWithCertificates(cmd.Context(), namespaceID)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get namespace (%s)", namespaceID), err)
	}

	rows := [][]string{
		{"Id", ns.GetId()},
		{"Name", ns.GetName()},
		{"FQN", ns.GetFqn()},
		{"Certificate Count", fmt.Sprintf("%d", len(ns.GetRootCerts()))},
	}
	if mdRows := getMetadataRows(ns.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	// Add certificate details
	if len(ns.GetRootCerts()) > 0 {
		rows = append(rows, []string{"", ""}) // Spacer
		rows = append(rows, []string{"=== Root Certificates ===", ""})
		for i, cert := range ns.GetRootCerts() {
			rows = append(rows, []string{fmt.Sprintf("Certificate %d ID", i+1), cert.GetId()})
			rows = append(rows, []string{fmt.Sprintf("Certificate %d (preview)", i+1), truncateString(cert.GetX5C(), 80)})
		}
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, ns.GetId(), t, ns)
}

func policy_convertPEMToX5C(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)

	filePath := c.Flags.GetRequiredString("file")
	outputPEM := cmd.Flags().Lookup("output-pem").Value.String() == "true"

	data, err := os.ReadFile(filePath)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to read file (%s)", filePath), err)
	}

	if outputPEM {
		// Convert x5c back to PEM
		x5c := string(data)
		pemData, err := handlers.ConvertX5CToPEM(x5c)
		if err != nil {
			cli.ExitWithError("Failed to convert x5c to PEM", err)
		}
		fmt.Println(string(pemData))
	} else {
		// Convert PEM to x5c
		x5c, err := handlers.ConvertPEMToX5C(data)
		if err != nil {
			cli.ExitWithError("Failed to convert PEM to x5c", err)
		}
		fmt.Println(x5c)
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

	// Read and convert PEM file to x5c
	data, err := os.ReadFile(filePath)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to read certificate file (%s)", filePath), err)
	}

	x5c, err := handlers.ConvertPEMToX5C(data)
	if err != nil {
		cli.ExitWithError("Failed to convert PEM to x5c format", err)
	}

	// Get metadata from labels
	metadata := getMetadataMutable(metadataLabels)
	var labels map[string]string
	if metadata != nil {
		labels = metadata.Labels
	}

	// Assign certificate to namespace
	resp, err := h.AssignCertificateToNamespace(cmd.Context(), namespace, x5c, labels)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to assign certificate to namespace (%s)", namespace), err)
	}

	// Prepare and display the result
	rows := [][]string{
		{"Namespace ID", resp.GetNamespaceCertificate().GetNamespaceId()},
		{"Certificate ID", resp.GetNamespaceCertificate().GetCertificateId()},
		{"Certificate (preview)", truncateString(resp.GetCertificate().GetX5C(), 80)},
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
