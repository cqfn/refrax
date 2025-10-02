package env

import (
	"os"

	"github.com/cqfn/refrax/internal/log"
	"github.com/joho/godotenv"
)

// Token retrieves the token for the specified provider from the given .env file path.
// If the token is not found in the .env file, it falls back to the system environment variables.
func Token(path, provider string) string {
	switch provider {
	case "deepseek":
		return find(path, "DEEPSEEK_API_KEY")
	case "openai":
		return find(path, "OPENAI_API_KEY")
	case "mock":
		return find(path, "MOCK_TOKEN")
	case "ollama":
		return find(path, "OLLAMA_TOKEN")
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
