package env

import (
	"github.com/cqfn/refrax/internal/log"
	"github.com/joho/godotenv"
)

func Token(path string) string {
	envs, err := godotenv.Read(path)
	if err != nil {
		log.Info(".env file not found at %s, using default environment settings and parameters", path)
		return ""
	}
	return envs["TOKEN"]
}
