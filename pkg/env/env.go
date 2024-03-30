package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load("local.env")
}

func IsDebug() bool {
	return os.Getenv("DEBUG") == "true"
}

func GetStrWithFallback(env string, fallback string) string {
	value := os.Getenv(env)
	if value == "" {
		value = fallback
	}
	return value
}

func GetStrOrFail(env string) string {
	value := os.Getenv(env)
	if value == "" {
		panic("missing env variable: " + env)
	}
	return value
}

func GetIntWithFallback(env string, fallback int) int {
	value := os.Getenv(env)
	if value == "" {
		return fallback
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		panic("invalid int: " + env)
	}

	return intVal
}

func GetIntOrFail(env string) int {
	value := os.Getenv(env)
	if value == "" {
		panic("missing env variable: " + env)
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		panic("invalid int: " + env)
	}

	return intVal
}

func GetUint64WithFallback(env string, fallback uint64) uint64 {
	value := os.Getenv(env)
	if value == "" {
		return fallback
	}

	intVal, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		panic("invalid uint64: " + env)
	}

	return intVal
}

func GetUint64OrFail(env string) uint64 {
	value := os.Getenv(env)
	if value == "" {
		panic("missing env variable: " + env)
	}

	intVal, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		panic("invalid uint64: " + env)
	}

	return intVal
}

func GetUint32WithFallback(env string, fallback uint32) uint32 {
	value := os.Getenv(env)
	if value == "" {
		return fallback
	}

	intVal, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		panic("invalid uint32: " + env)
	}

	return uint32(intVal)
}

func GetUint32OrFail(env string) uint32 {
	value := os.Getenv(env)
	if value == "" {
		panic("missing env variable: " + env)
	}

	intVal, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		panic("invalid uint32: " + env)
	}

	return uint32(intVal)
}

func GetUint16WithFallback(env string, fallback uint16) uint16 {
	value := os.Getenv(env)
	if value == "" {
		return fallback
	}

	intVal, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		panic("invalid uint16: " + env)
	}

	return uint16(intVal)
}

func GetUint16OrFail(env string) uint16 {
	value := os.Getenv(env)
	if value == "" {
		panic("missing env variable: " + env)
	}

	intVal, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		panic("invalid uint16: " + env)
	}

	return uint16(intVal)
}
