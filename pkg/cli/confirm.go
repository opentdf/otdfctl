package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

func ConfirmDelete(resource string, id string) {
	var confirm bool
	err := huh.NewConfirm().
		Title(fmt.Sprintf("Are you sure you want to delete %s:\n\n\t%s", resource, id)).
		Affirmative("yes").
		Negative("no").
		Value(&confirm).
		Run()
	if err != nil {
		fmt.Println(ErrorMessage("Confirmation prompt failed", err))
		os.Exit(1)
	}

	if !confirm {
		fmt.Println(ErrorMessage("Aborted", nil))
		os.Exit(1)
	}
}
