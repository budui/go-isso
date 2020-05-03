// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// @budui copy from "miniflux.app/http/response/json", may be modified later.

package json

import (
	"encoding/json"
	"errors"
	"net/http"

	"wrong.wang/x/go-isso/isso/response"
	"wrong.wang/x/go-isso/logger"
)

const contentTypeHeader = `application/json`

// OK creates a new JSON response with a 200 status code.
func OK(builder response.Builder, body interface{}) {
	builder.WithHeader("Content-Type", contentTypeHeader)
	builder.WithBody(toJSON(body))
	builder.Write()
}

// Created sends a created response to the client.
func Created(builder response.Builder, body interface{}) {
	builder.WithStatus(http.StatusCreated)
	builder.WithHeader("Content-Type", contentTypeHeader)
	builder.WithBody(toJSON(body))
	builder.Write()
}

// ServerError sends an internal error to the client.
func ServerError(builder response.Builder, err error) {
	builder.WithError(err)
	builder.WithStatus(http.StatusInternalServerError)
	builder.WithHeader("Content-Type", contentTypeHeader)
	builder.WithBody(toJSONError(err))
	builder.Write()
}

// BadRequest sends a bad request error to the client.
func BadRequest(builder response.Builder, err error) {
	builder.WithError(err)
	builder.WithStatus(http.StatusBadRequest)
	builder.WithHeader("Content-Type", contentTypeHeader)
	builder.WithBody(toJSONError(err))
	builder.Write()
}

// Unauthorized sends a not authorized error to the client.
func Unauthorized(builder response.Builder) {
	builder.WithStatus(http.StatusUnauthorized)
	builder.WithHeader("Content-Type", contentTypeHeader)
	builder.WithBody(toJSONError(errors.New("Access Unauthorized")))
	builder.Write()
}

// Forbidden sends a forbidden error to the client.
func Forbidden(builder response.Builder) {
	builder.WithStatus(http.StatusForbidden)
	builder.WithHeader("Content-Type", contentTypeHeader)
	builder.WithBody(toJSONError(errors.New("Access Forbidden")))
	builder.Write()
}

// NotFound sends a page not found error to the client.
func NotFound(builder response.Builder) {
	builder.WithStatus(http.StatusNotFound)
	builder.WithHeader("Content-Type", contentTypeHeader)
	builder.WithBody(toJSONError(errors.New("Resource Not Found")))
	builder.Write()
}

func toJSONError(err error) []byte {
	type errorMsg struct {
		ErrorMessage string `json:"error_message"`
	}

	return toJSON(errorMsg{ErrorMessage: err.Error()})
}

func toJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		logger.Error("[HTTP:JSON] %v", err)
		return []byte("")
	}

	return b
}
