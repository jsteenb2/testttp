package testttp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// GET is a test utility to test a svr handles a GET call.
func GET(t *testing.T, svr http.Handler, addr string, assertFns ...RecTestFn) {
	t.Helper()
	HTTP(t, svr, http.MethodGet, addr, nil, assertFns...)
}

// POST is a test utility to test a svr handles a POST call.
func POST(t *testing.T, svr http.Handler, addr string, body io.Reader, assertFns ...RecTestFn) {
	t.Helper()
	HTTP(t, svr, http.MethodPost, addr, body, assertFns...)
}

// PATCH is a test utility to test a svr handles a PATCH call.
func PATCH(t *testing.T, svr http.Handler, addr string, body io.Reader, assertFns ...RecTestFn) {
	t.Helper()
	HTTP(t, svr, http.MethodPatch, addr, body, assertFns...)
}

// PUT is a test utility to test a svr handles a PUT call.
func PUT(t *testing.T, svr http.Handler, addr string, body io.Reader, assertFns ...RecTestFn) {
	t.Helper()
	HTTP(t, svr, http.MethodPut, addr, body, assertFns...)
}

// DELETE is a test utility to test a svr handles a DELETE call.
func DELETE(t *testing.T, svr http.Handler, addr string, assertFns ...RecTestFn) {
	t.Helper()
	HTTP(t, svr, http.MethodDelete, addr, nil, assertFns...)
}

// HTTP is a test utility to test a svr handles whatever call you provide it.
func HTTP(
	t *testing.T, svr http.Handler,
	method, addr string, body io.Reader,
	assertFns ...RecTestFn,
) {
	t.Helper()
	req, w := httptest.NewRequest(method, addr, body), httptest.NewRecorder()

	svr.ServeHTTP(w, req)

	for _, fn := range assertFns {
		fn(t, w)
	}
}

// RecTestFn is a functional option to run assertions against the response of the http request
// being made by the HTTP func or its relatives.
type RecTestFn func(t *testing.T, w *httptest.ResponseRecorder)

// Resp allows the body to be asserted and viewed per the function.
func Resp(assertFn func(t *testing.T, w *httptest.ResponseRecorder)) RecTestFn {
	return func(t *testing.T, w *httptest.ResponseRecorder) {
		t.Helper()
		assertFn(t, w)
	}
}

// Status verifies the status code of the response matches the status provided.
func Status(status int) RecTestFn {
	return func(t *testing.T, w *httptest.ResponseRecorder) {
		t.Helper()
		if w.Code == status {
			return
		}

		t.Fatalf("received incorrect status code:\twant: %d\tgot: %d", status, w.Code)
	}
}

// StatusOK verifies the status code is 200 (Status OK).
func StatusOK() RecTestFn {
	return Status(http.StatusOK)
}

// StatusCreated verifies the status code is 201 (Status Created).
func StatusCreated() RecTestFn {
	return Status(http.StatusCreated)
}

// StatusAccepted verifies the status code is 202 (Status Accepted).
func StatusAccepted() RecTestFn {
	return Status(http.StatusAccepted)
}

// StatusNoContent verifies the status code is 204 (Status No Content).
func StatusNoContent() RecTestFn {
	return Status(http.StatusNoContent)
}

// StatusPartialContent verifies the status code is 206 (Status Partial Content).
func StatusPartialContent() RecTestFn {
	return Status(http.StatusPartialContent)
}

// StatusNotFound verifies the status code is 404 (Status Not Found).
func StatusNotFound() RecTestFn {
	return Status(http.StatusNotFound)
}

// StatusUnprocessableEntity verifies the status code is 422 (Status Unprocessable Entity).
func StatusUnprocessableEntity() RecTestFn {
	return Status(http.StatusUnprocessableEntity)
}

// StatusInternalServerError verifies the status code is 500 (Status Internal Server Error).
func StatusInternalServerError() RecTestFn {
	return Status(http.StatusInternalServerError)
}

// StatusNotImplemented verifies the status code is 501 (Status Not Implemented).
func StatusNotImplemented() RecTestFn {
	return Status(http.StatusNotImplemented)
}
