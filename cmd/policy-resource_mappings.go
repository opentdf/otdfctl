package cmd

import (
	_ "embed"
)

// TODO: add metadata to outputs once [https://github.com/opentdf/tructl/issues/30] is addressed

var (
// policy_resource_mappingsTerms []string

// policy_resource_mappingsCmd = &cobra.Command{
// 	Use:     man.PolicyResourceMappings["en"].Command,
// 	Aliases: man.PolicyResourceMappings["en"].Aliases,
// 	Short: man.PolicyResourceMappings["en"].ShortWithSubCommands([]string{
// 		policy_resource_mappingsCreateCmd.Use,
// 		policy_resource_mappingsGetCmd.Use,
// 		policy_resource_mappingsListCmd.Use,
// 		policy_resource_mappingsUpdateCmd.Use,
// 		policy_resource_mappingsDeleteCmd.Use,
// 	}),
// 	Long: man.PolicyResourceMappings["en"].Long,
// }

// policy_resource_mappingsCreateCmd = &cobra.Command{
// 	Use:   "create",
// 	Short: "Create resource mappings",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		h := cli.NewHandler(cmd)
// 		defer h.Close()

// 		flagHelper := cli.NewFlagHelper(cmd)
// 		attrId := flagHelper.GetRequiredString("attribute-value-id")
// 		terms := flagHelper.GetStringSlice("terms", policy_resource_mappingsTerms, cli.FlagHelperStringSliceOptions{
// 			Min: 1,
// 		})
// 		metadataLabels := flagHelper.GetStringSlice("label", newMetadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

// 		resourceMapping, err := h.CreateResourceMapping(attrId, terms, getMetadata(metadataLabels))
// 		if err != nil {
// 			cli.ExitWithError("Failed to create resource mapping", err)
// 		}

// 		fmt.Println(cli.SuccessMessage("Resource mapping created"))
// 		fmt.Println(cli.NewTabular().Rows([][]string{
// 			{"Id", resourceMapping.Id},
// 			{"Attribute Value Id", resourceMapping.AttributeValue.Id},
// 			{"Attribute Value", resourceMapping.AttributeValue.Value},
// 			{"Terms", strings.Join(resourceMapping.Terms, ", ")},
// 		}...).Render())
// 	},
// }

// policy_resource_mappingsGetCmd = &cobra.Command{
// 	Use:   "get",
// 	Short: "Get resource mappings",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		h := cli.NewHandler(cmd)
// 		defer h.Close()

// 		flagHelper := cli.NewFlagHelper(cmd)
// 		id := flagHelper.GetRequiredString("id")

// 		resourceMapping, err := h.GetResourceMapping(id)
// 		if err != nil {
// 			cli.ExitWithError("Failed to get resource mapping", err)
// 		}

// 		fmt.Println(cli.NewTabular().Rows([][]string{
// 			{"Id", resourceMapping.Id},
// 			{"Attribute Value Id", resourceMapping.AttributeValue.Id},
// 			{"Attribute Value", resourceMapping.AttributeValue.Value},
// 			{"Terms", strings.Join(resourceMapping.Terms, ", ")},
// 		}...).Render())
// 	},
// }

// policy_resource_mappingsListCmd = &cobra.Command{
// 	Use:   "list",
// 	Short: "List resource mappings",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		h := cli.NewHandler(cmd)
// 		defer h.Close()

// 		r, err := h.ListResourceMappings()
// 		if err != nil {
// 			cli.ExitWithError("Failed to list resource mappings", err)
// 		}

// 		t := cli.NewTable()
// 		t.Headers("Id", "Attribute Value Id", "Attribute Value", "Terms")
// 		for _, resourceMapping := range r {
// 			t.Row(resourceMapping.Id, resourceMapping.AttributeValue.Id, resourceMapping.AttributeValue.Value, strings.Join(resourceMapping.Terms, ", "))
// 		}
// 		fmt.Println(t.Render())
// 	},
// }

// policy_resource_mappingsUpdateCmd = &cobra.Command{
// 	Use:   "update",
// 	Short: "Update resource mappings",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		h := cli.NewHandler(cmd)
// 		defer h.Close()

// 		flagHelper := cli.NewFlagHelper(cmd)
// 		id := flagHelper.GetRequiredString("id")
// 		attrValueId := flagHelper.GetOptionalString("attribute-value-id")
// 		terms := flagHelper.GetStringSlice("terms", policy_resource_mappingsTerms, cli.FlagHelperStringSliceOptions{})
// 		newLabels := flagHelper.GetStringSlice("label-new", newMetadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})
// 		replacedLabels := flagHelper.GetStringSlice("label-replace", updatedMetadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

// 		metadata, behavior := processUpdateMetadata(newLabels, replacedLabels, func() (*common.Metadata, error) {
// 			rm, err := h.GetResourceMapping(id)
// 			if err != nil {
// 				errMsg := fmt.Sprintf("Could not find resource mapping (%s)", id)
// 				cli.ExitWithNotFoundError(errMsg, err)
// 				cli.ExitWithError(errMsg, err)
// 			}
// 			return rm.Metadata, nil
// 		},
// 		)

// 		resourceMapping, err := h.UpdateResourceMapping(id, attrValueId, terms, metadata, behavior)
// 		if err != nil {
// 			cli.ExitWithError("Failed to update resource mapping", err)
// 		}

// 		fmt.Println(cli.SuccessMessage("Resource mapping updated"))
// 		fmt.Println(cli.NewTabular().Rows([][]string{
// 			{"Id", resourceMapping.Id},
// 			{"Attribute Value Id", resourceMapping.AttributeValue.Id},
// 			{"Attribute Value", resourceMapping.AttributeValue.Value},
// 			{"Terms", strings.Join(resourceMapping.Terms, ", ")},
// 		}...).Render())
// 	},
// }

// policy_resource_mappingsDeleteCmd = &cobra.Command{
// 	Use:   "delete",
// 	Short: "Delete resource mappings",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		h := cli.NewHandler(cmd)
// 		defer h.Close()

// 		flagHelper := cli.NewFlagHelper(cmd)
// 		id := flagHelper.GetRequiredString("id")

// 		cli.ConfirmDelete("resource-mapping", id)

// 		resourceMapping, err := h.DeleteResourceMapping(id)
// 		if err != nil {
// 			cli.ExitWithError("Failed to delete resource mapping", err)
// 		}

//			fmt.Println(cli.SuccessMessage("Resource mapping deleted"))
//			fmt.Println(cli.NewTabular().Rows([][]string{
//				{"Id", resourceMapping.Id},
//				{"Attribute Value Id", resourceMapping.AttributeValue.Id},
//				{"Attribute Value", resourceMapping.AttributeValue.Value},
//				{"Terms", strings.Join(resourceMapping.Terms, ", ")},
//			}...).Render())
//		},
//	}
)

func init() {
	// policyCmd.AddCommand(policy_resource_mappingsCmd)

	// policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsCreateCmd)
	// policy_resource_mappingsCreateCmd.Flags().String("attribute-value-id", "", "Attribute Value ID")
	// policy_resource_mappingsCreateCmd.Flags().StringSliceVar(&policy_resource_mappingsTerms, "terms", []string{}, "Synonym terms")
	// policy_resource_mappingsCreateCmd.Flags().StringSliceVarP(&newMetadataLabels, "label", "l", []string{}, "Optional metadata 'labels' in the format: key=value")

	// policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsGetCmd)
	// policy_resource_mappingsGetCmd.Flags().String("id", "", "Resource Mapping ID")

	// policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsListCmd)

	// policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsUpdateCmd)
	// policy_resource_mappingsUpdateCmd.Flags().String("id", "", "Resource Mapping ID")
	// policy_resource_mappingsUpdateCmd.Flags().String("attribute-value-id", "", "Attribute Value ID")
	// policy_resource_mappingsUpdateCmd.Flags().StringSliceVar(&policy_resource_mappingsTerms, "terms", []string{}, "Synonym terms")
	// policy_resource_mappingsUpdateCmd.Flags().StringSliceVarP(&newMetadataLabels, "label-new", "n", []string{}, "Optional metadata 'labels' in the format: key=value")
	// policy_resource_mappingsUpdateCmd.Flags().StringSliceVarP(&updatedMetadataLabels, "label-replace", "r", []string{}, "Optional metadata 'labels' in the format: key=value. Note: providing one destructively replaces entire set of labels.")

	// policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsDeleteCmd)
	// policy_resource_mappingsDeleteCmd.Flags().String("id", "", "Resource Mapping ID")
}
