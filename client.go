package httptracing

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
)

// TracingClient wraps a http.Client with an OpenTracing tracer
type TracingClient struct {
	*http.Client
	tracer opentracing.Tracer
}

// Trace wraps the client with a tracer so each request has
// its own span.
func Trace(tracer opentracing.Tracer, client *http.Client) *TracingClient {
	wrappedClient := *client
	wrappedClient.Transport = &nethttp.Transport{
		RoundTripper: client.Transport,
	}

	return &TracingClient{
		Client: &wrappedClient,
		tracer: tracer,
	}
}

// Do is the same as http.Client.Do but the request is traced
func (client *TracingClient) Do(req *http.Request) (resp *http.Response, err error) {
	req, ht := nethttp.TraceRequest(client.tracer, req)
	defer ht.Finish()

	return client.Client.Do(req)
}

// Get is the same as http.Client.Get but the request is traced
func (client *TracingClient) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}

// Post is the same as http.Client.Post but the request is traced
func (client *TracingClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", contentType)

	return client.Do(req)
}

// Head is the same as http.Client.Head but the request is traced
func (client *TracingClient) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}

// PostForm is the same as http.Client.PostForm but the request is traced
func (client *TracingClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return client.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
