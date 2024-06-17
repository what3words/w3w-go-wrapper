package v3_test

import (
	"context"
	"encoding/json"
	"sync"

	"net/http"
	"os"
	"reflect"

	"testing"

	v3 "github.com/what3words/w3w-go-wrapper/pkg/apis/v3"
	"github.com/what3words/w3w-go-wrapper/pkg/core"
)

var (
	apiURL = os.Getenv("API_URL")
)

var c2cJson = `{"country":"GB","square":{"southwest":{"lng":-0.195543,"lat":51.520833},"northeast":{"lng":-0.195499,"lat":51.52086}},"nearestPlace":"Bayswater, London","coordinates":{"lng":-0.195521,"lat":51.520847},"words":"filled.count.soap","language":"en","map":"https:\/\/w3w.co\/filled.count.soap"}`
var c2cGeoJson = `{"features":[{"bbox":[-0.195543,51.520833,-0.195499,51.52086],"geometry":{"coordinates":[-0.195521,51.520847],"type":"Point"},"type":"Feature","properties":{"country":"GB","nearestPlace":"Bayswater, London","words":"filled.count.soap","language":"en","map":"https:\/\/w3w.co\/filled.count.soap"}}],"type":"FeatureCollection"}`
var c23waJson = `{"country":"GB","square":{"southwest":{"lng":-1.246252,"lat":51.751159},"northeast":{"lng":-1.246208,"lat":51.751186}},"nearestPlace":"Oxford, Oxfordshire","coordinates":{"lng":-1.24623,"lat":51.751172},"words":"pretty.needed.chill","language":"en","map":"https:\/\/w3w.co\/pretty.needed.chill"}`
var c23waGeoJson = `{"features":[{"bbox":[-1.246252,51.751159,-1.246208,51.751186],"geometry":{"coordinates":[-1.24623,51.751172],"type":"Point"},"type":"Feature","properties":{"country":"GB","nearestPlace":"Oxford, Oxfordshire","words":"pretty.needed.chill","language":"en","map":"https:\/\/w3w.co\/pretty.needed.chill"}}],"type":"FeatureCollection"}`

// var responses map[string]string = map[string]string{
// 	"/v3/convert-to-coordinates?words=filled.count.soap":                  c2cJson,
// 	"/v3/convert-to-coordinates?format=geojson&words=filled.count.soap":   c2cGeoJson,
// 	"/v3/convert-to-coordinates?words=fill.fake.fill":                     "",
// 	"/v3/convert-to-3wa?coordinates=51.520833%2C-0.195543":                c23waJson,
// 	"/v3/convert-to-3wa?coordinates=51.520833%2C-0.195543&format=geojson": c23waGeoJson,
// 	"/v3/available-languages":                                             "",
// }

// type MockClient struct {
// }

// func (mc MockClient) Do(req *http.Request) (*http.Response, error) {
// 	if req.Method != http.MethodGet {
// 		return nil, errors.New("invalid method. Test error")
// 	}
// 	url := fmt.Sprintf("%s?%s", req.URL.Path, req.URL.RawQuery)
// 	if responseBody, ok := responses[url]; !ok {
// 		return nil, errors.New("Test setup error")
// 	} else {
// 		resp := http.Response{}
// 		resp.StatusCode = http.StatusOK
// 		resp.Body = io.NopCloser(strings.NewReader(responseBody))

// 		return &resp, nil
// 	}
// }

func setupAPI(t *testing.T) *v3.API {
	apiKey := os.Getenv("X_API_KEY")
	if apiKey == "" {
		t.Fatal("ERROR: X_API_KEY is empty or not found")
	}
	httpClient := http.DefaultClient
	svc := v3.NewAPI(apiKey, v3.WithClient(httpClient), v3.WithCustomHeader("x-temp-header", "temp"), v3.WithCustomBaseURL(apiURL))
	return svc
}

func TestConvertToCoordinatesJSON(t *testing.T) {
	resp, err := setupAPI(t).ConvertToCoordinates(context.Background(), "filled.count.soap", nil)
	if err != nil {
		t.Fatalf("ERROR: Failed to get cordinates from API due to err : %v", err)
	}
	if resp.GeoJson != nil {
		t.Fatal("ERROR: GeoJson output should be set to nil since format by default is json")
	}
	if resp.Json == nil {
		t.Fatal("ERROR: Json should have value since format by default is json")
	}

	var expected v3.ConvertAPIJsonResponse
	json.Unmarshal([]byte(c2cJson), &expected)
	if !reflect.DeepEqual(expected, *resp.Json) {
		t.Fatalf("ERROR: Expected output '%v' recieved '%v'", expected, *resp.Json)
	}
}

func TestConvertToCoordinatesGeoJSON(t *testing.T) {
	resp, err := setupAPI(t).ConvertToCoordinates(context.Background(), "filled.count.soap", &v3.ConvertAPIOpts{
		Format: v3.ResponseFormatGeoJson,
	})
	if err != nil {
		t.Fatalf("ERROR: Failed to get cordinates from API due to err : %v", err)
	}
	if resp.GeoJson == nil {
		t.Fatal("ERROR: GeoJson output should be set since format has been overidden to use geojson")
	}
	if resp.Json != nil {
		t.Fatal("ERROR: Json output should be nil since format has been overidden to use geojson")
	}

	var expected v3.ConvertAPIGeoJsonResponse
	json.Unmarshal([]byte(c2cGeoJson), &expected)
	if !reflect.DeepEqual(expected, *resp.GeoJson) {
		t.Fatalf("ERROR: Expected output '%v' recieved '%v'", expected, *resp.GeoJson)
	}
}

func TestConvertToCoordinatesInvalidWords(t *testing.T) {
	_, err := setupAPI(t).ConvertToCoordinates(context.Background(), "fill.fake.fill", nil)
	if err == nil {
		t.Fatal("ERROR: error should be set to BadWords")
	}
}

func TestConvertTo3WAJSON(t *testing.T) {
	resp, err := setupAPI(t).ConvertTo3wa(context.Background(), core.Coordinates{
		Lng: -1.24623,
		Lat: 51.751172,
	}, nil)
	if err != nil {
		t.Fatalf("ERROR: Failed to get cordinates from API due to err : %v", err)
	}
	if resp.GeoJson != nil {
		t.Fatal("ERROR: GeoJson output should be set to nil since format by default is json")
	}
	if resp.Json == nil {
		t.Fatal("ERROR: Json should have value since format by default is json")
	}

	var expected v3.ConvertAPIJsonResponse
	json.Unmarshal([]byte(c23waJson), &expected)
	if !reflect.DeepEqual(expected, *resp.Json) {
		t.Fatalf("ERROR: Expected output '%v' recieved '%v'", expected.Words, resp.Json.Words)
	}
}

func TestConvertTo3WAGeoJSON(t *testing.T) {
	resp, err := setupAPI(t).ConvertTo3wa(context.Background(), core.Coordinates{
		Lng: -1.24623,
		Lat: 51.751172,
	}, &v3.ConvertAPIOpts{
		Format: v3.ResponseFormatGeoJson,
	})
	if err != nil {
		t.Fatalf("ERROR: Failed to get cordinates from API due to err : %v", err)
	}
	if resp.GeoJson == nil {
		t.Fatal("ERROR: GeoJson output should be set since format has been overidden to use geojson")
	}
	if resp.Json != nil {
		t.Fatal("ERROR: Json output should be nil since format has been overidden to use geojson")
	}

	var expected v3.ConvertAPIGeoJsonResponse
	json.Unmarshal([]byte(c23waGeoJson), &expected)
	if !reflect.DeepEqual(expected, *resp.GeoJson) {
		t.Fatalf("ERROR: Expected output '%v' recieved '%v'", expected, *resp.GeoJson)
	}
}

func TestConvertTo3WAInvalidWords(t *testing.T) {
	_, err := setupAPI(t).ConvertToCoordinates(context.Background(), "fill.fake.fill", nil)
	if err == nil {
		t.Fatal("ERROR: error should be set to BadWords")
	}
	_, ok := err.(*v3.ErrorResponse)
	if !ok {
		t.Fatal("ERROR: error should be of type ErrorResponse")
	}
}

func TestGridSectionJSON(t *testing.T) {

	resp, err := setupAPI(t).GridSection(context.Background(), v3.BoundingBox{
		SouthWest: core.Coordinates{
			Lat: 52.207988,
			Lng: 0.116126,
		},
		NorthEast: core.Coordinates{
			Lat: 52.208867,
			Lng: 0.117540,
		},
	}, nil)
	if err != nil {
		t.Fatalf("ERROR: Failed to get grid section from API - %v", err)
	}
	if err != nil {
		t.Fatalf("ERROR: Failed to get cordinates from API due to err : %v", err)
	}
	if resp.GeoJson != nil {
		t.Fatal("ERROR: GeoJson output should be set to nil since format by default is json")
	}
	if resp.Json == nil {
		t.Fatal("ERROR: Json should have value since format by default is json")
	}
}

func TestGridSectionGeoJSON(t *testing.T) {
	resp, err := setupAPI(t).GridSection(context.Background(), v3.BoundingBox{
		SouthWest: core.Coordinates{
			Lat: 52.207988,
			Lng: 0.116126,
		},
		NorthEast: core.Coordinates{
			Lat: 52.208867,
			Lng: 0.117540,
		},
	}, &v3.GridSectionOpts{
		Format: v3.ResponseFormatGeoJson,
	})
	if err != nil {
		t.Fatalf("ERROR: Failed to get grid section from API - %v", err)
	}
	if err != nil {
		t.Fatalf("ERROR: Failed to get cordinates from API due to err : %v", err)
	}
	if resp.GeoJson == nil {
		t.Fatal("ERROR: GeoJson output should be set since format has been overidden to use geojson")
	}
	if resp.Json != nil {
		t.Fatal("ERROR: Json output should be nil since format has been overidden to use geojson")
	}
}

func TestAvailableLanguages(t *testing.T) {
	_, err := setupAPI(t).AvailableLanguages(context.Background())
	if err != nil {
		t.Fatalf("ERROR: Error occured trying to retrieve languages")
	}
}

// AutoSuggest
// input - plan.clips.a
const autoSuggest = `{"suggestions":[{"country":"US","nearestPlace":"Absecon, New Jersey","words":"plan.clips.also","rank":1,"language":"en"},{"country":"US","nearestPlace":"Sunland, California","words":"plan.clips.back","rank":2,"language":"en"},{"country":"US","nearestPlace":"Keego Harbor, Michigan","words":"plan.clips.each","rank":3,"language":"en"}]}`

// focus - 51.521251,-0.203586
const autoSuggestFocus = `{"suggestions":[{"country":"GB","nearestPlace":"Brixton Hill, London","words":"plan.clips.area","rank":1,"distanceToFocusKm":11,"language":"en"},{"country":"GB","nearestPlace":"Borehamwood, Hertfordshire","words":"plan.clips.arts","rank":2,"distanceToFocusKm":16,"language":"en"},{"country":"GB","nearestPlace":"Wood Green, London","words":"plan.slips.cage","rank":3,"distanceToFocusKm":13,"language":"en"}]}`

// clip-to-country - NZ,AU
const autoSuggestClipToCountry = `{"suggestions":[{"country":"AU","nearestPlace":"Emerald, Queensland","words":"plan.clips.bias","rank":1,"language":"en"},{"country":"AU","nearestPlace":"Kumpupintil, Western Australia","words":"plan.clips.atop","rank":2,"language":"en"},{"country":"AU","nearestPlace":"Melville, Western Australia","words":"plan.clips.clad","rank":3,"language":"en"}]}`

// clip-to-bounding-box - 51.521,-0.343,52.6,2.3324
const autoSuggestClipToBoundingBox = `{"suggestions":[{"country":"GB","nearestPlace":"Borehamwood, Hertfordshire","words":"plan.clips.arts","rank":1,"language":"en"},{"country":"GB","nearestPlace":"Cambridge, Cambridgeshire","words":"plan.clips.boat","rank":2,"language":"en"},{"country":"GB","nearestPlace":"Rayleigh, Essex","words":"plan.clip.deal","rank":3,"language":"en"}]}`

// clip-to-circle - 51.521,-0.343,142
const autoSuggestClipToCirle = `{"suggestions":[{"country":"GB","nearestPlace":"Brixton Hill, London","words":"plan.clips.area","rank":1,"language":"en"},{"country":"GB","nearestPlace":"Borehamwood, Hertfordshire","words":"plan.clips.arts","rank":2,"language":"en"},{"country":"GB","nearestPlace":"Cambridge, Cambridgeshire","words":"plan.clips.boat","rank":3,"language":"en"}]}`

// clip-to-polygon=51.521,-0.343,52.6,2.3324,54.234,8.343,51.521,-0.343
const autoSuggestClipToPolygon = `{"suggestions":[{"country":"GB","nearestPlace":"High Ongar, Essex","words":"plan.clip.bags","rank":1,"language":"en"},{"country":"GB","nearestPlace":"Wood Green, London","words":"plan.slips.cage","rank":2,"language":"en"},{"country":"GB","nearestPlace":"High Ongar, Essex","words":"plan.flips.ants","rank":3,"language":"en"}]}`

// prefer-land - false
const autoSuggestPreferLandFalse = `{"suggestions":[{"country":"US","nearestPlace":"Absecon, New Jersey","words":"plan.clips.also","rank":1,"language":"en"},{"country":"US","nearestPlace":"Sunland, California","words":"plan.clips.back","rank":2,"language":"en"},{"country":"US","nearestPlace":"Keego Harbor, Michigan","words":"plan.clips.each","rank":3,"language":"en"}]}`

// multiple policies
// clip-to-circle - 51.521,-0.343,142
// clip-to-polygon - 51.521,-0.343,52.6,2.3324,54.234,8.343,51.521,-0.343
// clip-to-bounding-box - 51.521,-0.343,52.6,2.3324
// clip-to-country=GB
// language=eng
const autoSuggestMultiplePolicies = `{"suggestions":[{"country":"GB","nearestPlace":"High Ongar, Essex","words":"plan.clip.bags","rank":1,"language":"en"},{"country":"GB","nearestPlace":"Wood Green, London","words":"plan.slips.cage","rank":2,"language":"en"},{"country":"GB","nearestPlace":"High Ongar, Essex","words":"plan.flips.ants","rank":3,"language":"en"}]}`

func TestAutoSuggest(t *testing.T) {
	svc := setupAPI(t)
	preferLand := false
	asTests := []struct {
		name     string
		opts     *v3.AutoSuggestOpts
		expected string
	}{
		{
			"OnlyInputs",
			nil,
			autoSuggest,
		},
		{
			"Focus",
			&v3.AutoSuggestOpts{
				Focus: &core.Coordinates{
					Lat: 51.521251,
					Lng: -0.203586,
				},
			},
			autoSuggestFocus,
		},
		{
			"ClipToCountry",
			&v3.AutoSuggestOpts{
				ClipToCountry: []string{"NZ", "AU"},
			},
			autoSuggestClipToCountry,
		},
		{
			"ClipToBoundingBox",
			&v3.AutoSuggestOpts{
				ClipToBoundingBox: &v3.BoundingBox{
					core.Coordinates{
						Lat: 51.521,
						Lng: -0.343,
					},
					core.Coordinates{
						Lat: 52.6,
						Lng: 2.3324,
					},
				},
			},
			autoSuggestClipToBoundingBox,
		},
		{
			"ClipToCircle",
			&v3.AutoSuggestOpts{
				ClipToCircle: &v3.Circle{
					Center: core.Coordinates{
						Lat: 51.521,
						Lng: -0.343,
					},
					RadiusKm: 142,
				},
			},
			autoSuggestClipToCirle,
		},
		{
			"ClipToPolygon",
			&v3.AutoSuggestOpts{
				ClipToPolygon: []core.Coordinates{
					{
						Lat: 51.521,
						Lng: -0.343,
					},
					{
						Lat: 52.6,
						Lng: 2.3324,
					},
					{
						Lat: 54.234,
						Lng: 8.343,
					},
					{
						Lat: 51.521,
						Lng: -0.343,
					},
				},
			},
			autoSuggestClipToPolygon,
		},
		{
			"PreferLandFalse",
			&v3.AutoSuggestOpts{
				PreferLand: &preferLand,
			},
			autoSuggestPreferLandFalse,
		},
		{
			"MultiplePolicies",
			&v3.AutoSuggestOpts{
				ClipToBoundingBox: &v3.BoundingBox{
					core.Coordinates{
						Lat: 51.521,
						Lng: -0.343,
					},
					core.Coordinates{
						Lat: 52.6,
						Lng: 2.3324,
					},
				},
				ClipToCountry: []string{"GB"},
				ClipToPolygon: []core.Coordinates{
					{
						Lat: 51.521,
						Lng: -0.343,
					},
					{
						Lat: 52.6,
						Lng: 2.3324,
					},
					{
						Lat: 54.234,
						Lng: 8.343,
					},
					{
						Lat: 51.521,
						Lng: -0.343,
					},
				},
				ClipToCircle: &v3.Circle{
					Center: core.Coordinates{
						Lat: 51.521,
						Lng: -0.343,
					},
					RadiusKm: 142,
				},
				Language: "en",
			},
			autoSuggestMultiplePolicies,
		},
	}

	for _, asTest := range asTests {
		test := asTest
		t.Run(asTest.name, func(t *testing.T) {
			t.Parallel()
			resp, err := svc.AutoSuggest(context.Background(), "plan.clips.a", test.opts)
			if err != nil {
				t.Fatalf("ERROR: Failed to get auto suggest from API - %v", err)
			}
			var expected v3.AutoSuggestResponse
			json.Unmarshal([]byte(test.expected), &expected)

			buffer, err := json.Marshal(resp)
			if err != nil {
				t.Fatalf("ERROR: Failed to marshal response - %v", err)
			}

			if !reflect.DeepEqual(expected, *resp) {
				t.Fatalf("ERROR: Expected output '%s' recieved '%s'", test.expected, string(buffer))
			}
		})
	}
}

// Test thread safety
func TestThreadSafety(t *testing.T) {
	svc := setupAPI(t)
	var wg sync.WaitGroup
	count := 100
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			_, err := svc.ConvertToCoordinates(context.Background(), "filled.count.soap", nil)
			if err != nil {
				t.Errorf("ERROR: Failed to get coordinates from API - %v", err)
			}
		}()
	}
	wg.Wait()
}
