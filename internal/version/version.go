package version

import (
	"fmt"
)

//go:generate go run gen_version.go
//go:generate go fmt wrapper_version.go

// WRAPPER_PREFIX prefixes the value of the x-w3w-wrapper header
// in the following format <prefix>/<verison>. Used to collect
// metrics from requests recieved.
const WRAPPER_PREFIX = "what3words-golang"

// ResolveWrapperHeader to resolve x-w3w-wrapper header value in
// valid What3Words format. Seperated into a function to be able
// easily ingect version code from CI/CD
func ResolveWrapperHeader() string {
	version := wrapper_version()
	return fmt.Sprintf("%s/%s", WRAPPER_PREFIX, version)
}
