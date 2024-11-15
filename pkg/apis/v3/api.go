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

// API models the What3Words public v3 API. Each endpoint has a corresponding method
// that returns strictly typed structures or errors. APIOptions can be used to
// configure or modify the API Controller when creating an instance with the `NewAPI` function.
//
// # Default Configuration:
// - `baseURL`: https://api.what3words.com
// - `headers`: None
// - `client`: http.DefaultClient
//
//go:generate mockery --name API --output ./mocks --outpkg mocks --case underscore
type API interface {
	// Configuration Setters

	// SetBaseURL sets the base URL for the What3Words API after initialization.
	SetBaseURL(baseURL string)
	// SetHeader sets a single HTTP header to include in all API requests after initialization.
	SetHeader(headerKey, headerValue string)
	// SetHeaderMap sets multiple HTTP headers at once for all API requests after initialization.
	SetHeaderMap(headers map[string]string)
	// SetClient sets a custom HTTP client for API requests after initialization.
	SetClient(client client.HttpClient)

	// Endpoints

	// ConvertTo3wa converts a latitude and longitude pair to a three-word address.
	// It also provides additional information such as country, grid square bounds,
	// a nearby place, and a link to the map site.
	// Returns the response in JSON format.
	ConvertTo3wa(ctx context.Context, coordinates core.Coordinates, opts *ConvertAPIOpts) (*ConvertAPIJsonResponse, error)
	// ConvertTo3waGeoJson performs the same conversion as `ConvertTo3wa` but
	// returns the response in GeoJSON format.
	ConvertTo3waGeoJson(ctx context.Context, coordinates core.Coordinates, opts *ConvertAPIOpts) (*ConvertAPIGeoJsonResponse, error)
	// ConvertToCoordinates wraps around /v3/convert-to-coordinates which will
	// convert a 3 word address to a latitude and longitude pair. It also returns
	// country, the bounds of the grid square, a nearest place (such as a local town)
	// and a link to our map site. Returns the response in the JSON format.
	ConvertToCoordinates(ctx context.Context, words string, opts *ConvertAPIOpts) (*ConvertAPIJsonResponse, error)
	// ConvertTo3waGeoJson performs the same conversion as `ConvertToCoordinates` but
	// returns the response in GeoJSON format.
	ConvertToCoordinatesGeoJson(ctx context.Context, words string, opts *ConvertAPIOpts) (*ConvertAPIGeoJsonResponse, error)
	// GridSection wraps around the /v3/grid-section endpoint, returning a section
	// of the 3m x 3m What3Words grid for a specified bounding box.
	//
	// The bounding box is defined by two coordinates:
	// - `SouthWest` (latitude and longitude of the bottom-left corner)
	// - `NorthEast` (latitude and longitude of the top-right corner).

	// Response is returned in the JSON format.
	GridSection(ctx context.Context, boundingBox BoundingBox) (*GridSectionJsonResponse, error)
	// GridSectionGeoJson wraps around the /v3/grid-section endpoint, returning a section
	// of the 3m x 3m What3Words grid for a specified bounding box in GeoJSON format.
	//
	// The bounding box is defined by two coordinates:
	// - `SouthWest` (latitude and longitude of the bottom-left corner)
	// - `NorthEast` (latitude and longitude of the top-right corner).
	//
	// GeoJSON format is particularly useful for rendering on maps or integrating
	// with GIS tools, as it provides structured geospatial data.
	GridSectionGeoJson(ctx context.Context, boundingBox BoundingBox) (*GridSectionGeoJsonResponse, error)
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
	AutoSuggest(ctx context.Context, input string, opts *AutoSuggestOpts) (*AutoSuggestResponse, error)
	// AvailableLanguages wraps around /v3/available-languages which will
	// retrieve a list of all available 3 word address languages,
	// including the ISO 3166-1 alpha-2 2 letter code, English name and native name.
	// Bosnian-Croatian-Montenegrin-Serbian is available using the language code 'oo' with
	// Cyrillic and Latin locales ('oo_cy' and 'oo_la')
	AvailableLanguages(ctx context.Context) (*AvailableLanguagesResponse, error)
}

type api struct {
	baseURL string
	headers map[string]string
	client  client.HttpClient
}

func (a *api) SetBaseURL(baseURL string) {
	a.baseURL = fmt.Sprintf("%s/v3", baseURL)
}

func (a *api) SetHeader(headerKey, headerValue string) {
	a.headers[headerKey] = headerValue
}

func (a *api) SetHeaderMap(headers map[string]string) {
	a.headers = headers
}

func (a *api) SetClient(client client.HttpClient) {
	a.client = client
}

type APIOption func(*api)

// WithCustomHeader sets a custom HTTP header to be included with every request
// made through this API. This is useful for scenarios like adding authentication
// tokens or other custom headers required by the API.
//
// Example usage:
//
//	api := NewAPI("your-api-key", WithCustomHeader("X-Custom-Header", "value"))
func WithCustomHeader(key, value string) APIOption {
	return func(vs *api) {
		vs.headers[key] = value
	}
}

// WithClient sets a custom HTTP client to be used by the API. This is useful when
// you need to configure a custom transport layer, use a specific client for testing,
// or mock the API client.
//
// Example usage:
//
//	customClient := &http.Client{Timeout: 10 * time.Second}
//	api := NewAPI("your-api-key", WithClient(customClient))
func WithClient(client client.HttpClient) func(*api) {
	return func(vs *api) {
		vs.client = client
	}
}

// WithCustomBaseURL sets a custom base URL for the API, overriding the default base
// URL. The provided URL will be suffixed with `/v3`, forming the full base URL
// for API requests, following the pattern: <base_url>/v3/<endpoint>.
//
// Example usage:
//
//	api := NewAPI("your-api-key", WithCustomBaseURL("https://custom-url.example.com"))
func WithCustomBaseURL(baseURL string) func(*api) {
	return func(vs *api) {
		vs.baseURL = fmt.Sprintf("%s/v3", baseURL)
	}
}

// NewAPI creates a new What3Words V3 API Controller instance.
//
// This function initializes an API controller with the provided API key and
// applies any optional configurations specified through APIOption functions.
//
// The default configuration includes:
// - `baseURL`: https://api.what3words.com/v3
// - `headers`: Contains the API key and wrapper version header.
// - `client`: Uses the default HTTP client (http.DefaultClient).
//
// You can customize the API controller by passing options such as custom headers,
// HTTP clients, or base URLs.
//
// Example usage:
//
//	api := NewAPI("your-api-key",
//	    WithCustomHeader("X-Custom-Header", "value"),
//	    WithClient(customClient),
//	)
//
// # Parameters:
// - `apiKey`: The API key for authenticating requests to the What3Words API.
// - `opts`: Optional configuration functions to modify the API controller.
//
// # Returns:
// A new API controller instance with the specified configuration.
func NewAPI(apiKey string, opts ...APIOption) API {
	headers := make(map[string]string)
	headers[core.HEADER_API_KEY] = apiKey
	headers[core.HEADER_WRAPPER] = version.ResolveWrapperHeader()
	baseURL := core.BASE_URL
	a := &api{
		fmt.Sprintf("%s/v3", baseURL),
		headers,
		http.DefaultClient,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a api) convertTo3wa(ctx context.Context, coordinates core.Coordinates, opts *ConvertAPIOpts, format string) (*convertAPIResponse, error) {
	var c2cResponse convertAPIResponse
	queryParams := make(map[string]string)
	queryParams["coordinates"] = coordinates.AsQueryParam()
	queryParams["format"] = format
	if opts != nil {
		maps.Copy(queryParams, opts.asOptionsMap())
	}
	err := core.MakeGetRequest(
		ctx,
		a.client,
		a.baseURL,
		queryParams,
		a.headers,
		&c2cResponse,
		"convert-to-3wa",
	)
	if err != nil {
		return nil, err
	}
	return &c2cResponse, nil
}

func (a api) ConvertTo3wa(ctx context.Context, coordinates core.Coordinates, opts *ConvertAPIOpts) (*ConvertAPIJsonResponse, error) {
	resp, err := a.convertTo3wa(ctx, coordinates, opts, "json")
	if err != nil {
		return nil, err
	}
	return resp.ConvertAPIJsonResponse, nil
}

func (a api) ConvertTo3waGeoJson(ctx context.Context, coordinates core.Coordinates, opts *ConvertAPIOpts) (*ConvertAPIGeoJsonResponse, error) {
	resp, err := a.convertTo3wa(ctx, coordinates, opts, "geojson")
	if err != nil {
		return nil, err
	}
	return resp.ConvertAPIGeoJsonResponse, nil
}

func (a api) convertToCoordinates(ctx context.Context, words string, opts *ConvertAPIOpts, format string) (*convertAPIResponse, error) {
	var c2cResponse convertAPIResponse
	queryParams := make(map[string]string)
	queryParams["words"] = words
	queryParams["format"] = format
	if opts != nil {
		maps.Copy(queryParams, opts.asOptionsMap())
	}
	err := core.MakeGetRequest(
		ctx,
		a.client,
		a.baseURL,
		queryParams,
		a.headers,
		&c2cResponse,
		"convert-to-coordinates",
	)
	if err != nil {
		return nil, err
	}
	return &c2cResponse, nil
}

func (a api) ConvertToCoordinates(ctx context.Context, words string, opts *ConvertAPIOpts) (*ConvertAPIJsonResponse, error) {
	resp, err := a.convertToCoordinates(ctx, words, opts, "json")
	if err != nil {
		return nil, err
	}
	return resp.ConvertAPIJsonResponse, nil
}

func (a api) ConvertToCoordinatesGeoJson(ctx context.Context, words string, opts *ConvertAPIOpts) (*ConvertAPIGeoJsonResponse, error) {
	resp, err := a.convertToCoordinates(ctx, words, opts, "geojson")
	if err != nil {
		return nil, err
	}
	return resp.ConvertAPIGeoJsonResponse, nil
}

func (a api) AutoSuggest(ctx context.Context, input string, opts *AutoSuggestOpts) (*AutoSuggestResponse, error) {
	var autoSuggest autoSuggestResponse
	queryParams := make(map[string]string)
	queryParams["input"] = input
	if opts != nil {
		mOpts := opts.asOptionsMap()
		maps.Copy(queryParams, mOpts)
	}
	err := core.MakeGetRequest(
		ctx,
		a.client,
		a.baseURL,
		queryParams,
		a.headers,
		&autoSuggest,
		"autosuggest",
	)
	if err != nil {
		return nil, err
	}
	return &autoSuggest.AutoSuggestResponse, nil
}

func (a api) gridSection(ctx context.Context, boundingBox BoundingBox, format string) (*gridSectionResponse, error) {
	var gridSection gridSectionResponse
	queryParams := make(map[string]string)
	queryParams["bounding-box"] = boundingBox.asQueryParam()
	queryParams["format"] = format
	err := core.MakeGetRequest(
		ctx,
		a.client,
		a.baseURL,
		queryParams,
		a.headers,
		&gridSection,
		"grid-section",
	)
	if err != nil {
		return nil, err
	}
	return &gridSection, nil
}

func (a api) GridSection(ctx context.Context, boundingBox BoundingBox) (*GridSectionJsonResponse, error) {
	resp, err := a.gridSection(ctx, boundingBox, "json")
	if err != nil {
		return nil, err
	}
	return resp.GridSectionJsonResponse, nil
}

func (a api) GridSectionGeoJson(ctx context.Context, boundingBox BoundingBox) (*GridSectionGeoJsonResponse, error) {
	resp, err := a.gridSection(ctx, boundingBox, "geojson")
	if err != nil {
		return nil, err
	}
	return resp.GridSectionGeoJsonResponse, nil
}

func (a api) AvailableLanguages(ctx context.Context) (*AvailableLanguagesResponse, error) {
	var availableLanguages availableLanguagesResponse
	err := core.MakeGetRequest(ctx, a.client, a.baseURL, map[string]string{}, a.headers, &availableLanguages, "available-languages")
	if err != nil {
		return nil, err
	}
	return &availableLanguages.AvailableLanguagesResponse, nil
}
