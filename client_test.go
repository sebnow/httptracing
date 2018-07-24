package httptracing

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/opentracing/opentracing-go/mocktracer"
)

// Client is an interface for the http.Client implementation
type Client interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (resp *http.Response, err error)
	Head(url string) (resp *http.Response, err error)
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

func TestTracingClientImplementInterface(t *testing.T) {
	var client interface{} = &TracingClient{}
	if _, ok := client.(Client); !ok {
		t.Errorf("TracingClient does not implement the Client interface")
	}
}

func TestTracingClientDoInjectsCarrierHeader(t *testing.T) {
	server := httptest.NewServer(assertRequest(t, "GET"))
	defer server.Close()

	tracer := mocktracer.New()
	client := Trace(tracer, http.DefaultClient)
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		panic(err)
	}

	_, err = client.Do(req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestTracingClientGetInjectsCarrierHeader(t *testing.T) {
	server := httptest.NewServer(assertRequest(t, "GET"))
	defer server.Close()

	tracer := mocktracer.New()
	client := Trace(tracer, http.DefaultClient)

	_, err := client.Get(server.URL)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestTracingClientGetErrorsOnInvalidURL(t *testing.T) {
	tracer := mocktracer.New()
	client := Trace(tracer, http.DefaultClient)

	_, err := client.Get("://foo.com")
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestTracingClientPostInjectsCarrierHeader(t *testing.T) {
	server := httptest.NewServer(assertRequest(t, "POST"))
	defer server.Close()

	tracer := mocktracer.New()
	client := Trace(tracer, http.DefaultClient)

	_, err := client.Post(server.URL, "application/json", nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestTracingClientPostErrorsOnInvalidURL(t *testing.T) {
	tracer := mocktracer.New()
	client := Trace(tracer, http.DefaultClient)

	_, err := client.Post("://foo.com", "application/json", nil)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestTracingClientHeadInjectsCarrierHeader(t *testing.T) {
	server := httptest.NewServer(assertRequest(t, "HEAD"))
	defer server.Close()

	tracer := mocktracer.New()
	client := Trace(tracer, http.DefaultClient)

	_, err := client.Head(server.URL)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestTracingClientHeadErrorsOnInvalidURL(t *testing.T) {
	tracer := mocktracer.New()
	client := Trace(tracer, http.DefaultClient)

	_, err := client.Head("://foo.com")
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestTracingClientPostFormInjectsCarrierHeader(t *testing.T) {
	server := httptest.NewServer(assertRequest(t, "POST"))
	defer server.Close()

	tracer := mocktracer.New()
	client := Trace(tracer, http.DefaultClient)

	_, err := client.PostForm(server.URL, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestTracingClientPostFormErrorsOnInvalidURL(t *testing.T) {
	tracer := mocktracer.New()
	client := Trace(tracer, http.DefaultClient)

	_, err := client.PostForm("://foo.com", nil)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func assertRequest(t *testing.T, method string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("Expected method %s, got %s", method, r.Method)
		}

		if r.Header.Get("Mockpfx-Ids-Spanid") == "" {
			t.Errorf("Expected Mockpfx-Ids-Spanid header in the request")
		}

		if r.Header.Get("Mockpfx-Ids-Traceid") == "" {
			t.Errorf("Expected Mockpfx-Ids-Traceid header in the request")
		}

		w.WriteHeader(200)
	})
}
