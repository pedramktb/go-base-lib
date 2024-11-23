package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load(".env")
}

func IsDebug() bool {
	return GetEnvWithFallback("DEBUG", false)
}

func GetEnvOrFail[
	T string |
		bool |
		uintptr |
		int | int64 | int32 | int16 | int8 |
		uint | uint64 | uint32 | uint16 | uint8 |
		float64 | float32 |
		complex128 | complex64](env string) T {
	str := os.Getenv(env)
	if str == "" {
		panic("missing env variable: " + env)
	}

	var val any

	switch any(*new(T)).(type) {
	case string:
		val = str
	case bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			panic("invalid bool: " + env)
		}
		val = b
	case uintptr:
		u64, err := strconv.ParseUint(str, 10, 0)
		if err != nil {
			panic("invalid uintptr: " + env)
		}
		val = uintptr(u64)
	case int:
		i64, err := strconv.ParseInt(str, 10, 0)
		if err != nil {
			panic("invalid int: " + env)
		}
		val = int(i64)
	case int64:
		i64, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			panic("invalid int64: " + env)
		}
		val = i64
	case int32:
		i64, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			panic("invalid int32: " + env)
		}
		val = int32(i64)
	case int16:
		i64, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			panic("invalid int16: " + env)
		}
		val = int16(i64)
	case int8:
		i64, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			panic("invalid int8: " + env)
		}
		val = int8(i64)
	case uint:
		i64, err := strconv.ParseUint(str, 10, 0)
		if err != nil {
			panic("invalid uint: " + env)
		}
		val = uint(i64)
	case uint64:
		i64, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			panic("invalid uint64: " + env)
		}
		val = i64
	case uint32:
		i64, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			panic("invalid uint32: " + env)
		}
		val = uint32(i64)
	case uint16:
		i64, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			panic("invalid uint16: " + env)
		}
		val = uint16(i64)
	case uint8:
		i64, err := strconv.ParseUint(str, 10, 8)
		if err != nil {
			panic("invalid uint8: " + env)
		}
		val = uint8(i64)
	case float64:
		f64, err := strconv.ParseFloat(str, 64)
		if err != nil {
			panic("invalid float64: " + env)
		}
		val = f64
	case float32:
		f64, err := strconv.ParseFloat(str, 32)
		if err != nil {
			panic("invalid float32: " + env)
		}
		val = float32(f64)
	case complex128:
		c128, err := strconv.ParseComplex(str, 128)
		if err != nil {
			panic("invalid complex128: " + env)
		}
		val = c128
	case complex64:
		c128, err := strconv.ParseComplex(str, 64)
		if err != nil {
			panic("invalid complex64: " + env)
		}
		val = complex64(c128)
	default:
		panic("unsupported type")
	}

	return val.(T)
}

func GetEnvWithFallback[
	T string |
		bool |
		uintptr |
		int | int64 | int32 | int16 | int8 |
		uint | uint64 | uint32 | uint16 | uint8 |
		float64 | float32 |
		complex128 | complex64](env string, fallback T) T {
	str := os.Getenv(env)
	if env == "" {
		return fallback
	}

	var val any

	switch any(*new(T)).(type) {
	case string:
		val = str
	case bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			panic("invalid bool: " + env)
		}
		val = b
	case uintptr:
		u64, err := strconv.ParseUint(str, 10, 0)
		if err != nil {
			panic("invalid uintptr: " + env)
		}
		val = uintptr(u64)
	case int:
		i64, err := strconv.ParseInt(str, 10, 0)
		if err != nil {
			panic("invalid int: " + env)
		}
		val = int(i64)
	case int64:
		i64, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			panic("invalid int64: " + env)
		}
		val = i64
	case int32:
		i64, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			panic("invalid int32: " + env)
		}
		val = int32(i64)
	case int16:
		i64, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			panic("invalid int16: " + env)
		}
		val = int16(i64)
	case int8:
		i64, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			panic("invalid int8: " + env)
		}
		val = int8(i64)
	case uint:
		i64, err := strconv.ParseUint(str, 10, 0)
		if err != nil {
			panic("invalid uint: " + env)
		}
		val = uint(i64)
	case uint64:
		i64, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			panic("invalid uint64: " + env)
		}
		val = i64
	case uint32:
		i64, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			panic("invalid uint32: " + env)
		}
		val = uint32(i64)
	case uint16:
		i64, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			panic("invalid uint16: " + env)
		}
		val = uint16(i64)
	case uint8:
		i64, err := strconv.ParseUint(str, 10, 8)
		if err != nil {
			panic("invalid uint8: " + env)
		}
		val = uint8(i64)
	case float64:
		f64, err := strconv.ParseFloat(str, 64)
		if err != nil {
			panic("invalid float64: " + env)
		}
		val = f64
	case float32:
		f64, err := strconv.ParseFloat(str, 32)
		if err != nil {
			panic("invalid float32: " + env)
		}
		val = float32(f64)
	case complex128:
		c128, err := strconv.ParseComplex(str, 128)
		if err != nil {
			panic("invalid complex128: " + env)
		}
		val = c128
	case complex64:
		c128, err := strconv.ParseComplex(str, 64)
		if err != nil {
			panic("invalid complex64: " + env)
		}
		val = complex64(c128)
	default:
		panic("unsupported type")
	}

	return val.(T)
}
