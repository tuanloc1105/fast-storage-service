package config

import (
	"os"
)

func MountFolderLocationConfig() bool {
	_, isMountFolderFromEnvSet := os.LookupEnv("MOUNT_FOLDER")

	return isMountFolderFromEnvSet
}
