package core_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/what3words/w3w-go-wrapper/pkg/core"
)

type FakeResponse map[string]interface{}

func (fk FakeResponse) GetError() error {
	return nil
}

func TestMakeRequest(t *testing.T) {
	var fk FakeResponse
	err := core.MakeGetRequest(
		context.Background(),
		http.DefaultClient,
		"https://httpbin.org",
		map[string]string{
			"random": "123",
		},
		map[string]string{
			"accept": "application/json",
		},
		&fk,
		"json",
	)
	if err != nil {
		t.Fatalf("ERROR: Got error %v", err)
	}
}
