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

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	version               = 0.1
	projectName           = "Chemotion"
	defaultConfigFileName = "chemotion-cli.yml"
	logFileName           = "chemotion-cli.log"
	stateFPath            = "./version"
	defaultInstanceName   = "initial"
	composeURL            = "https://raw.githubusercontent.com/ptrxyz/chemotion/latest-release/docker-compose.yml"
)

var (
	// configuration
	currentState state
	configFile   string
	// logging
	zlog  zerolog.Logger
	zboth zerolog.Logger
	fs    afero.Fs = afero.NewOsFs() // pointer to the filesystem
)

// chemotionCmd represents the base command when called without any subcommands
var chemotionCmd = &cobra.Command{
	Use:       chemotionValues.use,
	Short:     chemotionValues.short,
	Long:      chemotionValues.long,
	ValidArgs: chemotionValues.options,
	// The following lines are the action associated with a bare application run i.e. without any arguments
	Run: func(cmd *cobra.Command, args []string) {
		zlog.Debug().Msgf("Used command: %s", cmd.CommandPath())
		confirmInteractive()
		zlog.Debug().Msgf("Selection Menu: Chemotion")
		if currentState.isInside {
			fmt.Printf("Welcome to %s! You are inside an instance called %s.\n", projectName, currentState.name)
		} else {
			fmt.Printf("Welcome to %s! You are on a host machine. Currently chosen instance, to manage, is %s.\n", projectName, currentState.name)
		}
		selected := selectOpt(chemotionValues.options)
		switch selected {
		case "system":
			systemCmd.Run(&cobra.Command{}, []string{})
			// case "instance":
			// 	instanceCmd.Run(&cobra.Command{}, []string{})
		}
		zlog.Debug().Msgf("Selected option: %s", selected)
	},
}

// This is called by main.main(). It only needs to happen once.
func Execute() {
	if err := chemotionCmd.Execute(); err == nil {
		zlog.Info().Msgf("%s exited gracefully.", projectName)
	} else {
		zboth.Fatal().Msgf("%s exited abruptly, check log file if necessary. ABORT!", projectName)
	}
}

func init() {
	// begin by setting up logging
	initLog() // in logger.go
	// flag 0: isInside, determined automatically whenever CLI runs
	currentState.isInside = existingFile(stateFPath)
	// initialize viper
	cobra.OnInitialize(initViper) // in configure.go, also takes care of first run, since its need is determined by existence of the configuration file
	// initialize flags
	zlog.Debug().Msg("Start: init()")
	// flag 1: instance, i.e. name of the instance to operate upon
	// terminal overrides config-file, default is `default`
	chemotionCmd.PersistentFlags().StringVarP(&currentState.name, "instance", "i", defaultInstanceName, "start "+projectName+" with a pre-selected instance")
	// flag 2: config, the configuration file
	// config as a flag cannot be read from the configuration file because that creates a circular dependency, default name is hard-coded
	chemotionCmd.PersistentFlags().StringVarP(&configFile, "config", "f", defaultConfigFileName, "path to the configuration file.")
	// flag 3: quiet, i.e. should the CLI run in interactive mode
	// terminal overrides config-file, default is false
	chemotionCmd.PersistentFlags().BoolVarP(&currentState.Quiet, "quiet", "q", false, "use "+projectName+" in scripted mode i.e. without an interactive prompt")
	// flag 4: debug, i.e. should debug messages be logged
	// terminal overrides config-file, default is false
	chemotionCmd.PersistentFlags().BoolVar(&currentState.Debug, "debug", false, "enable logging of debug messages")
	// viper bindings, one for each value in the struct called currentState
	viper.BindPFlag("chosen", chemotionCmd.PersistentFlags().Lookup("instance"))
	viper.BindPFlag(currentState.name+".quiet", chemotionCmd.PersistentFlags().Lookup("quiet"))
	viper.BindPFlag(currentState.name+".debug", chemotionCmd.PersistentFlags().Lookup("debug"))
	viper.Set(currentState.name+".Kind", &currentState.Kind)
	zlog.Debug().Msg("End: init()")
}
