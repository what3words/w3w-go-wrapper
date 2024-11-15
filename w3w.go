package w3wgowrapper

import (
	"context"
	"regexp"

	"github.com/what3words/w3w-go-wrapper/internal/client"
	v3 "github.com/what3words/w3w-go-wrapper/pkg/apis/v3"
	"github.com/what3words/w3w-go-wrapper/pkg/core"
)

var (
	regexFind3wa       = regexp.MustCompile(`[^0-9\x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]{1,}[.｡。･・︒។։။۔።।][^0-9\x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]{1,}[.｡。･・︒។։။۔።।][^0-9x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]{1,}`)
	regexIsPossible3wa = regexp.MustCompile(`^/*(?:[^0-9\x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]{1,}[.｡。･・︒។։။۔።।][^0-9\x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]{1,}[.｡。･・︒។։။۔።।][^0-9\x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]{1,}|[^0-9\x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]{1,}([\x{0020}\x{00A0}][^0-9\x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]+){1,3}[.｡。･・︒។։။۔።।][^0-9\x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]{1,}([\x{0020}\x{00A0}][^0-9\x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]+){1,3}[.｡。･・︒។։။۔።।][^0-9\x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]{1,}([\x{0020}\x{00A0}][^0-9\x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]+){1,3})$`)
	regexDidYouMean    = regexp.MustCompile(`^/?[^0-9\x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]{1,}[.\x{FF61}\x{3002}\x{FF65}\x{30FB}\x{FE12}\x{17D4}\x{0964}\x{1362}\x{3002}:။^_۔։ ,\\/+'&\\:;|\x{3000}-]{1,2}[^0-9\x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]{1,}[.\x{FF61}\x{3002}\x{FF65}\x{30FB}\x{FE12}\x{17D4}\x{0964}\x{1362}\x{3002}:။^_۔։ ,\\/+'&\\:;|\x{3000}-]{1,2}[^0-9\x60~!@#$%^&*()+\-_=\[\{\]}\\|'<>.,?/;:£§º©®\s]{1,}$`)
)

// Service wraps the What3Words public API, providing methods for each
// version of the API available under its own method. For example, a call
// to the v3/available-languages API would be made as follows:
//
//	svc := NewService(apiKey)
//
//	languages, err := svc.V3().AvailableLanguages(context.Background())
//
//	if err != nil {
//		return err
//	}
//
//go:generate mockery --name Service --output ./mocks --outpkg mocks --case underscore
type Service interface {

	// V3 provides access to methods/endpoints under the v3 version of the What3Words API.
	// Example usage:
	//   svc := NewService(apiKey)
	//   result, err := svc.V3().AvailableLanguages(ctx)
	//   if err != nil {
	//       // handle error
	//   }
	V3() v3.API
	// FindPossible3wa searches the string passed in for all substrings in the form of a three word address.
	FindPossible3wa(input string) []string
	// IsPossible3wa determines if the string passed in is in the form of a three word address.
	IsPossible3wa(input string) bool
	// DidYouMean determines if the string passed in is almost in the form of a three word address.
	DidYouMean(input string) bool
	// IsValid3wa validates the given string as a real three-word address by
	// making a call to the API. The context can be used to cancel the underlying call.
	IsValid3wa(ctx context.Context, input string) bool
}

type service struct {
	v3api v3.API
}

type ServiceOpts func(*service)

// WithCustomBaseURL allows you to set a custom base URL for the What3Words service.
// This is useful for testing purposes or when interacting with a self-hosted
// enterprise server.
//
// # Note:
// This custom base URL is applied to all API versions within the Service.
//
// Example usage:
//
//	service := NewService(apiKey, WithCustomBaseURL("https://custom-base-url.example.com"))
//
//	languages, err := service.V3().AvailableLanguages(context.Background())
//	if err != nil {
//	    return err
//	}
func WithCustomBaseURL(baseURL string) ServiceOpts {
	return func(svc *service) {
		svc.v3api.SetBaseURL(baseURL)
	}
}

// WithCustomHeader allows you to set a custom header for the What3Words service.
// All requests sent from this client will include this header.
//
// # Note:
// The custom header is applied to all API versions within the Service.
//
// Example usage:
//
//	service := NewService(apiKey, WithCustomHeader("X-Custom-Header", "header-value"))
//
//	languages, err := service.V3().AvailableLanguages(context.Background())
//	if err != nil {
//	    return err
//	}
func WithCustomHeader(key, value string) ServiceOpts {
	return func(svc *service) {
		svc.v3api.SetHeader(key, value)
	}
}

// WithClient allows you to set a custom HTTP client for the What3Words service.
// This is useful for testing purposes or when you need to configure a client
// with non-default settings for the net/http package.
//
// The client must implement the w3w-go-wrapper/internal/client/HttpClient interface.
//
// # Note:
// The custom HTTP client is applied to all API versions within the Service.
//
// Example usage:
//
//	service := NewService(apiKey, WithClient(&http.Client{}))
func WithClient(client client.HttpClient) ServiceOpts {
	return func(svc *service) {
		svc.v3api.SetClient(client)
	}
}

// WithV3API allows you to set a custom What3Words v3 service.
// You can construct a v3 service using the w3w-go-wrapper/pkg/v3 `NewService` function
// and configure it as needed before setting it.
//
// This is useful if you need to customize settings specific to the v3 API.
func WithV3API(v3 v3.API) ServiceOpts {
	return func(svc *service) {
		svc.v3api = v3
	}
}

// NewService creates a new What3Words service wrapper.
// This function initializes the service with the provided API key and applies
// any optional configurations specified through ServiceOpts.
//
// # Parameters:
//   - `apiKey`: The API key used for authenticating requests to the What3Words API.
//   - `opts`: Optional configurations for customizing the service, such as setting
//     a custom base URL, HTTP client, or v3 API service.
//
// Example usage:
//
//	service := NewService("your-api-key",
//		WithCustomBaseURL("https://custom-url.example.com"),
//		WithClient(&http.Client{}),
//	)
//
//	languages, err := service.V3().AvailableLanguages(context.Background())
//	if err != nil {
//	    // handle error
//	}
//
// # Returns:
// A Service interface that provides access to What3Words API methods.
func NewService(apiKey string, opts ...ServiceOpts) Service {
	svc := service{
		v3api: v3.NewAPI(apiKey),
	}
	for _, opt := range opts {
		opt(&svc)
	}
	return svc
}

func (svc service) V3() v3.API {
	return svc.v3api
}

func (svc service) FindPossible3wa(input string) []string {
	return regexFind3wa.FindAllString(input, -1)
}

func (svc service) IsPossible3wa(input string) bool {
	return regexIsPossible3wa.MatchString(input)
}

func (svc service) IsValid3wa(ctx context.Context, input string) bool {
	if svc.IsPossible3wa(input) {
		if resp, err := svc.V3().AutoSuggest(ctx, input, &v3.AutoSuggestOpts{
			NResults: core.Int(1),
		}); err != nil {
			if len(resp.Suggestions) >= 1 {
				return resp.Suggestions[0].Words == input
			}
		}
	}
	return false
}

func (svc service) DidYouMean(input string) bool {
	return regexDidYouMean.MatchString(input)
}
