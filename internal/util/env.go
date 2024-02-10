package util

import (
	"os"

	"petrichormud.com/app/internal/constant"
)

func IsProd() bool {
	return os.Getenv(constant.AppEnv) == "true"
}
