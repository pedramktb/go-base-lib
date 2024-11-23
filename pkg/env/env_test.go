package env

import (
	"os"
	"testing"
)

func Test_GetOrFail(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		expected any
	}{
		{
			name:     "rune",
			key:      "RUNE",
			value:    "1",
			expected: rune(1),
		},
		{
			name:     "byte",
			key:      "BYTE",
			value:    "1",
			expected: byte(1),
		},
		{
			name:     "string",
			key:      "STRING",
			value:    "string",
			expected: string("string"),
		},
		{
			name:     "bool",
			key:      "BOOL",
			value:    "true",
			expected: bool(true),
		},
		{
			name:     "uintptr",
			key:      "UINTPTR",
			value:    "1",
			expected: uintptr(1),
		},
		{
			name:     "int",
			key:      "INT",
			value:    "-1",
			expected: int(-1),
		},
		{
			name:     "int64",
			key:      "INT64",
			value:    "-1",
			expected: int64(-1),
		},
		{
			name:     "int32",
			key:      "INT32",
			value:    "-1",
			expected: int32(-1),
		},
		{
			name:     "int16",
			key:      "INT16",
			value:    "-1",
			expected: int16(-1),
		},
		{
			name:     "int8",
			key:      "INT8",
			value:    "-1",
			expected: int8(-1),
		},
		{
			name:     "uint",
			key:      "UINT",
			value:    "1",
			expected: uint(1),
		},
		{
			name:     "uint64",
			key:      "UINT64",
			value:    "1",
			expected: uint64(1),
		},
		{
			name:     "uint32",
			key:      "UINT32",
			value:    "1",
			expected: uint32(1),
		},
		{
			name:     "uint16",
			key:      "UINT16",
			value:    "1",
			expected: uint16(1),
		},
		{
			name:     "uint8",
			key:      "UINT8",
			value:    "1",
			expected: uint8(1),
		},
		{
			name:     "float64",
			key:      "FLOAT64",
			value:    "1.1",
			expected: float64(1.1),
		},
		{
			name:     "float32",
			key:      "FLOAT32",
			value:    "1.1",
			expected: float32(1.1),
		},
		{
			name:     "complex128",
			key:      "COMPLEX128",
			value:    "1.1+2.2i",
			expected: complex128(complex(float64(1.1), float64(2.2))),
		},
		{
			name:     "complex64",
			key:      "COMPLEX64",
			value:    "1.1+2.2i",
			expected: complex64(complex(float32(1.1), float32(2.2))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.key, tt.value)
			var got any
			switch tt.expected.(type) {
			case string:
				got = GetOrFail[string](tt.key)
			case bool:
				got = GetOrFail[bool](tt.key)
			case uintptr:
				got = GetOrFail[uintptr](tt.key)
			case int:
				got = GetOrFail[int](tt.key)
			case int64:
				got = GetOrFail[int64](tt.key)
			case int32:
				got = GetOrFail[int32](tt.key)
			case int16:
				got = GetOrFail[int16](tt.key)
			case int8:
				got = GetOrFail[int8](tt.key)
			case uint:
				got = GetOrFail[uint](tt.key)
			case uint64:
				got = GetOrFail[uint64](tt.key)
			case uint32:
				got = GetOrFail[uint32](tt.key)
			case uint16:
				got = GetOrFail[uint16](tt.key)
			case uint8:
				got = GetOrFail[uint8](tt.key)
			case float64:
				got = GetOrFail[float64](tt.key)
			case float32:
				got = GetOrFail[float32](tt.key)
			case complex128:
				got = GetOrFail[complex128](tt.key)
			case complex64:
				got = GetOrFail[complex64](tt.key)
			default:
				t.Fatalf("unsupported type: %T", tt.expected)
			}
			if got != tt.expected {
				t.Errorf("got %v, expected %v", got, tt.expected)
			}
		})
	}
}
