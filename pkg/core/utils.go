package core

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/what3words/w3w-go-wrapper/internal/client"
)

// MakeGetRequest makes a GET request to the specified URL.
// Reponses are unmarshalled into the response parameter, it
// is expected that the response parameter is a pointer to a struct
// which implements the ResponseErrorReader interface.
func MakeGetRequest(
	ctx context.Context,
	client client.HttpClient,
	baseURL string,
	queryParams map[string]string,
	headers map[string]string,
	response ResponseReader,
	paths ...string,
) error {

	preparedURL, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	preparedURL = preparedURL.JoinPath(paths...)
	query := preparedURL.Query()
	for qk, qv := range queryParams {
		query.Set(qk, qv)
	}
	preparedURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, preparedURL.String(), nil)
	if err != nil {
		return err
	}
	for hk, hv := range headers {
		req.Header.Set(hk, hv)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, response)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return response.GetError()
	}
	return nil
}
