// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// @budui copy from "miniflux.app/http/response/json", may be modified later.

package json

import (
	"encoding/json"
	"net/http"

	"wrong.wang/x/go-isso/isso/response"
	"wrong.wang/x/go-isso/logger"
)

const contentTypeHeader = `application/json`

func writeJSON(builder response.Builder, body interface{}, status int) {
	builder.WithHeader("Content-Type", contentTypeHeader)
	builder.WithStatus(status)
	builder.WithBody(toJSON(body))
	builder.Write()
}

func writeErrorJSON(builder response.Builder, err error, status int) {
	if err != nil {
		writeJSON(builder, map[string]string{"error": err.Error()}, status)
	}
	writeJSON(builder, map[string]string{"error": http.StatusText(status)}, status)
}

// OK creates a new JSON response with a 200 status code.
func OK(builder response.Builder, body interface{}) {
	writeJSON(builder, body, http.StatusOK)
}

// Created sends a created response to the client.
func Created(builder response.Builder, body interface{}) {
	writeJSON(builder, body, http.StatusCreated)
}

// Accepted sends a created response to the client.
func Accepted(builder response.Builder, body interface{}) {
	writeJSON(builder, body, http.StatusAccepted)
}

// ServerError sends an internal error to the client.
func ServerError(builder response.Builder, err error) {
	writeErrorJSON(builder, err, http.StatusInternalServerError)
}

// BadRequest sends a bad request error to the client.
func BadRequest(builder response.Builder, err error) {
	writeErrorJSON(builder, err, http.StatusBadRequest)
}

// Unauthorized sends a not authorized error to the client.
func Unauthorized(builder response.Builder) {
	writeErrorJSON(builder, nil, http.StatusUnauthorized)
}

// Forbidden sends a forbidden error to the client.
func Forbidden(builder response.Builder) {
	writeErrorJSON(builder, nil, http.StatusForbidden)
}

// NotFound sends a page not found error to the client.
func NotFound(builder response.Builder) {
	writeErrorJSON(builder, nil, http.StatusNotFound)
}

func toJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		logger.Error("[HTTP:JSON] %v", err)
		return []byte("")
	}

	return b
}
