package w3wgowrapper_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	w3w "github.com/what3words/w3w-go-wrapper"
	v3 "github.com/what3words/w3w-go-wrapper/pkg/apis/v3"
)

var (
	apiURL = os.Getenv("API_URL")
)

type MockAPI struct {
	t *testing.T
}

func (m MockAPI) Do(req *http.Request) (*http.Response, error) {

	if !strings.HasPrefix(req.URL.String(), apiURL) {
		m.t.Fatalf("ERROR: unexpected API URL: %s", req.URL.String())
	}

	if req.Header.Get("x-temp-header") != "temp" {
		m.t.Fatalf("ERROR: unexpected header: %s", req.Header.Get("x-temp-header"))
	}

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("{}")),
	}, nil
}

func setupSvc(t *testing.T) w3w.Service {
	apiKey := os.Getenv("X_API_KEY")
	if apiKey == "" {
		t.Fatal("ERROR: X_API_KEY is empty or not found")
	}
	httpClient := MockAPI{t}
	svc := w3w.NewService(apiKey, w3w.WithClient(httpClient), w3w.WithCustomHeader("x-temp-header", "temp"), w3w.WithCustomBaseURL(apiURL))
	return svc
}

func TestAPI(t *testing.T) {
	w3wAPI := setupSvc(t)
	_, err := w3wAPI.V3().AvailableLanguages(context.Background())
	if err != nil {
		t.Fatalf("ERROR: %+v", err)
	}
}

func TestFindPossible3wa(t *testing.T) {
	w3wAPI := setupSvc(t)
	source := "Can be found at filled.count.soap and at ///test.fake.words but not at test.fake. or test.fake"
	pa := w3wAPI.FindPossible3wa(source)
	expected := []string{"filled.count.soap", "test.fake.words"}
	if len(pa) != len(expected) {
		t.Fatalf("ERROR: expected to find %d possible 3 word addresses, but found %d", len(expected), len(pa))
	}
	if !reflect.DeepEqual(pa, expected) {
		t.Fatalf("ERROR: expected to find %v, but found %v", expected, pa)
	}
}

func ExampleService() {
	apiKey := os.Getenv("X_API_KEY")
	if apiKey == "" {
		panic("ERROR: X_API_KEY is empty or not found")
	}
	svc := w3w.NewService(apiKey)
	_, err := svc.V3().AvailableLanguages(context.Background())
	if err != nil {
		panic(err)
	}
}

func ExampleConvertToCoordinates() {
	apiKey := os.Getenv("X_API_KEY")
	if apiKey == "" {
		panic("ERROR: X_API_KEY is empty or not found")
	}
	svc := w3w.NewService(apiKey)
	resp, err := svc.V3().ConvertToCoordinates(context.Background(), "filled.count.soap", nil)
	if err != nil {
		panic(err)
	}
	// Json is the default output format used.
	fmt.Println(resp.Coordinates)
}

func ExampleConvertToCoordinatesLanguage() {
	apiKey := os.Getenv("X_API_KEY")
	if apiKey == "" {
		panic("ERROR: X_API_KEY is empty or not found")
	}
	svc := w3w.NewService(apiKey)
	resp, err := svc.V3().ConvertToCoordinates(context.Background(), "تجتمع.ضباط.ثقافية", &v3.ConvertAPIOpts{
		Language: "ar",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Coordinates)
}

func ExampleConvertToCoordinatesGeoJson() {
	apiKey := os.Getenv("X_API_KEY")
	if apiKey == "" {
		panic("ERROR: X_API_KEY is empty or not found")
	}
	svc := w3w.NewService(apiKey)
	resp, err := svc.V3().ConvertToCoordinatesGeoJson(context.Background(), "filled.count.soap", nil)
	if err != nil {
		panic(err)
	}
	// If format GeoJson is not set in options, GeoJson attribute of the response will be set to nil
	fmt.Println(resp.Features[0].Geometry.Coordinates)
}

func ExampleConvertTo3wa() {
	apiKey := os.Getenv("X_API_KEY")
	if apiKey == "" {
		panic("ERROR: X_API_KEY is empty or not found")
	}
	svc := w3w.NewService(apiKey)
	resp, err := svc.V3().ConvertTo3wa(context.Background(), v3.Coordinates{
		Lat: 51.520847,
		Lng: -0.195521,
	}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Words)
}

func ExampleGridSection() {
	apiKey := os.Getenv("X_API_KEY")
	if apiKey == "" {
		panic("ERROR: X_API_KEY is empty or not found")
	}
	svc := w3w.NewService(apiKey)
	resp, err := svc.V3().GridSection(context.Background(), v3.BoundingBox{
		SouthWest: v3.Coordinates{
			Lat: 52.207988,
			Lng: 0.116126,
		},
		NorthEast: v3.Coordinates{
			Lat: 52.208867,
			Lng: 0.117540,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Lines[0])
}

func ExampleAutoSuggest() {
	apiKey := os.Getenv("X_API_KEY")
	if apiKey == "" {
		panic("ERROR: X_API_KEY is empty or not found")
	}
	svc := w3w.NewService(apiKey)
	resp, err := svc.V3().AutoSuggest(context.Background(), "filled.count.so", &v3.AutoSuggestOpts{
		ClipToCircle: &v3.Circle{
			Center: v3.Coordinates{
				Lat: 51.520847,
				Lng: -0.195521,
			},
			RadiusKm: 10,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Suggestions[0].Words)
}

func ExampleExactWhat3WordsAPIErrors() {
	apiKey := os.Getenv("X_API_KEY")
	if apiKey == "" {
		panic("ERROR: X_API_KEY is empty or not found")
	}
	svc := w3w.NewService(apiKey)
	_, err := svc.V3().ConvertToCoordinates(context.Background(), "filled", nil)
	if err != nil {
		if err, ok := err.(*v3.ErrorResponse); ok {
			// Refer v3.ErrorCode for more types of errors
			fmt.Println(err.Code)
		}
	}
}
