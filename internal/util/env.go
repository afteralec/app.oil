package util

import (
	"os"

	"petrichormud.com/app/internal/constants"
)

func IsProd() bool {
	return os.Getenv(constants.AppEnv) == "true"
}
