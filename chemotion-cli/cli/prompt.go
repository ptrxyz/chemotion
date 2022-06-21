package cli

import (
	"fmt"
	"strconv"
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
	if err == nil {
		zboth.Debug().Msgf("Selected option: %s", result)
	} else if err == promptui.ErrInterrupt || err == promptui.ErrEOF {
		zboth.Fatal().Err(err).Msgf("Selection cancelled!")
	} else {
		zboth.Fatal().Err(err).Msgf("Selection failed! Check log. ABORT!")
	}
	return
}

// A simple yes or no question prompt. Yes = True, No = False.
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
	} else if err == promptui.ErrAbort {
		result = false
	} else if err == promptui.ErrInterrupt || err == promptui.ErrEOF {
		zboth.Fatal().Err(fmt.Errorf("yesno prompt cancelled")).Msgf("Selection cancelled.")
	} else {
		zboth.Fatal().Err(err).Msgf("Selection failed! Check log. ABORT!")
	}
	zlog.Debug().Msgf("Selected answer: %t", result)
	return
}

func textValidate(input string) (err error) {
	if len(strings.ReplaceAll(input, " ", "")) == 0 {
		err = fmt.Errorf("can not accept empty value")
	} else if len(strings.Fields(input)) > 1 {
		err = fmt.Errorf("can not have spaces in the input")
	} else {
		err = nil
	}
	return
}

func instanceValidate(input string) (err error) {
	err = textValidate(input)
	if err == nil {
		existingInstances := allInstances()
		if stringInArray(input, &existingInstances) > -1 {
			err = nil
		} else {
			err = fmt.Errorf("there is no instance called %s", input)
		}
	}
	return
}

func addressValidate(input string) (err error) {
	err = textValidate(input)
	protocol, address, found := strings.Cut(input, ":")
	if found {
		address = strings.TrimPrefix(address, "//")
		protocol += "://"
	}
	if err == nil {
		if !found || !((protocol == "http://") || (protocol == "https://")) {
			err = fmt.Errorf("address must start with protocol i.e. as `http://` or as `https://`")
		}
	}
	address, port, portGiven := strings.Cut(address, ":")
	if err == nil {
		if err = textValidate(address); err != nil {
			if err.Error() == "can not accept empty value" {
				err = fmt.Errorf("can not accept empty value for address")
			}
		}
	}
	if err == nil && portGiven {
		_, err = strconv.Atoi(port)
		if err != nil {
			err = fmt.Errorf("port must be an integer")
		}
	}
	return
}

// kind of opposite of instanceValidate
func newInstanceValidate(input string) (err error) {
	err = textValidate(input)
	if err == nil {
		if exists := instanceValidate(input); exists == nil {
			err = fmt.Errorf("this value is alredy taken")
		} else {
			err = nil
		}
	}
	return
}

// Get user input in form of a string by giving them the message.
func getString(message string, validator promptui.ValidateFunc) (result string) {
	zlog.Debug().Msgf("String prompt with message: %s", message)
	prompt := promptui.Prompt{
		Label:    message,
		Validate: validator,
	}
	if res, err := prompt.Run(); err == nil {
		zlog.Debug().Msgf("Given answer: %s", res)
		result = res
	} else if err == promptui.ErrInterrupt || err == promptui.ErrEOF {
		zboth.Fatal().Err(fmt.Errorf("prompt cancelled")).Msgf("Prompt cancelled. Can't proceed without. ABORT!")
	} else {
		zboth.Fatal().Err(err).Msgf("Prompt failed because: %s.", err.Error())
	}
	return
}

// to select an instance, gives a list to select from when less than 5, else a text input
func selectInstance(action string) (instance string) {
	existingInstances := allInstances()
	if len(existingInstances) < 5 {
		fmt.Printf("Please pick the instance to %s:\n", action)
		instance = selectOpt(existingInstances)
	} else {
		zlog.Debug().Msgf("String prompt to select instance")
		instance = getString("Please name the instance to "+action, instanceValidate)
	}
	return
}
