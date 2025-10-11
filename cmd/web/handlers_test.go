package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zrotrasukha/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	// ts provides a test server
	// this server uses self-signed tls certificates since it is a local test
	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	// The address where this test server is running is stored in ts.URL
	// I am going to use it with http.Client().Get()
	// NOTE: we are not using http.Get(ts.URL + "/ping")
	// because it will rather give us results like :
	// Get "https://127.0.0.1:xxxxx/ping": x509: certificate signed by unknown authority
	// and http.Get() will not trust the test server's self signed certificates
	// NOTE: ts.Client is configured to verify the certificates generaetd
	// by the test server
	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
