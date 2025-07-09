package env

import (
	"os"

	"github.com/cqfn/refrax/internal/log"
	"github.com/joho/godotenv"
)

func Token(path, provider string) string {
	switch provider {
	case "deepseek":
		return find(path, "DEEPSEEK_TOKEN")
	case "mock":
		return find(path, "MOCK_TOKEN")
	default:
		return find(path, "TOKEN")
	}
}

func find(path, variable string) string {
	envs, err := godotenv.Read(path)
	if err != nil {
		log.Info(".env file not found at %s, using default environment settings and parameters", path)
	}
	res := envs[variable]
	if res == "" {
		res = os.Getenv(variable)
	}
	return res
}
