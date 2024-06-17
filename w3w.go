package w3wgowrapper

import (
	"github.com/what3words/w3w-go-wrapper/internal/client"
	v3 "github.com/what3words/w3w-go-wrapper/pkg/apis/v3"
)

// Service wraps the what3words public Service with each
// version of the API available under its own
// attibute. A call for example to the v3/available-languages
// api would be made as such
//
// svc := NewService(ServiceKey)
//
// languages, err := svc.V3.AvailableLanguages(context.Background())
//
//	if err != nil {
//		return err
//	}

type Service struct {
	V3 *v3.API
}

type ServiceOpts func(*Service)

// WithCustomBaseURL allows you to set a custom base url
// for the what3words service. This is useful for testing or when calling an
// self-hosted enterprise server
//
// # This sets the base url for all versions of the API in the Service
//
// For example:
//
// Service := NewService(ServiceKey, WithCustomBaseURL("XXXXXXXXXXXXXXXXXXXXXXXXXX"))
//
// languages, err := Service.V3.AvailableLanguages(context.Background())
//
//	if err != nil {
//		return err
//	}
func WithCustomBaseURL(baseURL string) ServiceOpts {
	return func(Service *Service) {
		v3.WithCustomBaseURL(baseURL)(Service.V3)
	}
}

// WithCustomHeader allows you to set a custom header
// for the what3words Service. All requests sent from
// this client would contain this header
//
// # This sets the header for all versions in the Service
//
// For example:
//
// Service := NewService(ServiceKey, WithCustomHeader("X-Custom-Header", "XXXXXXXXXXXXXXXXXXXXXXXXXX"))
//
// languages, err := Service.V3.AvailableLanguages(context.Background())
//
//	if err != nil {
//		return err
//	}
func WithCustomHeader(key, value string) ServiceOpts {
	return func(Service *Service) {
		v3.WithCustomHeader(key, value)(Service.V3)
	}
}

// WithClient allows you to set a custom http client
// for the what3words Service. This is useful for testing or
// setting a client with non-default net/http settings.
//
// Client is anything that implements w3w-go-wrapper/internal/client/HttpClient
// interface
//
// # This sets the http client for all versions in the Service
//
// For example:
//
// Service := NewService(ServiceKey, WithClient(&http.Client{})
func WithClient(client client.HttpClient) ServiceOpts {
	return func(Service *Service) {
		v3.WithClient(client)(Service.V3)
	}
}

// WithV3API allows you to set a custom what3words v3
// service. Construct a v3 service with the w3w-go-wrapper/pkg/v3 NewService
// function and set it.
// Usefull if you want to set Custom anything specific to a version
func WithV3API(v3 *v3.API) ServiceOpts {
	return func(Service *Service) {
		Service.V3 = v3
	}
}

// NewService creates a new what3words Service wrapper
func NewService(ServiceKey string, opts ...ServiceOpts) *Service {
	Service := &Service{
		V3: v3.NewAPI(ServiceKey),
	}
	for _, opt := range opts {
		opt(Service)
	}
	return Service
}
