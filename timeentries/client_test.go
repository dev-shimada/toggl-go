package timeentries_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/dev-shimada/toggl-go/timeentries"
	"github.com/google/go-cmp/cmp"
)

func TestNewClient(t *testing.T) {
	want := timeentries.Client{
		HttpClient: &http.Client{},
		Token:      "token",
	}
	got := timeentries.NewClient("token")
	if !cmp.Equal(want, got) {
		t.Errorf("diff: %v", cmp.Diff(want, got))
	}
}

func TestGet(t *testing.T) {
	want := http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme:   "https",
			Host:     "api.track.toggl.com",
			User:     url.UserPassword("token", "api_token"),
			RawQuery: "key=value",
		},
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}
	client := timeentries.Client{
		Token: "token",
	}
	got := client.Get(url.URL{RawQuery: "key=value"})

	if want.Method != got.Method {
		t.Errorf("want: %v, got: %v", want.Method, got.Method)
	}
	if want.URL.RawQuery != got.URL.RawQuery {
		t.Errorf("wnt: %v, got: %v", want.URL.RawQuery, got.URL.RawQuery)
	}
	if !cmp.Equal(want.Header, got.Header) {
		t.Errorf("diff: %v", cmp.Diff(want.Header, got.Header))
	}
}

func TestPost(t *testing.T) {
	bodyJson, err := json.Marshal(struct{ Text string }{"body"})
	if err != nil {
		t.Fatal(err)
	}

	want := http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme:   "https",
			Host:     "api.track.toggl.com",
			User:     url.UserPassword("token", "api_token"),
			RawQuery: "key=value",
		},
		Body: io.ReadCloser(io.NopCloser(bytes.NewBuffer(bodyJson))),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}
	client := timeentries.Client{
		Token: "token",
	}
	got := client.Post(url.URL{RawQuery: "key=value"}, bodyJson)

	if want.Method != got.Method {
		t.Errorf("want: %v, got: %v", want.Method, got.Method)
	}
	if want.URL.RawQuery != got.URL.RawQuery {
		t.Errorf("wnt: %v, got: %v", want.URL.RawQuery, got.URL.RawQuery)
	}

	wantBody := make([]byte, len(bodyJson))
	gotBody := make([]byte, len(bodyJson))
	if _, err := want.Body.Read(wantBody); err != nil {
		t.Fatal(err)
	}
	if _, err := got.Body.Read(gotBody); err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(string(wantBody), string(gotBody)) {
		t.Errorf("diff: %v", cmp.Diff(wantBody, gotBody))
	}
	if !cmp.Equal(want.Header, got.Header) {
		t.Errorf("diff: %v", cmp.Diff(want.Header, got.Header))
	}
}

func TestPatch(t *testing.T) {
	bodyJson, err := json.Marshal(struct{ Text string }{"body"})
	if err != nil {
		t.Fatal(err)
	}

	want := http.Request{
		Method: http.MethodPatch,
		URL: &url.URL{
			Scheme:   "https",
			Host:     "api.track.toggl.com",
			User:     url.UserPassword("token", "api_token"),
			RawQuery: "key=value",
		},
		Body: io.ReadCloser(io.NopCloser(bytes.NewBuffer(bodyJson))),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}
	client := timeentries.Client{
		Token: "token",
	}
	got := client.Patch(url.URL{RawQuery: "key=value"}, bodyJson)

	if want.Method != got.Method {
		t.Errorf("want: %v, got: %v", want.Method, got.Method)
	}
	if want.URL.RawQuery != got.URL.RawQuery {
		t.Errorf("wnt: %v, got: %v", want.URL.RawQuery, got.URL.RawQuery)
	}

	wantBody := make([]byte, len(bodyJson))
	gotBody := make([]byte, len(bodyJson))
	if _, err := want.Body.Read(wantBody); err != nil {
		t.Fatal(err)
	}
	if _, err := got.Body.Read(gotBody); err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(string(wantBody), string(gotBody)) {
		t.Errorf("diff: %v", cmp.Diff(wantBody, gotBody))
	}
	if !cmp.Equal(want.Header, got.Header) {
		t.Errorf("diff: %v", cmp.Diff(want.Header, got.Header))
	}
}

func TestPut(t *testing.T) {
	bodyJson, err := json.Marshal(struct{ Text string }{"body"})
	if err != nil {
		t.Fatal(err)
	}

	want := http.Request{
		Method: http.MethodPut,
		URL: &url.URL{
			Scheme:   "https",
			Host:     "api.track.toggl.com",
			User:     url.UserPassword("token", "api_token"),
			RawQuery: "key=value",
		},
		Body: io.ReadCloser(io.NopCloser(bytes.NewBuffer(bodyJson))),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}
	client := timeentries.Client{
		Token: "token",
	}

	// test for PUT method
	got := client.Put(url.URL{RawQuery: "key=value"}, bodyJson)

	// compare
	if want.Method != got.Method {
		t.Errorf("want: %v, got: %v", want.Method, got.Method)
	}
	if want.URL.RawQuery != got.URL.RawQuery {
		t.Errorf("wnt: %v, got: %v", want.URL.RawQuery, got.URL.RawQuery)
	}

	wantBody := make([]byte, len(bodyJson))
	gotBody := make([]byte, len(bodyJson))
	if _, err := want.Body.Read(wantBody); err != nil {
		t.Fatal(err)
	}
	if _, err := got.Body.Read(gotBody); err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(string(wantBody), string(gotBody)) {
		t.Errorf("diff: %v", cmp.Diff(wantBody, gotBody))
	}
	if !cmp.Equal(want.Header, got.Header) {
		t.Errorf("diff: %v", cmp.Diff(want.Header, got.Header))
	}
}

func TestDelete(t *testing.T) {
	want := http.Request{
		Method: http.MethodDelete,
		URL: &url.URL{
			Scheme: "https",
			Host:   "api.track.toggl.com",
			User:   url.UserPassword("token", "api_token"),
		},
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}
	client := timeentries.Client{
		Token: "token",
	}

	// test for Delete method
	got := client.Delete(url.URL{})

	// compare
	if want.Method != got.Method {
		t.Errorf("want: %v, got: %v", want.Method, got.Method)
	}
	if want.URL.RawQuery != got.URL.RawQuery {
		t.Errorf("wnt: %v, got: %v", want.URL.RawQuery, got.URL.RawQuery)
	}
	if !cmp.Equal(want.Header, got.Header) {
		t.Errorf("diff: %v", cmp.Diff(want.Header, got.Header))
	}
}
