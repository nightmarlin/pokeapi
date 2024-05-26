package pokeapi_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nightmarlin/pokeapi"
)

func Test_ErrNotFound_IsHTTP404(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ts := httptest.NewServer(http.HandlerFunc(http.NotFound))
	t.Cleanup(ts.Close)

	expectErr := pokeapi.ErrNotFound

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("request generation error must succeed, but failed with: %v", err)
	}
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("request must succeed, but failed with: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	gotErr := pokeapi.NewHTTPError(resp)

	if !errors.Is(gotErr, expectErr) {
		t.Errorf("want error to be %v; got %v", expectErr, gotErr)
	}
}
