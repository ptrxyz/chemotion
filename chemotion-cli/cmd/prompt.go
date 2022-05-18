package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

// Prompt to select a value from a given set of values.
// Also displays the currently selected instance if not same as defaultInstanceName.
func selectOpt(acceptedOpts []string) (result string) {
	var labelPrefix string
	if currentState.name == defaultInstanceName {
		labelPrefix = ""
	} else {
		labelPrefix = fmt.Sprintf("[%s] ", currentState.name)
	}
	selection := promptui.Select{
		Label: labelPrefix + "Select one of the following",
		Items: acceptedOpts,
	}
	_, result, err := selection.Run()
	if err != nil {
		if err.Error() == "^C" {
			zboth.Warn().Err(err).Msg("Selection cancelled!")
		} else {
			zboth.Fatal().Err(err).Msg("Selection failed! ABORT!")
		}
	}
	return
}
