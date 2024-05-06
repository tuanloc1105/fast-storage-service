package config

import (
	"fast-storage-go-service/constant"
	"os"
)

func MountFolderLocationConfig() bool {
	mountFolderFromEnv, isMountFolderFromEnvSet := os.LookupEnv("MOUNT_FOLDER")

	constant.MountFolder = mountFolderFromEnv

	return isMountFolderFromEnvSet
}
