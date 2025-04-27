package taggederror

import (
	"errors"
	"net/http"

	"github.com/go-faster/jx"
)

// Handler handles errors for http.Handler
func Handler(err error, w http.ResponseWriter, r *http.Request) {
	var taggedErr *Error
	if !errors.As(err, &taggedErr) {
		taggedErr = ErrInternal.Wrap(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(taggedErr.Code())

	e := jx.GetEncoder()
	e.ObjStart()
	e.FieldStart("code")
	e.Int(taggedErr.Code())

	e.FieldStart("tag")
	e.StrEscape(taggedErr.Tag())

	e.FieldStart("detail")
	e.StrEscape(taggedErr.Error())

	e.ObjEnd()

	_, _ = e.WriteTo(w)
}
