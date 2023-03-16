package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	API_URL1 = "http://localhost:1000/article"
	API_URL2 = "http://localhost:1000/articles"
)

func testApiCall(b *testing.B, small bool, debug bool) {
	url := API_URL1
	if !small {
		url = API_URL2
	}
	var (
		body *bytes.Buffer
	)
	for i := 0; i < b.N; i++ {
		body = bytes.NewBufferString(``)
		req, err := http.NewRequest("GET", url, body)
		require.NoError(b, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(b, err)
		// log.Println(resp)
		// if debug {
		txt, err := io.ReadAll(resp.Body)
		require.NoError(b, err)
		log.Println(string(txt))
		// 	require.NoError(b, err)
		// }
		resp.Body.Close()
	}
}

// func setupHttpMock(b *testing.B, tiny bool) func() {
// 	httpmock.Activate()
// 	if tiny {
// 		httpmock.RegisterResponder("GET", API_URL,
// 			httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Article"}]`))
// 	} else {
// 		var msg string
// 		msg += `[`
// 		for i := 0; i < 100; i++ {
// 			if i > 0 {
// 				msg += `,`
// 			}
// 			msg += fmt.Sprintf(`{"id": %d, "name": "My Great Article"}`, i+1)
// 		}
// 		msg += `]`
// 		httpmock.RegisterResponder("GET", API_URL,
// 			httpmock.NewStringResponder(200, msg))
// 	}

// 	return httpmock.DeactivateAndReset
// }

func BenchmarkProcessNothing(b *testing.B) {
	// reset := setupHttpMock(b, true)
	// defer reset()
	testApiCall(b, true, false)
}

func BenchmarkProcessNothingBig(b *testing.B) {
	// reset := setupHttpMock(b, true)
	// defer reset()
	testApiCall(b, false, false)
}

func BenchmarkProcessWithLog(b *testing.B) {
	// reset := setupHttpMock(b, true)
	// defer reset()
	testApiCall(b, true, true)
}

func BenchmarkProcessWithLogBigMsg(b *testing.B) {
	// reset := setupHttpMock(b, false)
	// defer reset()
	testApiCall(b, false, true)
}

func testTracing(b *testing.B, small bool) {
	url := API_URL1
	if !small {
		url = API_URL2
	}

	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		panic(err)
	}

	otel.SetTracerProvider(tp)

	// TODO: setup trace provider
	var (
		body *bytes.Buffer
		ctx  context.Context
		span trace.Span
	)
	for i := 0; i < b.N; i++ {
		ctx = context.Background()
		body = bytes.NewBufferString(``)
		req, err := http.NewRequest("GET", url, body)
		require.NoError(b, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(b, err)
		tr := tp.Tracer("httpRequest")
		_, span = tr.Start(ctx, url)
		txt, err := io.ReadAll(resp.Body)
		require.NoError(b, err)
		err = resp.Body.Close()
		require.NoError(b, err)
		span.SetAttributes(attribute.Key("testset").String(string(txt)))
		span.End()
	}
}

func BenchmarkProcessWithTrace(b *testing.B) {
	// reset := setupHttpMock(b, true)
	// defer reset()
	testTracing(b, true)
}

func BenchmarkProcessWithTraceBigMsg(b *testing.B) {
	// reset := setupHttpMock(b, false)
	// defer reset()
	testTracing(b, false)
}
