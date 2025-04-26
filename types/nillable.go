package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Nillable is a type that can be used to represent a nillable/nullable value.
// The flag NotNil is used because certain Unmarshalers return the zero value
// of the type when facing a null value (e.g. BSON Unmarshaler)
type Nillable[T any] struct {
	Val    T
	NotNil bool
}

func NewNillable[T any](value *T) Nillable[T] {
	if value == nil {
		return Nillable[T]{}
	} else {
		return Nillable[T]{Val: *value, NotNil: true}
	}
}

func (n Nillable[T]) ToType() T {
	if n.NotNil {
		return n.Val
	}
	return *new(T)
}

func (n Nillable[T]) MarshalJSON() ([]byte, error) {
	if n.NotNil {
		return json.Marshal(n.Val)
	} else {
		return []byte("null"), nil
	}
}

func (n *Nillable[T]) UnmarshalJSON(data []byte) error {
	var t *T
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}
	if t != nil {
		n.Val = *t
		n.NotNil = true
	}
	return nil
}

func (n Nillable[T]) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if n.NotNil {
		return bson.MarshalValue(n.Val)
	}
	var t *T
	return bson.MarshalValue(t)
}

func (n *Nillable[T]) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t == bson.TypeNull {
		return nil
	}

	err := bson.UnmarshalValue(t, data, &n.Val)
	if err != nil {
		return err
	}
	n.NotNil = true
	return nil
}

func (n Nillable[T]) Value() (driver.Value, error) {
	if n.NotNil {
		if valuer, ok := any(n.Val).(driver.Valuer); ok {
			return valuer.Value()
		} else {
			return driver.Value(n.Val), nil
		}
	}
	return driver.Value(nil), nil
}

func (n *Nillable[T]) Scan(src any) error {
	if src == nil {
		return nil
	}
	if scanner, ok := any(&n.Val).(sql.Scanner); ok {
		err := scanner.Scan(src)
		if err != nil {
			return err
		}
		n.NotNil = true
	} else if n.Val, ok = src.(T); ok {
		n.NotNil = true
	} else {
		return fmt.Errorf("cannot scan %v into Nillable[%T]", src, n.Val)
	}
	return nil
}
