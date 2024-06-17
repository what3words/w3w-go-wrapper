package client

import (
	"net/http"
)

// HTTPClient defined so that net/http.Client satisfies it
// while allowing usage of mock or fake clients for easy
// testing
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
