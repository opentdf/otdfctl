package forms

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/opentdf/opentdf-v2-poc/sdk/attributes"
)

type AttributeDefinition struct {
	Name        string
	Namespace   string
	Description string
	Labels      map[string]string
	Type        string
	Rule        attributes.AttributeRuleTypeEnum
	Values      []string
}

func AddAttribute() (AttributeDefinition, error) {
	attr := AttributeDefinition{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Namespace").
				Description("Select a namespace. To create a namespace go back and select 'Add Namespace'").
				Options(
					huh.NewOption("demo.com", "demo.com"),
				).
				Value(&attr.Namespace),

			huh.NewInput().
				Title("Attribute Name").
				Value(&attr.Name),

			// Description
			huh.NewText().
				Title("Description").
				Value(&attr.Description),

			// Select Rule
			huh.NewSelect[attributes.AttributeRuleTypeEnum]().
				Title("Rule").
				Options(
					huh.NewOption("All Of", attributes.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_ALL_OF),
					huh.NewOption("Any Of", attributes.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_ANY_OF),
					huh.NewOption("Hierarchical", attributes.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_HIERARCHY),
				).
				Value(&attr.Rule),
		),
	)

	if err := form.Run(); err != nil {
		return attr, err
	}

	for {
		value, another, err := addValue()
		if err != nil {
			return attr, err
		}

		if value == "" {
			fmt.Print("Value cannot be empty\n")
			continue
		}

		attr.Values = append(attr.Values, value)

		if !another {
			break
		}
	}

	return attr, nil
}

func addValue() (value string, another bool, err error) {
	valueForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Value").
				Value(&value),

			huh.NewConfirm().
				Title("Add Another Value").
				Value(&another),
		),
	)

	err = valueForm.Run()

	return value, another, err
}
