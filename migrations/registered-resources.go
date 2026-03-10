//nolint:forbidigo // migration output requires direct terminal printing for interactive prompts and styled output
package migrations

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/charmbracelet/huh"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
	"github.com/opentdf/platform/protocol/go/policy/registeredresources"
)

const (
	optSkipResource = "skip-resource"
	optAbortAll     = "abort-all"
)

// MigrationHandler defines the handler methods needed for registered resource migration.
// handlers.Handler satisfies this interface implicitly.
type MigrationHandler interface {
	ListRegisteredResources(ctx context.Context, limit, offset int32, namespace string) (*registeredresources.ListRegisteredResourcesResponse, error)
	ListRegisteredResourceValues(ctx context.Context, resourceID string, limit, offset int32) (*registeredresources.ListRegisteredResourceValuesResponse, error)
	CreateRegisteredResource(ctx context.Context, namespace, name string, values []string, metadata *common.MetadataMutable) (*policy.RegisteredResource, error)
	CreateRegisteredResourceValue(ctx context.Context, resourceID string, value string, actionAttributeValues []*registeredresources.ActionAttributeValue, metadata *common.MetadataMutable) (*policy.RegisteredResourceValue, error)
	DeleteRegisteredResource(ctx context.Context, id string) error
	ListNamespaces(ctx context.Context, state common.ActiveStateEnum, limit, offset int32) (*namespaces.ListNamespacesResponse, error)
}

// RegisteredResourceMigrationPlan holds an existing resource with its values and the target namespace.
type RegisteredResourceMigrationPlan struct {
	Resource        *policy.RegisteredResource
	Values          []*policy.RegisteredResourceValue
	TargetNamespace string // namespace FQN or ID to migrate to
	Commit          bool
}

// MigrateRegisteredResources is the main entry point for migrating registered resources to namespaces.
func MigrateRegisteredResources(ctx context.Context, h MigrationHandler, commit, interactive bool) error {
	styles := initMigrationDisplayStyles()

	didBackup, err := backupForm()
	if err != nil {
		return err
	}
	if !didBackup {
		return errors.New("user did not confirm backup")
	}

	plan, err := buildRegisteredResourcePlan(ctx, h)
	if err != nil {
		return err
	}

	if len(plan) == 0 {
		fmt.Println(styles.styleWarning.Render("No registered resources found that need namespace migration."))
		return nil
	}

	availableNamespaces, err := listAvailableNamespaces(ctx, h)
	if err != nil {
		return err
	}

	if len(availableNamespaces) == 0 {
		return errors.New("no namespaces available - please create at least one namespace before running migration")
	}

	switch {
	case interactive && commit:
		return runInteractiveRegisteredResourceMigration(ctx, h, styles, plan, availableNamespaces)
	case commit:
		return runBatchRegisteredResourceMigration(ctx, h, styles, plan, availableNamespaces)
	default:
		displayRegisteredResourcePlan(styles, plan)
		if interactive {
			fmt.Println(styles.styleInfo.Render("\nNote: --interactive without --commit only shows a preview. Add --commit to apply changes."))
		}
	}

	return nil
}

// buildRegisteredResourcePlan fetches all registered resources without namespaces and their values.
func buildRegisteredResourcePlan(ctx context.Context, h MigrationHandler) ([]RegisteredResourceMigrationPlan, error) {
	var (
		plans    []RegisteredResourceMigrationPlan
		offset   int32
		pageSize int32 = 100
	)

	for {
		resp, err := h.ListRegisteredResources(ctx, pageSize, offset, "")
		if err != nil {
			return nil, fmt.Errorf("failed to list registered resources: %w", err)
		}

		resources := resp.GetResources()
		if len(resources) == 0 {
			break
		}

		for _, resource := range resources {
			// Only include resources that have no namespace
			if resource.GetNamespace() != nil && resource.GetNamespace().GetId() != "" {
				continue
			}

			values, err := fetchAllResourceValues(ctx, h, resource.GetId())
			if err != nil {
				return nil, fmt.Errorf("failed to fetch values for resource %s: %w", resource.GetId(), err)
			}

			plans = append(plans, RegisteredResourceMigrationPlan{
				Resource: resource,
				Values:   values,
			})
		}

		qty := len(resources)
		if qty > math.MaxInt32 || offset+int32(qty) < 0 {
			return nil, errors.New("resource count exceeded safe limit")
		}
		offset += int32(qty)

		if int32(qty) < pageSize {
			break
		}
	}

	return plans, nil
}

// fetchAllResourceValues paginates through all values for a resource.
func fetchAllResourceValues(ctx context.Context, h MigrationHandler, resourceID string) ([]*policy.RegisteredResourceValue, error) {
	var (
		allValues []*policy.RegisteredResourceValue
		offset    int32
		pageSize  int32 = 100
	)

	for {
		resp, err := h.ListRegisteredResourceValues(ctx, resourceID, pageSize, offset)
		if err != nil {
			return nil, err
		}

		values := resp.GetValues()
		if len(values) == 0 {
			break
		}

		allValues = append(allValues, values...)

		qty := len(values)
		if qty > math.MaxInt32 || offset+int32(qty) < 0 {
			break
		}
		offset += int32(qty)

		if int32(qty) < pageSize {
			break
		}
	}

	return allValues, nil
}

// listAvailableNamespaces fetches all active namespaces.
func listAvailableNamespaces(ctx context.Context, h MigrationHandler) ([]*policy.Namespace, error) {
	var (
		all      []*policy.Namespace
		offset   int32
		pageSize int32 = 100
	)

	for {
		resp, err := h.ListNamespaces(ctx, common.ActiveStateEnum_ACTIVE_STATE_ENUM_ACTIVE, pageSize, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to list namespaces: %w", err)
		}

		nsList := resp.GetNamespaces()
		if len(nsList) == 0 {
			break
		}

		all = append(all, nsList...)

		qty := len(nsList)
		if qty > math.MaxInt32 || offset+int32(qty) < 0 {
			break
		}
		offset += int32(qty)

		if int32(qty) < pageSize {
			break
		}
	}

	return all, nil
}

// buildNamespaceOptions creates huh options from a list of namespaces.
func buildNamespaceOptions(nsList []*policy.Namespace) []huh.Option[string] {
	opts := make([]huh.Option[string], 0, len(nsList))
	for _, ns := range nsList {
		label := ns.GetFqn()
		if label == "" {
			label = ns.GetName() + " (" + ns.GetId() + ")"
		}
		opts = append(opts, huh.NewOption(label, ns.GetFqn()))
	}
	return opts
}

// displayRegisteredResourcePlan shows a preview of resources that would be migrated.
func displayRegisteredResourcePlan(styles *migrationDisplayStyles, plan []RegisteredResourceMigrationPlan) {
	fmt.Println(styles.styleTitle.Render("\nRegistered Resources Migration Plan"))
	fmt.Println(styles.styleSeparator.Render(styles.separatorText))
	fmt.Printf("%s %d\n\n",
		styles.styleInfo.Render("Resources requiring namespace assignment:"),
		len(plan),
	)

	for i, p := range plan {
		fmt.Printf("%s %s\n",
			styles.styleInfo.Render(fmt.Sprintf("%d. Resource ID:", i+1)),
			styles.styleResourceID.Render(p.Resource.GetId()),
		)
		fmt.Printf("   %s %s\n",
			styles.styleInfo.Render("Name:"),
			styles.styleName.Render(p.Resource.GetName()),
		)
		if len(p.Values) > 0 {
			fmt.Printf("   %s\n", styles.styleInfo.Render("Values:"))
			for _, v := range p.Values {
				aavCount := len(v.GetActionAttributeValues())
				fmt.Printf("     - %s (ID: %s, %d action-attribute mapping(s))\n",
					styles.styleValue.Render(v.GetValue()),
					styles.styleID.Render(v.GetId()),
					aavCount,
				)
			}
		} else {
			fmt.Printf("   %s\n", styles.styleInfo.Render("Values: (none)"))
		}
		fmt.Println()
	}

	fmt.Println(styles.styleSeparator.Render(styles.separatorText))
	fmt.Println(styles.styleInfo.Render("\nRun with --commit to assign namespaces in batch mode."))
	fmt.Println(styles.styleInfo.Render("Run with --interactive --commit for per-resource namespace assignment."))
}

// runBatchRegisteredResourceMigration prompts once for a namespace, then migrates all resources.
func runBatchRegisteredResourceMigration(ctx context.Context, h MigrationHandler, styles *migrationDisplayStyles, plan []RegisteredResourceMigrationPlan, nsList []*policy.Namespace) error {
	displayRegisteredResourcePlan(styles, plan)

	// Prompt for target namespace
	var targetNamespace string
	nsOpts := buildNamespaceOptions(nsList)

	namespaceForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a target namespace for ALL registered resources:").
				Options(nsOpts...).
				Value(&targetNamespace),
		),
	)

	if err := namespaceForm.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return errors.New("migration aborted by user")
		}
		return fmt.Errorf("namespace selection failed: %w", err)
	}

	fmt.Println(styles.styleInfo.Render(fmt.Sprintf("\nMigrating all %d resources to namespace: %s\n", len(plan), styles.styleNamespace.Render(targetNamespace))))

	successCount := 0
	failedResources := make(map[string]string)

	for _, p := range plan {
		p.TargetNamespace = targetNamespace
		p.Commit = true

		fmt.Printf("%s %s (%s)...\n",
			styles.styleInfo.Render("Migrating resource"),
			styles.styleName.Render(p.Resource.GetName()),
			styles.styleResourceID.Render(p.Resource.GetId()),
		)

		if err := commitRegisteredResourceMigration(ctx, h, p); err != nil {
			errMsg := fmt.Sprintf("Failed to migrate resource %s: %v", p.Resource.GetId(), err)
			fmt.Println(styles.styleWarning.Render(errMsg))
			failedResources[p.Resource.GetId()] = err.Error()
		} else {
			fmt.Println(styles.styleAction.Render(fmt.Sprintf("  Successfully migrated resource %s", p.Resource.GetName())))
			successCount++
		}
	}

	// Print summary
	fmt.Println(styles.styleTitle.Render("\nBatch Migration Summary:"))
	fmt.Printf("  Total Resources: %d\n", len(plan))
	fmt.Printf("  Successfully Migrated: %d\n", successCount)
	fmt.Printf("  Failed: %d\n", len(failedResources))
	if len(failedResources) > 0 {
		fmt.Println(styles.styleWarning.Render("  Failed Resources:"))
		for id, errMsg := range failedResources {
			fmt.Printf("    - Resource ID %s: %s\n", styles.styleResourceID.Render(id), errMsg)
		}
	}

	return nil
}

// runInteractiveRegisteredResourceMigration prompts per-resource for namespace assignment.
func runInteractiveRegisteredResourceMigration(ctx context.Context, h MigrationHandler, styles *migrationDisplayStyles, plan []RegisteredResourceMigrationPlan, nsList []*policy.Namespace) error {
	fmt.Println(styles.styleInfo.Render("Interactive mode: processing resources one by one..."))

	nsOpts := buildNamespaceOptions(nsList)
	// Add skip and abort options
	skipOpt := huh.NewOption("Skip this resource", optSkipResource)
	abortOpt := huh.NewOption("Abort entire migration", optAbortAll)
	nsOptsWithControls := append(append([]huh.Option[string]{}, nsOpts...), skipOpt, abortOpt)

	successCount := 0
	skippedCount := 0
	failedResources := make(map[string]string)

	for i, p := range plan {
		fmt.Println(styles.styleSeparator.Render(styles.separatorText))
		fmt.Printf("%s %s (%s %s)\n",
			styles.styleTitle.Render(fmt.Sprintf("Resource %d/%d:", i+1, len(plan))),
			styles.styleName.Render(p.Resource.GetName()),
			styles.styleInfo.Render("ID:"),
			styles.styleResourceID.Render(p.Resource.GetId()),
		)

		if len(p.Values) > 0 {
			fmt.Printf("  %s\n", styles.styleInfo.Render("Values:"))
			for _, v := range p.Values {
				aavCount := len(v.GetActionAttributeValues())
				fmt.Printf("    - %s (%d action-attribute mapping(s))\n",
					styles.styleValue.Render(v.GetValue()),
					aavCount,
				)
			}
		}

		var targetNamespace string
		namespaceForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title(fmt.Sprintf("Select namespace for resource '%s':", p.Resource.GetName())).
					Options(nsOptsWithControls...).
					Value(&targetNamespace),
			),
		)

		if err := namespaceForm.Run(); err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				fmt.Println(styles.styleWarning.Render("Migration aborted by user."))
				break
			}
			fmt.Println(styles.styleWarning.Render(fmt.Sprintf("Error during prompt: %v. Skipping resource.", err)))
			skippedCount++
			continue
		}

		switch targetNamespace {
		case optSkipResource:
			fmt.Println(styles.styleInfo.Render(fmt.Sprintf("Skipping resource %s.", p.Resource.GetName())))
			skippedCount++
			continue
		case optAbortAll:
			fmt.Println(styles.styleWarning.Render("Aborting migration."))
			goto summary
		}

		p.TargetNamespace = targetNamespace
		p.Commit = true

		fmt.Printf("%s %s to namespace %s...\n",
			styles.styleAction.Render("  Migrating"),
			styles.styleName.Render(p.Resource.GetName()),
			styles.styleNamespace.Render(targetNamespace),
		)

		if err := commitRegisteredResourceMigration(ctx, h, p); err != nil {
			errMsg := fmt.Sprintf("Failed to migrate resource %s: %v", p.Resource.GetId(), err)
			fmt.Println(styles.styleWarning.Render(errMsg))
			failedResources[p.Resource.GetId()] = err.Error()
		} else {
			fmt.Println(styles.styleAction.Render(fmt.Sprintf("  Successfully migrated resource %s", p.Resource.GetName())))
			successCount++
		}
	}

summary:
	fmt.Println(styles.styleTitle.Render("\nInteractive Migration Summary:"))
	fmt.Printf("  Total Resources: %d\n", len(plan))
	fmt.Printf("  Successfully Migrated: %d\n", successCount)
	fmt.Printf("  Skipped: %d\n", skippedCount)
	fmt.Printf("  Failed: %d\n", len(failedResources))
	if len(failedResources) > 0 {
		fmt.Println(styles.styleWarning.Render("  Failed Resources:"))
		for id, errMsg := range failedResources {
			fmt.Printf("    - Resource ID %s: %s\n", styles.styleResourceID.Render(id), errMsg)
		}
	}

	return nil
}

// commitRegisteredResourceMigration re-creates a resource under a target namespace, then deletes the old one.
func commitRegisteredResourceMigration(ctx context.Context, h MigrationHandler, plan RegisteredResourceMigrationPlan) error {
	if !plan.Commit || plan.TargetNamespace == "" {
		return errors.New("migration plan is not ready for commit")
	}

	resource := plan.Resource

	// Build metadata for the new resource
	var metadata *common.MetadataMutable
	if resource.GetMetadata() != nil && len(resource.GetMetadata().GetLabels()) > 0 {
		metadata = &common.MetadataMutable{
			Labels: resource.GetMetadata().GetLabels(),
		}
	}

	// Step 1: Create new resource under target namespace (without values — we create them individually)
	newResource, err := h.CreateRegisteredResource(ctx, plan.TargetNamespace, resource.GetName(), nil, metadata)
	if err != nil {
		return fmt.Errorf("failed to create resource under namespace %s: %w", plan.TargetNamespace, err)
	}

	// Step 2: Create each value individually, preserving action-attribute mappings
	for _, oldValue := range plan.Values {
		oldAAVs := oldValue.GetActionAttributeValues()
		var aavRequests []*registeredresources.ActionAttributeValue
		if len(oldAAVs) > 0 {
			aavRequests = convertActionAttributeValues(oldAAVs)
		}

		var valueMetadata *common.MetadataMutable
		if oldValue.GetMetadata() != nil && len(oldValue.GetMetadata().GetLabels()) > 0 {
			valueMetadata = &common.MetadataMutable{
				Labels: oldValue.GetMetadata().GetLabels(),
			}
		}

		_, err := h.CreateRegisteredResourceValue(ctx, newResource.GetId(), oldValue.GetValue(), aavRequests, valueMetadata)
		if err != nil {
			return fmt.Errorf("failed to create value %s for resource %s: %w", oldValue.GetValue(), newResource.GetId(), err)
		}
	}

	// Step 3: Delete old resource (cascades to its values)
	if err := h.DeleteRegisteredResource(ctx, resource.GetId()); err != nil {
		return fmt.Errorf("failed to delete old resource %s (new resource %s was created successfully - manual cleanup may be needed): %w",
			resource.GetId(), newResource.GetId(), err)
	}

	return nil
}

// convertActionAttributeValues converts from policy object AAVs to request AAVs.
func convertActionAttributeValues(aavs []*policy.RegisteredResourceValue_ActionAttributeValue) []*registeredresources.ActionAttributeValue {
	result := make([]*registeredresources.ActionAttributeValue, 0, len(aavs))
	for _, aav := range aavs {
		req := &registeredresources.ActionAttributeValue{}

		// Use action ID if available
		if action := aav.GetAction(); action != nil && action.GetId() != "" {
			req.ActionIdentifier = &registeredresources.ActionAttributeValue_ActionId{
				ActionId: action.GetId(),
			}
		}

		// Use attribute value ID if available
		if attrValue := aav.GetAttributeValue(); attrValue != nil && attrValue.GetId() != "" {
			req.AttributeValueIdentifier = &registeredresources.ActionAttributeValue_AttributeValueId{
				AttributeValueId: attrValue.GetId(),
			}
		}

		result = append(result, req)
	}
	return result
}

// backupForm prompts the user to confirm they have taken a backup.
func backupForm() (bool, error) {
	var backupResponse bool
	styles := initMigrationDisplayStyles()

	fmt.Println(styles.styleWarning.Render("WARNING: This operation will delete and re-create registered resources under new namespaces."))
	fmt.Println(styles.styleWarning.Render("It is STRONGLY recommended to take a complete backup of your system before proceeding.\n"))

	backupResponseForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title("Have you taken a complete backup? (yes/no): ").
				Options(
					huh.NewOption("yes", true),
					huh.NewOption("no", false),
					huh.NewOption("cancel", false),
				).
				Value(&backupResponse),
		),
	)

	if err := backupResponseForm.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return false, errors.New("user aborted backup form")
		}
		return false, err
	}
	return backupResponse, nil
}
