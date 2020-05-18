package json

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"wrong.wang/x/go-isso/logger"
	"wrong.wang/x/go-isso/version"
)

const contentTypeHeader = `application/json`

func writeJSON(w http.ResponseWriter, body interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	_ = encoder.Encode(body)
}

func writeErrorJSON(w http.ResponseWriter, err error, requestID string, desc string, status int) {
	var caller string
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		caller = "unkown"
	} else {
		fn := runtime.FuncForPC(pc)
		caller = strings.TrimPrefix(fn.Name(), version.Mod)
	}

	errstr := fmt.Sprintf("%s %s - %s", requestID, caller, desc)
	if err != nil {
		errstr = fmt.Sprintf("%s\n\t%v", errstr, err)
	}
	logger.Error("%s", errstr)

	reason := desc
	if reason == "" {
		reason = http.StatusText(status)
	}
	writeJSON(w, map[string]string{"error": reason}, status)
}

// OK creates a new JSON response with a 200 status code.
func OK(w http.ResponseWriter, body interface{}) {
	writeJSON(w, body, http.StatusOK)
}

// Created sends a created response to the client.
func Created(w http.ResponseWriter, body interface{}) {
	writeJSON(w, body, http.StatusCreated)
}

// Accepted sends a created response to the client.
func Accepted(w http.ResponseWriter, body interface{}) {
	writeJSON(w, body, http.StatusAccepted)
}

// ServerError sends an internal error to the client.
func ServerError(requestID string, w http.ResponseWriter, err error, desc string) {
	writeErrorJSON(w, err, requestID, desc, http.StatusInternalServerError)
}

// BadRequest sends a bad request error to the client.
func BadRequest(requestID string, w http.ResponseWriter, err error, desc string) {
	writeErrorJSON(w, err, requestID, desc, http.StatusBadRequest)
}

// Unauthorized sends a not authorized error to the client.
func Unauthorized(requestID string, w http.ResponseWriter, err error, desc string) {
	writeErrorJSON(w, err, requestID, desc, http.StatusUnauthorized)
}

// Forbidden sends a forbidden error to the client.
func Forbidden(requestID string, w http.ResponseWriter, err error, desc string) {
	writeErrorJSON(w, err, requestID, desc, http.StatusForbidden)
}

// NotFound sends a page not found error to the client.
func NotFound(requestID string, w http.ResponseWriter, err error, desc string) {
	writeErrorJSON(w, err, requestID, desc, http.StatusNotFound)
}
