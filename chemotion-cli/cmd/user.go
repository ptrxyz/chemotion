package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Create a new user of Chemotion
var addUser = &cobra.Command{
	Use:        "add <name_of_user>",
	SuggestFor: []string{"add"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("We are now going to add a new user with the name", args[0])
		//TODO
	},
}

// Show details related to a user of Chemotion
var showUser = &cobra.Command{
	Use:        "show <name_of_user>",
	SuggestFor: []string{"sho"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This is what we know about the user", args[0])
		//TODO
	},
}

// Change password for a user of Chemotion
var passwdUser = &cobra.Command{
	Use:        "passwd <name_of_user>",
	SuggestFor: []string{"pas"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Changing password for the user", args[0]+"...")
		//TODO
	},
}

// Remove an existing user of Chemotion
var removeUser = &cobra.Command{
	Use:        "remove <name_of_user>",
	SuggestFor: []string{"rem"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Removing user with the name", args[0]+"...")
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
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Chemotion user actions:")
		switch selectOpt([]string{"add", "show", "passwd", "remove"}) {
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
