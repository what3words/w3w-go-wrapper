package core

// ResponseReader contains functions that all API responses
// need to implement. Provides easy abstraction to check if
// an error occurred within the response.
type ResponseReader interface {
	GetError() error
}
