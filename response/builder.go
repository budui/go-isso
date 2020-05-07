package response

// Builder generates HTTP or serverless responses.
type Builder interface {
	WithStatus(statusCode int)
	WithHeader(key, value string)
	WithBody(body interface{})
	WithError(error)
	Write()
}
