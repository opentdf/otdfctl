package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

const (
	ActionGet        = "get"
	ActionList       = "list"
	ActionCreate     = "create"
	ActionUpdate     = "update"
	ActionDeactivate = "deactivate"
	ActionDelete     = "delete"
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
