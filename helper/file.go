package helper

import (
	"fmt"
	"os"
	"github.com/tfonfara/plexsmarthome/constants"
)

func getDataPath() string {
	dataPath := os.Getenv("PLEX_CONFIG")
	if len(dataPath) == 0 {
		dataPath = constants.DefaultConfigDir
	}
	return dataPath
}

func ConfigurationFilePath() string {
	return fmt.Sprintf("%s/%s", getDataPath(), constants.ConfigFileName)
}
