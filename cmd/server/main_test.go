package main

import (
	"strings"
	"testing"

	"github.com/connorvoisey/shgrid_api/pkg/load"
	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestGetGreeting(t *testing.T) {
    state, err := load.Init()
    if err != nil{
        panic(err)
    }
	_, api := humatest.New(t)

	addRoutes(api, state.Log, state.Db)

	resp := api.Get("/greeting/world")
	if !strings.Contains(resp.Body.String(), "Hello, world!") {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestPutReview(t *testing.T) {
    state, err := load.Init()
    if err != nil{
        panic(err)
    }
	_, api := humatest.New(t)

	addRoutes(api, state.Log, state.Db)

	resp := api.Post("/reviews", map[string]any{
		"author": "daniel",
		"rating": 5,
	})

	if resp.Code != 201 {
		t.Fatalf("Unexpected status code: %d", resp.Code)
	}
}

func TestPutReviewError(t *testing.T) {
    state, err := load.Init()
    if err != nil{
        panic(err)
    }
	_, api := humatest.New(t)

	addRoutes(api, state.Log, state.Db)

	resp := api.Post("/reviews", map[string]any{
		"rating": 10,
	})

	if resp.Code != 422 {
		t.Fatalf("Unexpected status code: %d", resp.Code)
	}
}
