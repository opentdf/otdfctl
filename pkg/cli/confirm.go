package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

const (
	// top level actions
	ActionGet        = "get"
	ActionList       = "list"
	ActionCreate     = "create"
	ActionUpdate     = "update"
	ActionDeactivate = "deactivate"
	ActionReactivate = "reactivate"
	ActionDelete     = "delete"

	// member actions
	ActionMemberAdd     = "add members"
	ActionMemberRemove  = "remove members"
	ActionMemberReplace = "replace all existing members"

	// text input names
	InputNameFQN = "fully qualified name (FQN)"
)

func ConfirmAction(action, resource, id string) {
	var confirm bool
	err := huh.NewConfirm().
		Title(fmt.Sprintf("Are you sure you want to %s %s:\n\n\t%s", action, resource, id)).
		Affirmative("yes").
		Negative("no").
		Value(&confirm).
		Run()
	if err != nil {
		ExitWithError("Confirmation prompt failed", err)
	}

	if !confirm {
		fmt.Println(ErrorMessage("Aborted", nil))
		os.Exit(0)
	}
}

func ConfirmTextInput(action, resource, inputName, shouldMatchValue string) {
	var input string
	err := huh.NewInput().
		Title(fmt.Sprintf("To confirm you want to %s this %s and accept any side effects, please enter the %s to proceed: %s", action, resource, inputName, shouldMatchValue)).
		Value(&input).
		Validate(func(s string) error {
			if s != shouldMatchValue {
				return fmt.Errorf(fmt.Sprintf("FQN entered [%s] does not match required %s: %s", s, inputName, shouldMatchValue))
			}
			return nil
		}).Run()
	if err != nil {
		ExitWithError("Confirmation prompt failed", err)
	}
}
