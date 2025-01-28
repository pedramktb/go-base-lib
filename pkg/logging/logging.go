package logging

import (
	"context"
	"time"

	"github.com/pedramktb/go-base-lib/pkg/env"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	if env.IsDebug() {
		logger = zap.Must(zap.NewDevelopment())
	} else {
		logger = zap.Must(zap.NewProduction())
	}
}

func Logger() *zap.Logger {
	return logger
}

var contextKeys = []string{}

func AddContextKeys(keys []string) {
	contextKeys = append(contextKeys, keys...)
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return Logger()
	}
	return Logger().With(mapToZapFields(ctxToMap(ctx))...)
}

func ctxToMap(ctx context.Context) map[string]any {
	m := make(map[string]any)
	for _, key := range contextKeys {
		if value := ctx.Value(key); value != nil {
			m[key] = value
		}
	}
	return m
}

func mapToZapFields(m map[string]any) []zap.Field {
	var fields []zap.Field
	for key, value := range m {
		strKey := string(key)
		switch v := value.(type) {
		case string:
			fields = append(fields, zap.String(strKey, v))
		case bool:
			fields = append(fields, zap.Bool(strKey, v))
		case uintptr:
			fields = append(fields, zap.Uintptr(strKey, v))
		case int:
			fields = append(fields, zap.Int(strKey, v))
		case int64:
			fields = append(fields, zap.Int64(strKey, v))
		case int32:
			fields = append(fields, zap.Int32(strKey, v))
		case int16:
			fields = append(fields, zap.Int16(strKey, v))
		case int8:
			fields = append(fields, zap.Int8(strKey, v))
		case uint:
			fields = append(fields, zap.Uint(strKey, v))
		case uint64:
			fields = append(fields, zap.Uint64(strKey, v))
		case uint32:
			fields = append(fields, zap.Uint32(strKey, v))
		case uint16:
			fields = append(fields, zap.Uint16(strKey, v))
		case uint8:
			fields = append(fields, zap.Uint8(strKey, v))
		case float64:
			fields = append(fields, zap.Float64(strKey, v))
		case float32:
			fields = append(fields, zap.Float32(strKey, v))
		case complex128:
			fields = append(fields, zap.Complex128(strKey, v))
		case complex64:
			fields = append(fields, zap.Complex64(strKey, v))
		case time.Time:
			fields = append(fields, zap.Time(strKey, v))
		case error:
			fields = append(fields, zap.NamedError(strKey, v))
		case *string:
			fields = append(fields, zap.Stringp(strKey, v))
		case *bool:
			fields = append(fields, zap.Boolp(strKey, v))
		case *uintptr:
			fields = append(fields, zap.Uintptrp(strKey, v))
		case *int:
			fields = append(fields, zap.Intp(strKey, v))
		case *int64:
			fields = append(fields, zap.Int64p(strKey, v))
		case *int32:
			fields = append(fields, zap.Int32p(strKey, v))
		case *int16:
			fields = append(fields, zap.Int16p(strKey, v))
		case *int8:
			fields = append(fields, zap.Int8p(strKey, v))
		case *uint:
			fields = append(fields, zap.Uintp(strKey, v))
		case *uint64:
			fields = append(fields, zap.Uint64p(strKey, v))
		case *uint32:
			fields = append(fields, zap.Uint32p(strKey, v))
		case *uint16:
			fields = append(fields, zap.Uint16p(strKey, v))
		case *uint8:
			fields = append(fields, zap.Uint8p(strKey, v))
		case *float64:
			fields = append(fields, zap.Float64p(strKey, v))
		case *float32:
			fields = append(fields, zap.Float32p(strKey, v))
		case *complex128:
			fields = append(fields, zap.Complex128p(strKey, v))
		case *complex64:
			fields = append(fields, zap.Complex64p(strKey, v))
		case *time.Time:
			fields = append(fields, zap.Timep(strKey, v))
		case []string:
			fields = append(fields, zap.Strings(strKey, v))
		case []bool:
			fields = append(fields, zap.Bools(strKey, v))
		case []uintptr:
			fields = append(fields, zap.Uintptrs(strKey, v))
		case []int:
			fields = append(fields, zap.Ints(strKey, v))
		case []int64:
			fields = append(fields, zap.Int64s(strKey, v))
		case []int32:
			fields = append(fields, zap.Int32s(strKey, v))
		case []int16:
			fields = append(fields, zap.Int16s(strKey, v))
		case []int8:
			fields = append(fields, zap.Int8s(strKey, v))
		case []uint:
			fields = append(fields, zap.Uints(strKey, v))
		case []uint64:
			fields = append(fields, zap.Uint64s(strKey, v))
		case []uint32:
			fields = append(fields, zap.Uint32s(strKey, v))
		case []uint16:
			fields = append(fields, zap.Uint16s(strKey, v))
		case []uint8:
			fields = append(fields, zap.Uint8s(strKey, v))
		case []float64:
			fields = append(fields, zap.Float64s(strKey, v))
		case []float32:
			fields = append(fields, zap.Float32s(strKey, v))
		case []complex128:
			fields = append(fields, zap.Complex128s(strKey, v))
		case []complex64:
			fields = append(fields, zap.Complex64s(strKey, v))
		case []time.Time:
			fields = append(fields, zap.Times(strKey, v))
		case []error:
			fields = append(fields, zap.Errors(strKey, v))
		case [][]byte:
			fields = append(fields, zap.ByteStrings(strKey, v))
		default:
			fields = append(fields, zap.Any(strKey, value))
		}
	}
	return fields
}
