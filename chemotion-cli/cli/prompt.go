package cli

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// Prompt to select a value from a given set of values.
// Also displays the currently selected instance.
func selectOpt(acceptedOpts []string) (result string) {
	zlog.Debug().Msgf("Selection prompt with options %s:", acceptedOpts)
	selection := promptui.Select{
		Label: fmt.Sprintf("%s%s%s%s Select one of the following", string("\033[31m"), string("\033[1m"), currentState.name, string("\033[0m")),
		Items: acceptedOpts,
	}
	_, result, err := selection.Run()
	if err != nil {
		if err.Error() == "^C" {
			zboth.Warn().Err(err).Msgf("Selection cancelled!")
		} else {
			zboth.Fatal().Err(err).Msgf("Selection failed! Check log. ABORT!")
		}
	}
	zlog.Debug().Msgf("Selected option: %s", result)
	return
}

// A simple yes or no question prompt. Yes = True, No = False.
// Important: Cancelled selection leads to result that is defValue. Therefore, make defValue the safer option.
func selectYesNo(question string, defValue bool) (result bool) {
	zlog.Debug().Msgf("Binary question: %s; default is: %t", question, defValue)
	var defValueStr string
	if defValue {
		defValueStr = "y"
	} else {
		defValueStr = "n"
	}
	answer := promptui.Prompt{
		Label:     question,
		IsConfirm: true,
		Default:   defValueStr,
	}
	if _, err := answer.Run(); err == nil {
		result = true
	} else if err.Error() == "" {
		result = false
	} else if err.Error() == "^C" {
		zboth.Warn().Err(fmt.Errorf("yesno prompt cancelled")).Msgf("Selection cancelled. Going for the default (safer) option i.e %s.", strings.ToUpper(defValueStr))
		result = defValue
	} else {
		zboth.Fatal().Err(err).Msgf("Selection failed! Check log. ABORT!")
	}
	zlog.Debug().Msgf("Selected answer: %t", result)
	return
}

// Get user input in form of a string by giving them the message.
func getString(message string) (result string) {
	zlog.Debug().Msgf("String prompt with message: %s", message)
	prompt := promptui.Prompt{
		Label: message,
		Validate: func(input string) (err error) {
			if input == "^C" {
				err = fmt.Errorf("^C")
			} else if len(input) == 0 {
				err = fmt.Errorf("can not accept empty value")
			} else {
				err = nil
			}
			return
		},
	}
	if res, err := prompt.Run(); err == nil {
		zlog.Debug().Msgf("Given answer: %s", res)
		result = res
	} else if err.Error() == "^C" {
		zboth.Fatal().Err(fmt.Errorf("string prompt cancelled")).Msg("Input cancelled. Can't proceed without. Abort!")
	} else {
		zboth.Fatal().Err(err).Msgf("Prompt failed for unknown reason. Check log. Abort!")
	}
	return
}
