package evasion

import (
	"errors"
	"net/http"

	"github.com/pedramktb/go-base-lib/env"
	"github.com/pedramktb/go-base-lib/taggederror"
)

// evasiveError makes sure that all errors are sent with status code 200
var evasiveError = taggederror.NewRoot(
	errors.New("error"),
	"ERROR",
	http.StatusOK,
)

func ErrorHandler(err error, trusted bool, w http.ResponseWriter, r *http.Request) {
	if trusted {
		trustedHandler(err, w, r)
	} else {
		untrustedHandler(err, w, r)
	}
}

// trustedHandler uses evasiveError to handle errors
func trustedHandler(err error, w http.ResponseWriter, r *http.Request) {
	var taggedErr *taggederror.Error
	if !errors.As(err, &taggedErr) {
		taggedErr = taggederror.ErrInternal.Wrap(err)
	}
	taggederror.Handler(evasiveError.Wrap(err), w, r)
}

// untrustedHandler uses FailStatusCode with no detail
func untrustedHandler(err error, w http.ResponseWriter, r *http.Request) {
	if env.GetEnvironment() != env.EnvironmentProd {
		taggederror.Handler(err, w, r)
	} else {
		w.WriteHeader(FailStatusCode)
	}
}
