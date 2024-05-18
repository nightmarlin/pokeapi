package pokeapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_do(t *testing.T) {
	t.Parallel()

	t.Run(
		"fetches value on cache miss",
		func(t *testing.T) {
			t.Parallel()

			const pkPath = "/pokemon/togedemaru"
			var expectRes = echoResp{Path: pkPath}

			given(t, echoHandler).
				get(pkPath).
				verify(
					thatThereWasNoError,
					thatResponseIs(expectRes),
					thatNCacheLookupsOccurred(1),
					thatCacheWasWrittenTo(pkPath, expectRes),
				)
		},
	)

	t.Run(
		"uses cached value on hit",
		func(t *testing.T) {
			t.Parallel()

			const pkPath = "/pokemon/togedemaru"
			var expectRes = echoResp{Path: pkPath}

			given(t, http.NotFound, withCachedValue(pkPath, expectRes)).
				get(pkPath).
				verify(
					thatThereWasNoError,
					thatResponseIs(expectRes),
					thatNCacheLookupsOccurred(1),
					thatCacheWasNotWrittenTo,
				)
		},
	)

	t.Run(
		"returns expected error on server error response",
		func(t *testing.T) {
			t.Parallel()

			const pkPath = "/pokemon/dedenne"

			given(t, http.NotFound).
				get(pkPath).
				verify(
					thatErrorIs(ErrNotFound),
					thatNCacheLookupsOccurred(1),
					thatCacheWasNotWrittenTo,
				)
		},
	)
}

// region test helpers

// A recordingCache records the lookups performed on it. it is unsafe for concurrent use.
type recordingCache struct {
	cachedValues   map[string]any
	lookups        []string
	recordedValues map[string]any
}

func (r *recordingCache) Lookup(
	ctx context.Context,
	url string,
	loadOnMiss CacheLoader,
) (any, error) {
	r.lookups = append(r.lookups, url)

	if v, ok := r.cachedValues[url]; ok {
		return v, nil
	}

	v, err := loadOnMiss(ctx)
	if err != nil {
		return nil, err
	}

	if r.recordedValues == nil {
		r.recordedValues = make(map[string]any)
	}
	r.recordedValues[url] = v
	return v, nil
}

type echoResp struct {
	Path string `json:"path"`
}

type doSUT struct {
	t *testing.T

	server *httptest.Server
	cache  recordingCache

	done bool
	resp echoResp
	err  error
}

type (
	setup    func(s *doSUT)
	verifier func(s *doSUT)
)

func given(t *testing.T, serverFn http.HandlerFunc, setups ...setup) *doSUT {
	t.Helper()

	s := doSUT{t: t, server: httptest.NewServer(serverFn)}
	t.Cleanup(s.server.Close)

	for _, setup := range setups {
		setup(&s)
	}
	return &s
}

func (s *doSUT) endpointToURL(endpoint string) string {
	return fmt.Sprintf("%s/%s", s.server.URL, trimSlash(endpoint))
}

func withCachedValue(endpoint string, v any) setup {
	return func(s *doSUT) {
		if s.cache.cachedValues == nil {
			s.cache.cachedValues = make(map[string]any)
		}
		s.cache.cachedValues[s.endpointToURL(endpoint)] = v
	}
}

// the echoHandler writes an echoResp with Path set to the /path of the request
// encoded in JSON.
func echoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(echoResp{Path: r.URL.Path})
}

func (s *doSUT) get(endpoint string) *doSUT {
	s.t.Helper()

	if s.done {
		s.t.Fatalf("s.on() called multiple times")
	}

	s.resp, s.err = do[echoResp](
		context.Background(),
		NewClient(
			&ClientOpts{
				Cache:       &s.cache,
				PokeAPIRoot: s.server.URL,
				HTTPClient:  s.server.Client(),
			},
		),
		s.endpointToURL(endpoint),
		nil,
	)

	s.done = true
	return s
}

func (s *doSUT) verify(vs ...verifier) {
	s.t.Helper()

	if !s.done {
		s.t.Fatalf("s.verify() called before s.on()")
	}

	for _, v := range vs {
		v(s)
	}
}

func thatThereWasNoError(s *doSUT) {
	if s.err != nil {
		s.t.Errorf("want returned error to be nil; got %v", s.err)
	}
}

func thatErrorIs(targetErr error) verifier {
	return func(s *doSUT) {
		if !errors.Is(s.err, targetErr) {
			s.t.Errorf("want errors.Is(%v, %v) == true; got false", s.err, targetErr)
		}

		if s.resp != (echoResp{}) {
			s.t.Errorf("want response to be zero; got %v", s.resp)
		}
	}
}

func thatResponseIs(wantResp echoResp) verifier {
	return func(s *doSUT) {
		if s.resp != wantResp {
			s.t.Errorf("want response to be %v; got %v", wantResp, s.resp)
		}
	}
}

func thatNCacheLookupsOccurred(n int) verifier {
	return func(s *doSUT) {
		if got := len(s.cache.lookups); got != n {
			s.t.Errorf("want %d cache lookups to happen; got %d", n, got)
		}
	}
}

func thatCacheWasNotWrittenTo(s *doSUT) {
	if got := len(s.cache.recordedValues); got != 0 {
		s.t.Errorf("want 0 cache writes to occur; got %d", got)
	}
}

func thatCacheWasWrittenTo(endpoint string, wantValue any) verifier {
	return func(s *doSUT) {
		v, ok := s.cache.recordedValues[s.endpointToURL(endpoint)]
		if !ok {
			s.t.Errorf("want cache write for endpoint %q, but none found", endpoint)
		} else if v != wantValue {
			s.t.Errorf("want cached value for endpoint %q to be %v; got %v", endpoint, wantValue, v)
		}
	}
}

// endregion
