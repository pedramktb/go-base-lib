package env

import (
	"os"

	"github.com/joho/godotenv"
)

const StageEnv = "STAGE"

func init() {
	if IsLocal() {
		err := godotenv.Load("local.env")
		if err != nil {
			panic(err)
		}
	}
}

func GetEnvWithFallback(env string, fallback string) string {
	value := os.Getenv(env)
	if value == "" {
		value = fallback
	}
	return value
}

func GetEnvOrFail(env string) string {
	value := os.Getenv(env)
	if value == "" {
		panic("missing env variable: " + env)
	}
	return value
}

func IsProd() bool {
	return os.Getenv(StageEnv) == "prod"
}

func IsLocal() bool {
	stageEnv := os.Getenv(StageEnv)
	return stageEnv == "" || stageEnv == "local"
}
