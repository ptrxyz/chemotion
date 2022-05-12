/*
Copyright © 2022 Peter Krauß, Shashank S. Harivyasi
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	projectName     = "Chemotion"
	baseCommand     = "chemotion"
	version         = 0.1
	defaultInstance = "default"
	shellVariable   = "CHEMOTION_CONTAINER"
	unknownContName = "UNKNOWN"
)

var (
	chemotionCmdOpts        = []string{"instance", "user", "system"}
	containerName    string = ""
	isInside         bool
)

// chemotionCmd represents the base command when called without any subcommands
var chemotionCmd = &cobra.Command{
	Use:   baseCommand,
	Short: "CLI for Chemotion ELN",
	Long: `Chemotion ELN is an Electronic Lab Notebook solution.
Developed for, and by, researchers, the software aims
to work for you. See, https://www.chemotion.net.`,
	ValidArgs: chemotionCmdOpts,
	// The following lines are the action associated with
	// a bare application run i.e. without any arguments
	Run: func(cmd *cobra.Command, args []string) {
		zlog.Info().Str("container", containerName).Str("instance", config.Instance).Msg("In chemotionCmd. Used command: " + cmd.CommandPath())
		if confirmInteractive() {
			fmt.Println("Welcome to " + projectName + "!")
			if isInside {
				fmt.Println("You are running this " + baseCommand + " CLI tool inside a container called " + containerName + ".")
			} else {
				fmt.Println("You are running this " + baseCommand + " CLI tool on a host machine.")
			}
			fmt.Println("Currently chosen instance is: " + config.Instance + ".")
			selected := selectOpt(chemotionCmdOpts)
			switch selected {
			case "instance":
				// instanceCmd.Run(&cobra.Command{}, []string{})
			case "user":
				// userCmd.Run(&cobra.Command{}, []string{})
			case "system":
				systemCmd.Run(&cobra.Command{}, []string{})
			}
			zlog.Debug().Str("container", containerName).Str("instance", config.Instance).Msg("Selected option: " + selected)
		}
	},
}

// This is called by main.main(). It only needs to happen once.
func Execute() {
	if err := chemotionCmd.Execute(); err == nil {
		zlog.Info().Str("container", containerName).Str("instance", config.Instance).Msg("Chemotion exited gracefully.")
	} else {
		zboth.Fatal().Msg("Chemotion exited abruptly, check log file if necessary. ABORT!")
	}
}

func init() {
	// check where is it running
	containerName, isInside = os.LookupEnv(shellVariable) // using os because we don't want any dependancy here!
	if isInside && containerName == "" {
		containerName = unknownContName
	}
	// begin by setting up logging
	initLog()
	zlog.Info().Str("container", containerName).Str("instance", config.Instance).Msg("Chemotion started.")
	// setup viper
	cobra.OnInitialize(initViper) // in configurer.go
	chemotionCmd.PersistentFlags().StringVarP(&configFile, "config", "f", defaultConfig, "path to the configuration file.")
	// config-file as a flag, cannot be set in the configuration file because that creates a circular dependency
	chemotionCmd.PersistentFlags().BoolVarP(&config.Quiet, "quiet", "q", false, "use "+baseCommand+" in scripted mode i.e. without an interactive prompt")
	viper.BindPFlag("quiet", chemotionCmd.PersistentFlags().Lookup("quiet"))
	chemotionCmd.PersistentFlags().StringVarP(&config.Instance, "instance", "i", "", "start "+baseCommand+" with a pre-selected instance")
	viper.BindPFlag("instance", chemotionCmd.PersistentFlags().Lookup("instance"))
	chemotionCmd.PersistentFlags().BoolVar(&config.Debug, "debug", false, "enable logging of debug messages")
	viper.BindPFlag("debug", chemotionCmd.PersistentFlags().Lookup("debug"))
	zlog.Debug().Str("container", containerName).Str("instance", config.Instance).Msg("Finished root.init().")
}
