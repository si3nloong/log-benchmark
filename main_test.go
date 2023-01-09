package logbenchmark

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func apiCall(b *testing.B, debug bool) {
	body := bytes.NewBufferString(``)
	req, err := http.NewRequest("GET", "https://api.mybiz.com/articles", body)
	require.NoError(b, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(b, err)
	if debug {
		txt, err := io.ReadAll(resp.Body)
		require.NoError(b, err)
		b.Log(string(txt))
	}
}

func setupHttpMock(b *testing.B) func() {
	httpmock.Activate()
	httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles",
		httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Article"}]`))
	return httpmock.DeactivateAndReset
}

func BenchmarkProcessWithLog(b *testing.B) {
	reset := setupHttpMock(b)
	defer reset()
	for i := 0; i < b.N; i++ {
		apiCall(b, true)
	}
}

func BenchmarkProcessWithTrace(b *testing.B) {
	reset := setupHttpMock(b)
	defer reset()
	http.DefaultClient = &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	for i := 0; i < b.N; i++ {
		ctx := context.TODO()
		body := bytes.NewBufferString(``)
		req, err := http.NewRequest("GET", "https://api.mybiz.com/articles", body)
		require.NoError(b, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(b, err)
		tr := otel.Tracer("component-bar")
		_, span := tr.Start(ctx, "bar")
		txt, err := io.ReadAll(resp.Body)
		require.NoError(b, err)
		resp.Body.Close()
		span.SetAttributes(attribute.Key("testset").String(string(txt)))
		span.End()
	}
}

func BenchmarkProcessWithoutLogAndTrace(b *testing.B) {
	reset := setupHttpMock(b)
	defer reset()
	for i := 0; i < b.N; i++ {
		apiCall(b, false)
	}
}
