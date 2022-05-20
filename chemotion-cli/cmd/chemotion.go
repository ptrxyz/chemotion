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

	"github.com/chigopher/pathlib"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	versionCLI            = "0.1"
	versionYAML           = "1.0"
	nameCLI               = "Chemotion"
	defaultConfigFilepath = "./chemotion-cli.yml"
	logFileName           = "chemotion-cli.log"
	instanceDefault       = "initial"
	stateFile             = "./version"
	composeURL            = "https://raw.githubusercontent.com/ptrxyz/chemotion/1.1.2p220401/new-build/release/1.1.2p220401/docker-compose.yml"
	shell                 = "bash"
	virtualizer           = "docker"
	instancesFolder       = "instances" // the folder in which chemotion expects to find all the instances
	//composeURL            = "https://raw.githubusercontent.com/ptrxyz/chemotion/latest-release/docker-compose.yml"
)

var (
	// configuration
	currentState state
	configFile   string
	firstRun     bool
	// logging
	zlog  zerolog.Logger
	zboth zerolog.Logger
	// path of the working directory
	workDir pathlib.Path = *pathlib.NewPath(".")
	// minimum version for virtualizer (docker)
	minimumVirtualizer = "17.12" // so as to support docker compose files version 3.5
)

// struct to store information about the currently selected instance, which has implications for the current state of the program
type state struct {
	Debug    bool
	Quiet    bool
	Kind     string
	name     string
	isInside bool
}

// cmd represents the base command when called without any subcommands
var cmd = &cobra.Command{
	Use:     "chemotion",
	Short:   "CLI for Chemotion ELN",
	Long:    "Chemotion ELN is an Electronic Lab Notebook solution.\nDeveloped for researchers, the software aims to work for you.\nSee, https://www.chemotion.net.",
	Version: versionCLI,
	// The following lines are the action associated with a bare application run i.e. without any arguments
	Run: func(cmd *cobra.Command, args []string) {
		zlog.Debug().Msgf("Used command: %s", cmd.CalledAs())
		confirmInteractive()
		if firstRun {
			zboth.Info().Msgf("Welcome to your first run of %s.", nameCLI)
			if selectYesNo("Have you installed the prerequisites, particularly "+virtualizer, false) {
				confirmVirtualizer(minimumVirtualizer)
			} else {
				zboth.Info().Msgf("Please install %s (at least version %s) before proceeding. Thanks!", virtualizer, minimumVirtualizer)
			}
		} else {
			fmt.Printf("Welcome to %s! You are on a host machine. The instance you are currently managing is %s%s%s%s.\n", nameCLI, string("\033[31m"), string("\033[1m"), currentState.name, string("\033[0m"))
			zlog.Debug().Msgf("Selection Menu: Chemotion")
			acceptedOpts := []string{"on", "off ", "restart", "instance", "user", "system", "exit"}
			if firstRun {
				acceptedOpts = []string{"install", "system", "exit"}
			}
			selected := selectOpt(acceptedOpts)
			zlog.Debug().Msgf("Selected option: %s", selected)
			switch selected {
			case "install":
				installCmd.Run(&cobra.Command{}, []string{})
			case "system":
				systemCmd.Run(&cobra.Command{}, []string{})
			case "instance":
				instanceCmd.Run(&cobra.Command{}, []string{})
			case "on":
				// onCmd.Run(&cobra.Command{}, []string{})
			case "off":
				// offCmd.Run(&cobra.Command{}, []string{})
			case "restart":
				// offCmd.Run(&cobra.Command{}, []string{})
				// onCmd.Run(&cobra.Command{}, []string{})
			case "exit":
				zlog.Debug().Msg("Chose to exit.")
			}
		}
	},
}

// This is called by main.main(). It only needs to happen once.
func Execute() {
	if err := cmd.Execute(); err == nil {
		zlog.Debug().Msgf("%s exited gracefully.", nameCLI)
	} else {
		zboth.Fatal().Err(fmt.Errorf("unexplained")).Msgf("%s exited abruptly, check log file if necessary. ABORT!", nameCLI)
	}
}

func init() {
	// flag 0: isInside and firstRun, determined automatically whenever CLI runs
	currentState.isInside = existingFile(stateFile)
	firstRun = !existingFile(configFile)
	// begin by setting up logging
	initLog() // in logger.go
	// initialize flags
	zlog.Debug().Msg("Start: init(): initialize flags")
	// flag 1: instance, i.e. name of the instance to operate upon
	// terminal overrides config-file, default is `default`
	if !firstRun {
		cmd.PersistentFlags().StringVarP(&currentState.name, "select-instance", "i", "", "start "+nameCLI+" with a pre-selected instance")
	}
	// flag 2: config, the configuration file
	// config as a flag cannot be read from the configuration file because that creates a circular dependency, default name is hard-coded
	cmd.PersistentFlags().StringVarP(&configFile, "config", "f", defaultConfigFilepath, "path to the configuration file.")
	// flag 3: quiet, i.e. should the CLI run in interactive mode
	// terminal overrides config-file, default is false
	cmd.PersistentFlags().BoolVarP(&currentState.Quiet, "quiet", "q", false, "use "+nameCLI+" in scripted mode i.e. without an interactive prompt")
	// flag 4: debug, i.e. should debug messages be logged
	// terminal overrides config-file, default is false
	cmd.PersistentFlags().BoolVar(&currentState.Debug, "debug", false, "enable logging of debug messages")
	zlog.Debug().Msg("End: init(): initialize flags")
	zlog.Debug().Msg("Start: init(): bind flags")
	// viper bindings, one for each value in the struct called currentState
	viper.BindPFlag("selected", cmd.PersistentFlags().Lookup("select-instance"))
	viper.BindPFlag("instances."+currentState.name+".quiet", cmd.PersistentFlags().Lookup("quiet"))
	viper.BindPFlag("instances."+currentState.name+".debug", cmd.PersistentFlags().Lookup("debug"))
	viper.Set("version", versionYAML)
	zlog.Debug().Msg("End: init(): bind flags")
	zlog.Debug().Msg("Start: init(): add commands")
	// add commands to `chemotion`
	if firstRun {
		cmd.AddCommand(installCmd)
		cmd.AddCommand(systemCmd)
	} else {
		cmd.AddCommand(instanceCmd)
		cmd.AddCommand(systemCmd)
	}
	zlog.Debug().Msg("End: init(): add commands")
	// initialize viper (runs last, i.e. when cmd.Execute runs)
	cobra.OnInitialize(initViper) // in configure.go, also takes care of first run, since its need is determined by existence of the configuration file
}
