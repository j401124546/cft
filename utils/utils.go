package utils

import (
	"cft/config"
	"fmt"
)

func GetCheckpointFile(id string) string {
	return fmt.Sprintf("%s/%s.tar.gz", config.GetConfig().CheckpointDir, id)
}
