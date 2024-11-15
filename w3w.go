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

// Service wraps the what3words public Service with each
// version of the API available under its own
// method. A call for example to the v3/available-languages
// api would be made as such
//
// svc := NewService(apiKey)
//
// languages, err := svc.V3().AvailableLanguages(context.Background())
//
//	if err != nil {
//		return err
//	}
//go:generate mockery --name Service --output ./mocks --outpkg mocks --case underscore
type Service interface {
	V3() v3.API
	FindPossible3wa(input string) []string
	IsPossible3wa(input string) bool
	IsValid3wa(ctx context.Context, input string) bool
	DidYouMean(input string) bool
}

type service struct {
	v3api v3.API
}

type ServiceOpts func(*service)

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
	return func(svc *service) {
		svc.v3api.SetBaseURL(baseURL)
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
	return func(svc *service) {
		svc.v3api.SetHeader(key, value)
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
	return func(svc *service) {
		svc.v3api.SetClient(client)
	}
}

// WithV3API allows you to set a custom what3words v3
// service. Construct a v3 service with the w3w-go-wrapper/pkg/v3 NewService
// function and set it.
// Usefull if you want to set Custom anything specific to a version
func WithV3API(v3 v3.API) ServiceOpts {
	return func(svc *service) {
		svc.v3api = v3
	}
}

// NewService creates a new what3words Service wrapper
func NewService(ServiceKey string, opts ...ServiceOpts) Service {
	svc := service{
		v3api: v3.NewAPI(ServiceKey),
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
