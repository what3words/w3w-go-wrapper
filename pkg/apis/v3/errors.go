package v3

import "fmt"

type ErrorCode string

const (
	ErrorCodeMissingWords ErrorCode = "MissingWords"
	ErrorCodeBadWords     ErrorCode = "BadWords"
	ErrorCodeInvalidKey   ErrorCode = "InvalidKey"
	ErrorCodeBadLanguage  ErrorCode = "BadLanguage"
	// TODO: Add all error codes
)

// ErrorResponse models format of the error response
// recived from the API when a non 200 status code
// is recieved.
// Implements std `error` interface, so that it
// can be returned as an error.
//
// Example inferance:
//
//	errResp, ok := err.(*ErrorResponse)
type ErrorResponse struct {
	// Code can be ustized to programatically determine
	// the error.
	Code ErrorCode `json:"code"`
	// Message is intended to be helpful human readable
	// version of the error code.
	Message string `json:"message"`
}

func (er ErrorResponse) Error() string {
	return fmt.Sprintf("api: got error response '%s' with message '%s'", er.Code, er.Message)
}
