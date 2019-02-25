package testttp_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jsteenb2/testttp"
)

func TestHTTP(t *testing.T) {
	svr := newMux()
	t.Run("GET", func(t *testing.T) {
		testttp.GET(
			t, svr, "/",
			testttp.StatusOK(),
			testttp.Resp(assertResp(http.MethodGet)),
		)
	})

	t.Run("POST", func(t *testing.T) {
		testttp.POST(
			t, svr, "/", nil,
			testttp.StatusCreated(),
			testttp.Resp(assertResp(http.MethodPost)),
		)
	})

	t.Run("PUT", func(t *testing.T) {
		testttp.PUT(
			t, svr, "/", nil,
			testttp.StatusAccepted(),
			testttp.Resp(assertResp(http.MethodPut)),
		)
	})

	t.Run("PATCH", func(t *testing.T) {
		testttp.PATCH(
			t, svr, "/", nil,
			testttp.StatusPartialContent(),
			testttp.Resp(assertResp(http.MethodPatch)),
		)
	})

	t.Run("DELETE", func(t *testing.T) {
		testttp.DELETE(t, svr, "/", testttp.StatusNoContent())
	})
}

type foo struct {
	Name, Thing, Method string
}

func newMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			writeFn(w, req.Method, http.StatusOK)
		case http.MethodPost:
			writeFn(w, req.Method, http.StatusCreated)
		case http.MethodPut:
			writeFn(w, req.Method, http.StatusAccepted)
		case http.MethodPatch:
			writeFn(w, req.Method, http.StatusPartialContent)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	return mux
}

func assertResp(method string) func(t *testing.T, w *httptest.ResponseRecorder) {
	return func(t *testing.T, w *httptest.ResponseRecorder) {
		var f foo
		if err := json.NewDecoder(w.Body).Decode(&f); err != nil {
			t.Fatal(err)
		}
		expected := foo{Name: "name", Thing: "thing", Method: method}
		equals(t, expected, f)
	}
}

func writeFn(w http.ResponseWriter, method string, statusCode int) {
	f := foo{Name: "name", Thing: "thing", Method: method}
	r, err := encodeBuf(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	if _, err := io.Copy(w, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func equals(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected == actual {
		return
	}
	t.Errorf("expected: %v\tactual: %v", expected, actual)
}

func encodeBuf(v interface{}) (io.Reader, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return &buf, nil
}
