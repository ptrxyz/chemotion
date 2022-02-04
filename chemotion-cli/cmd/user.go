package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Create a new user of Chemotion
var addUser = &cobra.Command{
	Use:        "add <name_of_user>",
	SuggestFor: []string{"add"},
	Args:       cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		user := getArg(args, "Please enter username of the new user")
		fmt.Println("We are now going to add a new user with the name", user)
		//TODO
	},
}

// Show details related to a user of Chemotion
var showUser = &cobra.Command{
	Use:        "show <name_of_user>",
	SuggestFor: []string{"sho"},
	Args:       cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		user := getArg(args, "Please enter the username whose details are required")
		fmt.Println("This is what we know about the user", user)
		//TODO
	},
}

// Change password for a user of Chemotion
var passwdUser = &cobra.Command{
	Use:        "passwd <name_of_user> <new_passwd>",
	SuggestFor: []string{"pas"},
	Args:       cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		user := getArg(args, "Please enter the username whose password needs to be changed")
		var passwd string
		if len(args) == 2 {
			passwd = args[1]
		} else {
			passwd1 := promptPass("Please enter new password for " + user)
			passwd2 := promptPass("Please confirm the new password")
			if passwd1 == passwd2 {
				passwd = passwd1
			} else {
				fmt.Println("Passwords do not match. Exiting!")
				os.Exit(1)
			}
		}
		// TODO remove this clear text password !!!
		fmt.Println("The new password for", user, "is", passwd)
	},
}

// Remove an existing user of Chemotion
var removeUser = &cobra.Command{
	Use:        "remove <name_of_user>",
	SuggestFor: []string{"rem"},
	Args:       cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		user := getArg(args, "Please enter the username you want to remove")
		fmt.Println("Removing username", user)
		//TODO
	},
}

// Backbone for user-related commands
var userCmd = &cobra.Command{
	Use:        "user {show|add|remove|passwd} <name_of_user>",
	Aliases:    []string{"u"},
	SuggestFor: []string{"u"},
	Short:      "Perform user-related actions in " + baseCommand,
	Long:       "Perform user-related actions in " + baseCommand + " using one of the available actions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Chemotion. Actions on users.")
		confirmInteractive()
		acceptedOpts := []string{"add", "show", "passwd", "remove"}
		switch selectOpt(acceptedOpts, args) {
		case "add":
			addUser.Run(&cobra.Command{}, []string{})
		case "show":
			showUser.Run(&cobra.Command{}, []string{})
		case "passwd":
			passwdUser.Run(&cobra.Command{}, []string{})
		case "remove":
			removeUser.Run(&cobra.Command{}, []string{})
		}
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(addUser)
	userCmd.AddCommand(showUser)
	userCmd.AddCommand(passwdUser)
	userCmd.AddCommand(removeUser)
}
