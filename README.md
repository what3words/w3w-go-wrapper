# <img src="https://what3words.com/assets/images/w3w_square_red.png" width="64" height="64" alt="what3words">&nbsp;w3w-go-wrapper
[![Go Report Card](https://goreportcard.com/badge/github.com/what3words/w3w-go-wrapper)](https://goreportcard.com/report/github.com/what3words/w3w-go-wrapper)

A Go library to make requests to the [what3words REST API](https://developer.what3words.com/public-api/). See the what3words public API [documentation](https://developer.what3words.com/public-api/docs) for more information about how to use our REST API.

## Overview

The what3words Golang wrapper gives you programmatic access to:

- [Convert a 3 word address to coordinates](https://developer.what3words.com/public-api/docs#convert-to-coords)
- [Convert coordinates to a 3 word address](https://developer.what3words.com/public-api/docs#convert-to-3wa)
- [Autosuggest functionality which takes a slightly incorrect 3 word address, and suggests a list of valid 3 word addresses](https://developer.what3words.com/public-api/docs#autosuggest)
- [Obtain a section of the 3m x 3m what3words grid for a bounding box.](https://developer.what3words.com/public-api/docs#grid-section)
- [Determine the currently support 3 word address languages.](https://developer.what3words.com/public-api/docs#available-languages)

## Install

```sh
go get github.com/what3words/w3w-go-wrapper
```

## Usage

```go
package main

import (
    "context"
    "fmt"

	"github.com/what3words/w3w-go-wrapper/pkg/apis/v3"
    w3w "github.com/what3words/w3w-go-wrapper"
)

func main() {
    apiKey := "<YOUR_API_KEY>"
    svc := w3w.NewService(apiKey)
    // Using a custom api endpoint in cases if your using the enterprise server
    // svc := w3w.NewService(apiKey, WithCustomBaseURL("enterprise.sever.domain"))
}
```

## Documentation

> NOTE: All functions and structures part of the w3w-go-wrapper library are fully documented using godoc compatible in-line documentation

### what3words Service

The `w3w-go-wrapper.Service` provides a quick and easy way to instantiate the client that can be used to make requests against the what3words API. It also provides helper functions for setting API configuration across all versions of the What3Words API.

## Examples

### Autosuggest

```go
package main

import (
    "context"
    "fmt"

    w3w "github.com/what3words/w3w-go-wrapper"
	"github.com/what3words/w3w-go-wrapper/pkg/apis/v3"
    
)

func main() {
    apiKey := "<YOUR_API_KEY>"
    svc := w3w.NewService(apiKey)

    // Selected option clip to circle, multiple options can be selected, Refer https://developer.what3words.com/public-api/docs#autosuggest for options and limitations. Pass nil if options are not required
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
    fmt.Println(resp)

}
```

### Convert To Coordinates

```go
package main

import (
    "fmt"
    "context"
    "fmt"

    w3w "github.com/what3words/w3w-go-wrapper"
	"github.com/what3words/w3w-go-wrapper/pkg/apis/v3"
    
)

func main() {
    apiKey := "<YOUR_API_KEY>"
    svc := w3w.NewService(apiKey)

    resp, err := svc.V3().ConvertToCoordinates(context.Background(), "filled.count.soap", nil)
    if err != nil {
        panic(err)
    }
    // By default response JSON is used
    fmt.Println(resp.Coordinates)

    // Getting a geojson response
    geoResp, err := svc.V3().ConvertToCoordinatesGeoJson(context.Background(), "filled.count.soap", nil)
    if err != nil {
        panic(err)
    }
	fmt.Println(geoResp.Features[0].Geometry.Coordinates)
}
```

### Convert to Three Word Address

```go
package main

import (
    "context"
    "fmt"

    w3w "github.com/what3words/w3w-go-wrapper"
	"github.com/what3words/w3w-go-wrapper/pkg/apis/v3"
    
)

func main() {
    apiKey := "<YOUR_API_KEY>"
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
```

### Grid Section

```go
package main

import (
    "context"
    "fmt"

    w3w "github.com/what3words/w3w-go-wrapper"
	"github.com/what3words/w3w-go-wrapper/pkg/apis/v3"
    
)

func main() {
    apiKey := "<YOUR_API_KEY>"
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
	fmt.Println(resp.Lines)

}
```

> **The requested box must not exceed 4km from corner to corner, or a BadBoundingBoxTooBig error will be returned. Latitudes must be >= -90 and <= 90, but longitudes are allowed to wrap around 180. To specify a bounding-box that crosses the anti-meridian, use longitude greater than 180.**

### Available Languages

```go
package main

import (
    "context"
    "fmt"

    w3w "github.com/what3words/w3w-go-wrapper"
	"github.com/what3words/w3w-go-wrapper/pkg/apis/v3"
    
)

func main() {
    apiKey := "<YOUR_API_KEY>"
    svc := w3w.NewService(apiKey)
	resp, err := svc.V3().AvailableLanguages(context.Background())
	if err != nil {
		panic(err)
	}
    fmt.Println(resp)
}
```

### Find Possible 3 Word Addresses

FindPossible3wa searches the string passed in for all substrings in the form of a three word address.

```go
package main

import (
	"fmt"

	w3w "github.com/what3words/w3w-go-wrapper"
)

func main() {
	apiKey := "<YOUR_API_KEY>"
	svc := w3w.NewService(apiKey)
	text := `This is a valid 3 word address filled.count.soap`
	pa := svc.FindPossible3wa(text)
	fmt.Println(pa)
}
```

### Is Possible 3 Word Address

IsPossible3wa determines if the string passed in is in the form of a three word address.

```go
package main

import (
	"fmt"

	w3w "github.com/what3words/w3w-go-wrapper"
)

func main() {
	apiKey := "<YOUR_API_KEY>"
	svc := w3w.NewService(apiKey)
	pa := svc.IsPossible3wa("filled.count.fake")
	fmt.Println(pa)
}
```

### Did you Mean

DidYouMean determines if the string passed in is almost in the form of a three word address.

```go
package main

import (
	"fmt"

	w3w "github.com/what3words/w3w-go-wrapper"
)

func main() {
	apiKey := "<YOUR_API_KEY>"
	svc := w3w.NewService(apiKey)
	pa := svc.DidYouMean("filled-count-fake")
	fmt.Println(pa)
}
```

### Is Valid 3 word address

IsValid3wa validates the given string as a real three-word address by making a call to the API. The context can be used to cancel the underlying call.

```go
package main

import (
	"context"
	"fmt"

	w3w "github.com/what3words/w3w-go-wrapper"
)

func main() {
	apiKey := "<YOUR_API_KEY>"
	svc := w3w.NewService(apiKey)
	pa := svc.IsValid3wa(context.Background(), "filled.count.soap")
	fmt.Println(pa)
}
```