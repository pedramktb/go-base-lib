package entity

import (
	"errors"

	"github.com/pedramktb/go-base-lib/taggederror"
)

var (
	ErrDBUnhandled = taggederror.ErrInternal.Wrap(taggederror.New(
		errors.New("db unhandled error"),
		"UNHANDLED_DATABASE_ERROR",
	))
	ErrUnknownField = taggederror.ErrBadRequest.Wrap(taggederror.New(
		errors.New("unknown field"),
		"UNKNOWN_FIELD",
	))
)
