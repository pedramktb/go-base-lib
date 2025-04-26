package evasion

import (
	"math/rand"
	"net/http"
	"time"
)

var statusCodes = []int{
	http.StatusBadRequest,
	http.StatusUnauthorized,
	http.StatusForbidden,
	http.StatusNotFound,

	http.StatusInternalServerError,
	http.StatusNotImplemented,
	http.StatusBadGateway,
	http.StatusServiceUnavailable,
}

// FailStatusCode is a random failed status code initialized at the start of the program (401-404, 500-503)
var FailStatusCode = statusCodes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(statusCodes))]
