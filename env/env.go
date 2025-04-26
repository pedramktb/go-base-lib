package env

import (
	"errors"
	"os"
	"strconv"
)

type Environment string

const (
	EnvironmentProd    Environment = "prod"
	EnvironmentStaging Environment = "staging"
	EnvironmentDev     Environment = "dev"
	EnvironmentLocal   Environment = ""
)

func (e Environment) String() string {
	if e == EnvironmentLocal {
		return "local"
	} else {
		return string(e)
	}
}

var env Environment
var envSet bool

func SetEnvironment(e Environment) {
	env = e
	envSet = true
}

func GetEnvironment() Environment {
	if envSet {
		return env
	}
	e, _ := GetWithFallback("ENVIRONMENT", string(EnvironmentLocal))
	return Environment(e)
}

// GetOrFail returns the value of the environment variable or returns an error if it is not set
func GetOrFail[ //NOSONAR
	T string |
		bool |
		uintptr |
		int | int64 | int32 | int16 | int8 |
		uint | uint64 | uint32 | uint16 | uint8 |
		float64 | float32 |
		complex128 | complex64](env string) (T, error) {
	str := os.Getenv(env)
	if str == "" {
		return *new(T), errors.New("missing env variable: " + env)
	}

	var val any

	switch any(*new(T)).(type) {
	case string:
		val = str
	case bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			return *new(T), errors.New("invalid bool: " + env)
		}
		val = b
	case uintptr:
		u64, err := strconv.ParseUint(str, 10, 0)
		if err != nil {
			return *new(T), errors.New("invalid uintptr: " + env)
		}
		val = uintptr(u64)
	case int:
		i64, err := strconv.ParseInt(str, 10, 0)
		if err != nil {
			return *new(T), errors.New("invalid int: " + env)
		}
		val = int(i64)
	case int64:
		i64, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return *new(T), errors.New("invalid int64: " + env)
		}
		val = i64
	case int32:
		i64, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return *new(T), errors.New("invalid int32: " + env)
		}
		val = int32(i64)
	case int16:
		i64, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			return *new(T), errors.New("invalid int16: " + env)
		}
		val = int16(i64)
	case int8:
		i64, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			return *new(T), errors.New("invalid int8: " + env)
		}
		val = int8(i64)
	case uint:
		i64, err := strconv.ParseUint(str, 10, 0)
		if err != nil {
			return *new(T), errors.New("invalid uint: " + env)
		}
		val = uint(i64)
	case uint64:
		i64, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return *new(T), errors.New("invalid uint64: " + env)
		}
		val = i64
	case uint32:
		i64, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return *new(T), errors.New("invalid uint32: " + env)
		}
		val = uint32(i64)
	case uint16:
		i64, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			return *new(T), errors.New("invalid uint16: " + env)
		}
		val = uint16(i64)
	case uint8:
		i64, err := strconv.ParseUint(str, 10, 8)
		if err != nil {
			return *new(T), errors.New("invalid uint8: " + env)
		}
		val = uint8(i64)
	case float64:
		f64, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return *new(T), errors.New("invalid float64: " + env)
		}
		val = f64
	case float32:
		f64, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return *new(T), errors.New("invalid float32: " + env)
		}
		val = float32(f64)
	case complex128:
		c128, err := strconv.ParseComplex(str, 128)
		if err != nil {
			return *new(T), errors.New("invalid complex128: " + env)
		}
		val = c128
	case complex64:
		c128, err := strconv.ParseComplex(str, 64)
		if err != nil {
			return *new(T), errors.New("invalid complex64: " + env)
		}
		val = complex64(c128)
	default:
		return *new(T), errors.New("unsupported type")
	}

	return val.(T), nil
}

// GetWithFallback returns the value of the environment variable or the fallback value if it is not set
func GetWithFallback[
	T string |
		bool |
		uintptr |
		int | int64 | int32 | int16 | int8 |
		uint | uint64 | uint32 | uint16 | uint8 |
		float64 | float32 |
		complex128 | complex64](env string, fallback T) (T, error) {
	str := os.Getenv(env)
	if str == "" {
		return fallback, nil
	}

	return GetOrFail[T](env)
}
