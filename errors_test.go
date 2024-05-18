package pokeapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nightmarlin/pokeapi"
)

func Test_ErrNotFound_IsHTTP404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(http.NotFound))
	t.Cleanup(ts.Close)

	expectErr := pokeapi.ErrNotFound

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("request generation error must succeed, but failed with: %v", err)
	}
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("request must succeed, but failed with: %v", err)
	}

	gotErr := pokeapi.NewHTTPError(resp)

	if gotErr != expectErr {
		t.Errorf("want error to be %v; got %v", expectErr, gotErr)
	}
}
