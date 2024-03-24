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

var FailStatusCode = statusCodes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(statusCodes))]
