package logbenchmark

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

const API_URL = "https://api.mybiz.com/articles"

func apiCall(b *testing.B, debug bool) {
	body := bytes.NewBufferString(``)
	req, err := http.NewRequest("GET", API_URL, body)
	require.NoError(b, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(b, err)
	if debug {
		txt, err := io.ReadAll(resp.Body)
		require.NoError(b, err)
		b.Log(string(txt))
		err = resp.Body.Close()
		require.NoError(b, err)
	}
}

func setupHttpMock(b *testing.B, tiny bool) func() {
	httpmock.Activate()
	if tiny {
		httpmock.RegisterResponder("GET", API_URL,
			httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Article"}]`))
	} else {
		var msg string
		msg += `[`
		for i := 0; i < 100; i++ {
			if i > 0 {
				msg += `,`
			}
			msg += fmt.Sprintf(`{"id": %d, "name": "My Great Article"}`, i+1)
		}
		msg += `]`
		httpmock.RegisterResponder("GET", API_URL,
			httpmock.NewStringResponder(200, msg))
	}

	return httpmock.DeactivateAndReset
}

func BenchmarkProcessWithLog(b *testing.B) {
	reset := setupHttpMock(b, true)
	defer reset()
	for i := 0; i < b.N; i++ {
		apiCall(b, true)
	}
}

func BenchmarkProcessWithLogBigMsg(b *testing.B) {
	reset := setupHttpMock(b, false)
	defer reset()
	for i := 0; i < b.N; i++ {
		apiCall(b, true)
	}
}

func testTracing(b *testing.B) {
	http.DefaultClient = &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	// log.Println(otel.GetTracerProvider())
	// TODO: setup trace provider
	for i := 0; i < b.N; i++ {
		ctx := context.TODO()
		body := bytes.NewBufferString(``)
		req, err := http.NewRequest("GET", API_URL, body)
		require.NoError(b, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(b, err)
		tr := otel.Tracer("request")
		_, span := tr.Start(ctx, "bar")
		txt, err := io.ReadAll(resp.Body)
		require.NoError(b, err)
		err = resp.Body.Close()
		require.NoError(b, err)
		span.SetAttributes(attribute.Key("testset").String(string(txt)))
		span.End()
	}
}

func BenchmarkProcessWithTrace(b *testing.B) {
	reset := setupHttpMock(b, true)
	defer reset()
	testTracing(b)
}

func BenchmarkProcessWithTraceBigMsg(b *testing.B) {
	reset := setupHttpMock(b, false)
	defer reset()
	testTracing(b)
}

func BenchmarkProcessWithoutLogAndTrace(b *testing.B) {
	reset := setupHttpMock(b, true)
	defer reset()
	for i := 0; i < b.N; i++ {
		apiCall(b, false)
	}
}

func BenchmarkProcessWithoutLogAndTraceBigMsg(b *testing.B) {
	reset := setupHttpMock(b, false)
	defer reset()
	for i := 0; i < b.N; i++ {
		apiCall(b, false)
	}
}
