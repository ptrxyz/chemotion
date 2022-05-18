package cmd

import (
	"fmt"
	"os"

	"github.com/cavaliergopher/grab/v3"
)

func instanceCreate(current state) (success bool) {
	zboth.Info().Msgf("Creating a new default (production) instance of %s called %s.", projectName, currentState.name)
	// prerequisites: docker
	confirmDocker()
	fs.MkdirAll("instances/"+currentState.name, folderPerm) // make folder in case it doesn't exist
	os.Chdir("instances/" + currentState.name)
	zlog.Debug().Str("instance", currentState.name).Msgf("Changed working directory to: instances/%s", currentState.name)
	var composeFilename string
	if resp, err := grab.Get("x", composeURL); err == nil {
		composeFilename = resp.Filename
		fmt.Println("problem with file " + composeFilename)
		composeFilename = "docker-compose.yml"
		zlog.Info().Str("instance", currentState.name).Msgf("Downloaded docker-compose file saved as: %s", composeFilename)
	} else {
		zboth.Fatal().Err(err).Msg("Failed to download docker-compose file. Check log. ABORT!")
	}
	commandStr := fmt.Sprintf("compose -f %s create", composeFilename)
	zlog.Info().Str("instance", currentState.name).Msgf("Starting docker with command: %s", commandStr)
	if success = callDocker(commandStr); !success {
		zboth.Fatal().Err(fmt.Errorf("docker compose create failed")).Msgf("Failed to setup %. Check log. ABORT!", projectName)
	}
	os.Chdir("../..")
	return success
}
