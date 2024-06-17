package v3

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/what3words/w3w-go-wrapper/internal/client"
	"github.com/what3words/w3w-go-wrapper/internal/version"
	"github.com/what3words/w3w-go-wrapper/pkg/core"
)

// API models the what3words public V3 api, with each
// endpoint having its own corresponding methods that would
// return strictly typed structures or errors.
// APIOptions can be used to add/modify aspects of
// the API Controller.
//
// By default:
//
//	baseURL: https://api.what3words.com
//	headers: None
//	client: http.DefaultClient
//
// TOOD: Change to API
type API struct {
	baseURL string
	headers map[string]string
	client  client.HttpClient
}

type APIOption func(*API)

// WithCustomHeader sets a new custom header which
// would be sent with every request made through
// this API
func WithCustomHeader(key, value string) APIOption {
	return func(vs *API) {
		vs.headers[key] = value
	}
}

// WithClient sets a new custom http client can be used
// to set the underlying http client, If a non default
// http client is required. For example use a custom
// transport layer or to mock the API.
func WithClient(client client.HttpClient) func(*API) {
	return func(vs *API) {
		vs.client = client
	}
}

// WithCustomBaseURL sets a new custom base url which
// overdide the defaut base api url. V3 would be suffixed
// to the provided input making the whole address follow
// this pattern <base_url>/v3/<endpoint>
func WithCustomBaseURL(baseURL string) func(*API) {
	return func(vs *API) {
		vs.baseURL = fmt.Sprintf("%s/v3", baseURL)
	}
}

// NewAPI creates a new what3words API Controller.
func NewAPI(apiKey string, opts ...APIOption) *API {
	headers := make(map[string]string)
	headers[core.HEADER_API_KEY] = apiKey
	headers[core.HEADER_WRAPPER] = version.ResolveWrapperHeader()
	baseURL := core.BASE_URL
	svc := &API{
		fmt.Sprintf("%s/v3", baseURL),
		headers,
		http.DefaultClient,
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// ConvertTo3wa wraps around /v3/convert-to-3wa which will convert a latitude
// and longitude pair to a 3 word address, in the language of your choice. It also returns country,
// the bounds of the grid square, a nearby place (such as a local town) and a link to our map site.
func (svc API) ConvertTo3wa(ctx context.Context, coordinates core.Coordinates, opts *ConvertAPIOpts) (*ConvertAPIResponse, error) {
	var c2cResponse convertAPIResponse
	queryParams := make(map[string]string)
	queryParams["coordinates"] = coordinates.String()
	if opts != nil {
		maps.Copy(queryParams, opts.asOptionsMap())
	}
	err := core.MakeGetRequest(
		ctx,
		svc.client,
		svc.baseURL,
		queryParams,
		svc.headers,
		&c2cResponse,
		"convert-to-3wa",
	)
	if err != nil {
		return nil, err
	}
	return &ConvertAPIResponse{
		Json:    c2cResponse.ConvertAPIJsonResponse,
		GeoJson: c2cResponse.ConvertAPIGeoJsonResponse,
	}, nil
}

// ConvertToCoordinates wraps around /v3/convert-to-coordinates which will
// convert a 3 word address to a latitude and longitude pair. It also returns
// country, the bounds of the grid square, a nearest place (such as a local town)
// and a link to our map site.
func (svc API) ConvertToCoordinates(ctx context.Context, words string, opts *ConvertAPIOpts) (*ConvertAPIResponse, error) {
	var c2cResponse convertAPIResponse
	queryParams := make(map[string]string)
	queryParams["words"] = words
	if opts != nil {
		maps.Copy(queryParams, opts.asOptionsMap())
	}
	err := core.MakeGetRequest(
		ctx,
		svc.client,
		svc.baseURL,
		queryParams,
		svc.headers,
		&c2cResponse,
		"convert-to-coordinates",
	)
	if err != nil {
		return nil, err
	}
	return &ConvertAPIResponse{
		Json:    c2cResponse.ConvertAPIJsonResponse,
		GeoJson: c2cResponse.ConvertAPIGeoJsonResponse,
	}, nil
}

// AutoSuggest wraps around /v3/autosuggest endpoint which takes slightly
// incorrect 3 word address and suggest a list of valid 3 word addresses.
// It has powerful features that can, for example, optionally limit results
// to a country or area, and prioritise results that are near the user (see Clipping and Focus below).
//
// It provides corrections for the following types of input error:
//
// Typing errors
// - Spelling errors
// - Misremembered words (e.g. singular vs. plural)
//
// Input 3 word address
// AutoSuggest accepts either a full or partial 3 word address (it will be partial
// if the user is still typing in a search box, for example). A partial 3 word address
// must contain at least the first two words and first character of the third word.
// For example filled.count.s will return results, but anything shorter will not.
//
// Clipping and Focus
// Our clipping allows you to specify a country (or list of countries) and/or
// geographic area to exclude results that are not likely to be relevant to your
// users. To give a more targeted, shorter set of results to your users,
// we recommend you use the clipping parameters. If you know your user’s current
// location, we also strongly recommend that you use the focus to return results
// that are likely to be more relevant (i.e. results near the user).
//
// In summary, the clipping policy is used to optionally restrict the list of candidate
// AutoSuggest results, after which, if focus has been supplied, this will be used to weight
// the results in order of relevance to the focus.
//
// Multiple clipping policies can be specified, though only one of each type. For example you
// can clip to country and clip to circle in the same AutoSuggest call, and it will clip to the
// intersection of the two (results must be in the circle AND in the country).
//
// Language
// AutoSuggest will search in all languages. However, you can additionally specify a fallback language,
// to help the API in situations where the input is particularly messy. For normal text input,
// the language parameter is optional, and AutoSuggest will work well even without a language parameter.
// However, for voice input the language should always be specified.
func (svc API) AutoSuggest(ctx context.Context, input string, opts *AutoSuggestOpts) (*AutoSuggestResponse, error) {
	var autoSuggest autoSuggestResponse
	queryParams := make(map[string]string)
	queryParams["input"] = input
	if opts != nil {
		mOpts := opts.asOptionsMap()
		maps.Copy(queryParams, mOpts)
	}
	err := core.MakeGetRequest(
		ctx,
		svc.client,
		svc.baseURL,
		queryParams,
		svc.headers,
		&autoSuggest,
		"autosuggest",
	)
	if err != nil {
		return nil, err
	}
	return &autoSuggest.AutoSuggestResponse, nil
}

// GridSection wraps around /v3/grid-section which will return a section
// section of the 3m x 3m what3words grid for a bounding box. The bounding box
// is specified by lat,lng,lat,lng as south,west,north,east. You can request the
// grid in GeoJSON format, making it very simple to display on a map.
func (svc API) GridSection(ctx context.Context, boundingBox BoundingBox, opts *GridSectionOpts) (*GridSectionResponse, error) {
	var gridSection gridSectionResponse
	queryParams := make(map[string]string)
	queryParams["bounding-box"] = boundingBox.String()
	if opts != nil {
		maps.Copy(queryParams, opts.asOptionsMap())
	}
	err := core.MakeGetRequest(
		ctx,
		svc.client,
		svc.baseURL,
		queryParams,
		svc.headers,
		&gridSection,
		"grid-section",
	)
	if err != nil {
		return nil, err
	}
	return &GridSectionResponse{
		Json:    gridSection.GridSectionJsonResponse,
		GeoJson: gridSection.GridSectionGeoJsonResponse,
	}, nil
}

// AvailableLanguages wraps around /v3/available-languages which will
// retrieve a list of all available 3 word address languages,
// including the ISO 3166-1 alpha-2 2 letter code, English name and native name.
// Bosnian-Croatian-Montenegrin-Serbian is available using the language code 'oo' with
// Cyrillic and Latin locales ('oo_cy' and 'oo_la')
func (svc API) AvailableLanguages(ctx context.Context) (*AvailableLanguagesResponse, error) {
	var availableLanguages availableLanguagesResponse
	err := core.MakeGetRequest(ctx, svc.client, svc.baseURL, map[string]string{}, svc.headers, &availableLanguages, "available-languages")
	if err != nil {
		return nil, err
	}
	return &availableLanguages.AvailableLanguagesResponse, nil
}