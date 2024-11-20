package v3

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/what3words/w3w-go-wrapper/pkg/core"
)

type Coordinates = core.Coordinates

// Convert APIs

// ConvertAPIOpts models all optional options accepted
// by /v3/convert-to-3wa and /v3/convert-to-coordiates endpoints
// of the what3words public api,
type ConvertAPIOpts struct {
	// Locale to specify a variant of a language
	Locale string
	// A supported 3 word address language as an ISO 639-1 2 letter code.
	// For Bosnian-Croatian-Montenegrin-Serbian use oo. Defaults to en (English).
	// For a full list of 3 word address languages, see available-languages.
	Language string
}

func (cto ConvertAPIOpts) asOptionsMap() map[string]string {
	mapOpts := make(map[string]string)
	if cto.Locale != "" {
		mapOpts["locale"] = cto.Locale
	}
	if cto.Language != "" {
		mapOpts["language"] = cto.Language
	}
	return mapOpts
}

// ConvertAPIGeoJsonResponse models the format `geojson`
// returned by the what3words public api convert endpoints
type ConvertAPIGeoJsonResponse struct {
	Features []struct {
		Bbox     []float64 `json:"Bbox"`
		Geometry struct {
			Coordinates  []float64 `json:"coordinates"`
			GeometryType string    `json:"type"`
		}
		FeatureType string `json:"type"`
		Properties  struct {
			Country      string `json:"country"`
			NearestPlace string `json:"nearestPlace"`
			Words        string `json:"words"`
			Language     string `json:"language"`
			Locale       string `json:"locale,omitempty"`
			MapURL       string `json:"map"`
		} `json:"properties"`
	} `json:"features,omitempty"`
	GeoJsonType string `json:"type,omitempty"`
}

// ConvertAPIJsonResponse models the format `json` (Default)
// response returned by the what3words public api convert endpoints.
type ConvertAPIJsonResponse struct {
	Country      string      `json:"country,omitempty"`
	Square       Sqaure      `json:"square,omitempty"`
	NearestPlace string      `json:"nearestPlace,omitempty"`
	Coordinates  Coordinates `json:"coordinates,omitempty"`
	Words        string      `json:"words,omitempty"`
	Language     string      `json:"language,omitempty"`
	MapUrl       string      `json:"map,omitempty"`
}

type convertAPIResponse struct {
	*ConvertAPIJsonResponse
	*ConvertAPIGeoJsonResponse
	Error *ErrorResponse `json:"error,omitempty"`
}

func (ctr convertAPIResponse) GetError() error {
	return ctr.Error
}

// AutoSuggest API

type Circle struct {
	Center   Coordinates
	RadiusKm float64
}

func (c Circle) asQueryParam() string {
	return fmt.Sprintf("%s,%f", c.Center.AsQueryParam(), c.RadiusKm)
}

type Sqaure struct {
	SouthWest Coordinates `json:"southwest"`
	NorthEast Coordinates `json:"northeast"`
}

type BoundingBox Sqaure

func (bb BoundingBox) asQueryParam() string {
	return fmt.Sprintf("%.6f,%.6f,%.6f,%.6f", bb.SouthWest.Lat, bb.SouthWest.Lng, bb.NorthEast.Lat, bb.NorthEast.Lng)
}

type Polygon []Coordinates

// AutoSuggestOpts models all the possible
// optional options availabel for /v3/autosuggest
// api
type AutoSuggestOpts struct {
	// This is a location, specified as latitude,longitude (often where the user making the query is).
	// If specified, the results will be weighted to give preference to those near the `Focus`.
	// For convenience, longitude is allowed to wrap around the 180 line, so 361 is equivalent to 1.
	Focus *Coordinates
	// Restricts AutoSuggest to only return results inside the countries specified in the
	// list of uppercase ISO 3166-1 alpha-2 country codes (for example,
	// to restrict to Belgium and the UK, use []string{"GB","BE"}).
	// ClipToCountry will also accept lowercase country codes.
	// Entries must be two a-z letters.
	// WARNING: If the two-letter code does not correspond to a country,
	// there is no error: API simply returns no results.
	ClipToCountry []string
	// Restrict AutoSuggest results to a bounding box, specified by coordinates.
	// SouthLat, WestLng, NorthLat, EastLng, where:SouthLat less than or equal to NorthLat, WestLng
	// less than or equal to EastLng in other words, latitudes and longitudes should be specified order
	// of increasing size. Lng is allowed to wrap, so that you can specify bounding boxes which
	// cross the ante-meridian
	ClipToBoundingBox *BoundingBox
	// Restrict AutoSuggest results to a circle, specified by center (pair of coordinates) and radius in Km.
	// For convenience, longitude is allowed to wrap around 180 degrees.
	// For example 181 is equivalent to -179.
	ClipToCircle *Circle
	// Restrict AutoSuggest results to a polygon, list of coordinates.
	// The polygon should be closed, i.e. the first element should be repeated as the last element;
	// also the list should contain at least 4 entries.
	// The API is currently limited to accepting up to 25 pairs.
	ClipToPolygon Polygon
	// For normal text input, used to specify a fallback language which will help guide AutoSuggest if the input
	// is particularly messy. If specified, this parameter must be a supported 3 word address language
	// as an ISO 639-1 2 letter code. For Bosnian-Croatian-Montenegrin-Serbian use oo.
	Language string
	// Makes AutoSuggest prefer results on land to those in the sea. This setting is on by default.
	// Use false to disable this setting and receive more suggestions in the sea.
	// Set easily using `core.Bool(false)`
	PreferLand *bool
	// Locale to specify a variant of a language
	Locale string
	// The number of AutoSuggest results to return (max 100).
	// Set easily using `core.Int(10)`
	NResults *int
	// Number of results within the results set which will have a focus
	// Set easily using `core.Int(10)`
	NFocusResult *int
}

func (aso AutoSuggestOpts) asOptionsMap() map[string]string {
	mapOpts := make(map[string]string)
	if aso.Focus != nil {
		mapOpts["focus"] = aso.Focus.AsQueryParam()
	}
	if len(aso.ClipToCountry) > 0 {
		mapOpts["clip-to-country"] = strings.Join(aso.ClipToCountry, ",")
	}
	if aso.ClipToBoundingBox != nil {
		mapOpts["clip-to-bounding-box"] = aso.ClipToBoundingBox.asQueryParam()
	}
	if aso.ClipToCircle != nil {
		mapOpts["clip-to-circle"] = aso.ClipToCircle.asQueryParam()
	}

	if len(aso.ClipToPolygon) > 0 {
		pointsStr := make([]string, 0, len(aso.ClipToPolygon))
		for _, point := range aso.ClipToPolygon {
			pointsStr = append(pointsStr, point.AsQueryParam())
		}
		mapOpts["clip-to-polygon"] = strings.Join(pointsStr, ",")
	}
	if aso.Language != "" {
		mapOpts["language"] = aso.Language
	}
	if aso.PreferLand != nil {
		mapOpts["prefer-land"] = strconv.FormatBool(*aso.PreferLand)
	}
	if aso.Locale != "" {
		mapOpts["locale"] = aso.Locale
	}
	if aso.NResults != nil {
		mapOpts["n-results"] = strconv.Itoa(*aso.NResults)
	}
	if aso.NFocusResult != nil {
		mapOpts["n-focus-result"] = strconv.Itoa(*aso.NFocusResult)
	}

	return mapOpts
}

type AutoSuggestSuggestion struct {
	Country           string `json:"country"`
	NearestPlace      string `json:"nearestPlace"`
	Words             string `json:"words"`
	DistanceToFocusKm int    `json:"distanceToFocusKm"`
	Rank              int    `json:"rank"`
	Language          string `json:"language"`
	Locale            string `json:"locale"`
}

// AutoSuggestGeoJsonResponse models the response recieved
// from the what3words public api autosuggest endpoint
type AutoSuggestResponse struct {
	Suggestions []AutoSuggestSuggestion `json:"suggestions,omitempty"`
}

type AutoSuggestWithCoordinatesSuggestion struct {
	AutoSuggestSuggestion
	Coordinates Coordinates `json:"coordinates"`
	Square      Sqaure      `json:"square"`
	MapURL      string      `json:"map"`
}

type AutoSuggestWithCoordinatesResponse struct {
	Suggestions []AutoSuggestWithCoordinatesSuggestion `json:"suggestions,omitempty"`
}

type autoSuggestResponse struct {
	AutoSuggestResponse
	Error *ErrorResponse `json:"error"`
}

func (asr autoSuggestResponse) GetError() error {
	return asr.Error
}

type autoSuggestWithCoordinatesResponse struct {
	AutoSuggestWithCoordinatesResponse
	Error *ErrorResponse `json:"error"`
}

func (asr autoSuggestWithCoordinatesResponse) GetError() error {
	return asr.Error
}

// Grid Section API
// GridSectionJsonResponse models the response recieved when
// format set to json is provided by the /v3/grid-section endpoint
// of the what3words public api.
type GridSectionJsonResponse struct {
	Lines []struct {
		Start Coordinates `json:"start"`
		End   Coordinates `json:"end"`
	} `json:"lines"`
}

// GridSectionGeoJsonResponse models the response recieved when
// format set to json is provided by the /v3/grid-section endpoint
// of the what3words public api.
type GridSectionGeoJsonResponse struct {
	Features []struct {
		Geometry struct {
			Coordinates  [][][]float64 `json:"coordinates"`
			GeometryType string        `json:"type"`
		} `json:"geometry"`
		FeatureType string            `json:"type"`
		Properties  map[string]string `json:"properties"`
	} `json:"features,omitempty"`
	Type string `json:"type,omitempty"`
}

type gridSectionResponse struct {
	*GridSectionJsonResponse
	*GridSectionGeoJsonResponse
	Error *ErrorResponse `json:"error"`
}

func (gr gridSectionResponse) GetError() error {
	return gr.Error
}

// GridSectionResponse encapsulates 2 pointers to
// Json and GeoJson response structs, based on the
// Format provided, one will be set while other
// set to nil
type GridSectionResponse struct {
	Json    *GridSectionJsonResponse
	GeoJson *GridSectionGeoJsonResponse
}

// Available Languages API
type Language struct {
	NativeName string `json:"nativeName"`
	Code       string `json:"code"`
	Name       string `json:"name"`
}

// AvailableLanguagesResponse models the response recieved
// from the what3words public api available-languages endpoint.
type AvailableLanguagesResponse struct {
	Languages []Language `json:"languages"`
}

type availableLanguagesResponse struct {
	AvailableLanguagesResponse
	Error *ErrorResponse `json:"error"`
}

func (alr availableLanguagesResponse) GetError() error {
	return alr.Error
}
