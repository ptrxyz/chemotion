package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

// err check function
func check_err(err error, message string) {
	if err != nil {
		if len(message) == 0 {
			message = "An error occured!"
		}
		fmt.Println(message)
		fmt.Println("Error was:", err)
		os.Exit(1)
	}
}

// check if the CLI is running interactively; if not, then exit.
func confirmInteractive() bool {
	if quiet {
		fmt.Println(baseCommand, "has been started in quiet (scripting) mode.\nAll arguments to specify the desired action must be supplied correctly.\nUse the '--help' flag after your command to know what is missing.")
		os.Exit(1)
		return false
	} else {
		return true
	}
}

// selection function
func selectOpt(acceptedOpts []string, remainingArgs []string) (result string) {
	confirmInteractive()
	if len(remainingArgs) != 0 {
		res := false
		for _, opt := range acceptedOpts {
			if remainingArgs[0] == opt {
				panic("Bug: argument should have been handled by Cobra.")
			}
		}
		if !res {
			println("Ignoring invalid argument(s):", strings.Join(remainingArgs, " "))
		}
	}
	var inst string
	if activeInstance == "" {
		inst = "no installed instance found"
	} else {
		inst = activeInstance
	}
	selection := promptui.Select{
		Label: fmt.Sprintf("[%s] Select one of the following", inst),
		Items: acceptedOpts,
	}
	_, result, err := selection.Run()
	check_err(err, "Selection failed!")
	return
}

// string prompt function
func promptStr(label string) (result string) {
	confirmInteractive()
	prompt := promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("can not accept empty value")
			}
			return nil
		},
	}
	result, err := prompt.Run()
	check_err(err, "Prompt failed!")
	return
}

// password prompt
func promptPass(label string) (result string) {
	confirmInteractive()
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("can not accept empty value")
			}
			return nil
		},
	}
	result, err := prompt.Run()
	check_err(err, "Prompt failed!")
	return
}

// process args
func getArg(args []string, prompt string) (arg string) {
	if len(args) == 0 {
		arg = promptStr(prompt)
	} else {
		arg = args[0]
		if len(arg) == 0 { // i.e. an empty argument given as: ""
			println("Can't perform action on an empty argument.")
			os.Exit(1)
		}
	}
	return
}

// get Instance
func getInstance(args []string) (instance string) {
	if len(args) == 0 {
		instance = activeInstance
	} else {
		instance = args[0]
	}
	return
}

var failedMessage = "Failed to execute command: "

// execute command in bash shell
func execShell(command string) (result string) {
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return fmt.Sprintf(failedMessage+"%s with message %s", command, err.Error())
	}
	return string(out)
}

// find version of a software using --version in bash shell
func findVersion(software string) (result string) {
	result = execShell(software + " --version")
	if strings.HasPrefix(result, failedMessage) && strings.HasSuffix(result, "127") {
		return "Not installed!"
	} else {
		return result
	}
}
