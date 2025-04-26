package types

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_Nillable_JSON_Marhshal(t *testing.T) {
	randomID := uuid.New()
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{
			name:  "string",
			value: NewNillable(Pointer("example")),
			want:  `"example"`,
		},
		{
			name:  "null",
			value: NewNillable[any](nil),
			want:  `null`,
		},
		{
			name:  "object",
			value: NewNillable(&struct{ ID uuid.UUID }{ID: randomID}),
			want:  `{"ID":"` + randomID.String() + `"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.value)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func Test_Nillable_JSON_Unmarhshal(t *testing.T) {
	randomID := uuid.New()
	tests := []struct {
		name  string
		value string
		want  any
	}{
		{
			name:  "nil",
			value: `null`,
			want:  NewNillable[any](nil),
		},
		{
			name:  "string",
			value: `"example"`,
			want:  NewNillable(Pointer("example")),
		},
		{
			name:  "object",
			value: `{"ID":"` + randomID.String() + `"}`,
			want:  NewNillable(&struct{ ID uuid.UUID }{ID: randomID}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.want.(type) {
			case Nillable[any]:
				got := tt.want.(Nillable[any])
				err := json.Unmarshal([]byte(tt.value), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			case Nillable[string]:
				got := tt.want.(Nillable[string])
				err := json.Unmarshal([]byte(tt.value), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			case Nillable[struct{ ID uuid.UUID }]:
				got := tt.want.(Nillable[struct{ ID uuid.UUID }])
				err := json.Unmarshal([]byte(tt.value), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_Nillable_BSON_Marhshal(t *testing.T) {
	randomID := uuid.New()
	randomIDBSON, _ := bson.Marshal(struct{ ID uuid.UUID }{ID: randomID})
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{
			name:  "string",
			value: NewNillable(Pointer("example")),
			want:  "\b\x00\x00\x00example\x00",
		},
		{
			name:  "null",
			value: NewNillable[any](nil),
			want:  "",
		},
		{
			name:  "object",
			value: NewNillable(&struct{ ID uuid.UUID }{ID: randomID}),
			want:  string(randomIDBSON),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got, err := bson.MarshalValue(tt.value)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func Test_Nillable_BSON_Unmarhshal(t *testing.T) {
	randomID := uuid.New()
	randomIDBSON, _ := bson.Marshal(struct{ ID uuid.UUID }{ID: randomID})
	tests := []struct {
		name  string
		value string
		want  any
	}{
		{
			name:  "null",
			value: "",
			want:  NewNillable[any](nil),
		},
		{
			name:  "string",
			value: "\b\x00\x00\x00example\x00",
			want:  NewNillable(Pointer("example")),
		},
		{
			name:  "object",
			value: string(randomIDBSON),
			want:  NewNillable(&struct{ ID uuid.UUID }{ID: randomID}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.want.(type) {
			case Nillable[any]:
				got := Nillable[any]{}
				err := bson.UnmarshalValue(bson.TypeNull, []byte(tt.value), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			case Nillable[string]:
				got := Nillable[string]{}
				err := bson.UnmarshalValue(bson.TypeString, []byte(tt.value), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			case Nillable[struct{ ID uuid.UUID }]:
				got := Nillable[struct{ ID uuid.UUID }]{}
				err := bson.UnmarshalValue(bson.TypeEmbeddedDocument, []byte(tt.value), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_Nillable_Value(t *testing.T) {
	randomID := uuid.New()
	tests := []struct {
		name  string
		value any
		want  driver.Value
	}{
		{
			name:  "nil",
			value: NewNillable[any](nil),
			want:  nil,
		},
		{
			name:  "string",
			value: NewNillable(Pointer("example")),
			want:  "example",
		},
		{
			name:  "uuid",
			value: NewNillable(&randomID),
			want:  randomID.String(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.value.(driver.Valuer).Value()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Nillable_Scan(t *testing.T) {
	randomID := uuid.New()
	tests := []struct {
		name  string
		value any
		want  any
	}{
		{
			name:  "null",
			value: nil,
			want:  NewNillable[any](nil),
		},
		{
			name:  "string",
			value: "example",
			want:  NewNillable(Pointer("example")),
		},
		{
			name:  "object",
			value: randomID[:],
			want:  NewNillable(&randomID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.want.(type) {
			case Nillable[any]:
				got := Nillable[any]{}
				err := got.Scan(tt.value)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			case Nillable[string]:
				got := Nillable[string]{}
				err := got.Scan(tt.value)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			case Nillable[uuid.UUID]:
				got := Nillable[uuid.UUID]{}
				err := got.Scan(tt.value)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
