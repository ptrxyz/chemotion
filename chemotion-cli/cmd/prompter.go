package cmd

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// Prompt to select a value from a given set of values.
// Also displays the currently selected instance if not same as defaultInstance.
func selectOpt(acceptedOpts []string) (result string) {
	labelPrefix := fmt.Sprintf("[%s on %s] ", config.Instance, containerName)
	if strings.HasPrefix(labelPrefix, fmt.Sprintf("[%s on", defaultInstance)) {
		labelPrefix = strings.Replace(labelPrefix, fmt.Sprintf("[%s on", defaultInstance), "[on", 1)
	}
	if strings.HasSuffix(labelPrefix, "on ] ") {
		labelPrefix = strings.Replace(labelPrefix, "on ] ", "] ", 1)
	}
	if labelPrefix == "[] " {
		labelPrefix = ""
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
