package entity

import (
	"encoding/json"
	"errors"

	"github.com/pedramktb/go-base-lib/taggederror"
	"github.com/pedramktb/go-base-lib/types"
)

var (
	ErrInvalidCursor = taggederror.ErrBadRequest.Wrap(taggederror.New(
		errors.New("invalid cursor"),
		"INVALID_PAGINATION_CURSOR",
	))
)

// Cursor is the list of values in the same order of sorts
type Cursor []any

// Mapper Implementation of Cursor to DTO (string)
func (c Cursor) ToDTO() (types.Nillable[string], error) {
	str := types.Nillable[string]{}
	for i := range c {
		if i == 0 {
			str.NotNil = true
		} else {
			str.Val += ","
		}
		partBytes, err := json.Marshal(c[i])
		if err != nil {
			return types.Nillable[string]{}, err
		}
		str.Val += string(partBytes)
	}
	return str, nil
}
