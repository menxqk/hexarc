package frontend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

const apiUrl = "/api/"
const allUrl = "/all"

func TestWebserverNotAllowed(t *testing.T) {
	r := httptest.NewRequest(http.MethodOptions, apiUrl, nil)
	w := httptest.NewRecorder()
	testWebserverFrontend.notAllowedHandler(w, r)
	if want, got := http.StatusMethodNotAllowed, w.Result().StatusCode; want != got {
		t.Errorf("expected a %d, instead got: %d", want, got)
	}
}

func TestWebserverPutGetDelete(t *testing.T) {
	const key = "a-key"
	const value = "a-value"

	// Set mux routes for testing, otherwise mux.Vars(r) will
	// return empty map
	m := mux.NewRouter()
	m.HandleFunc(apiUrl+"{key}", testWebserverFrontend.putHandler).Methods("PUT")
	m.HandleFunc(apiUrl+"{key}", testWebserverFrontend.getHandler).Methods("GET")
	m.HandleFunc(apiUrl+"{key}", testWebserverFrontend.deleteHandler).Methods("DELETE")

	// Put key and value
	body := strings.NewReader(value)
	r := httptest.NewRequest(http.MethodPut, apiUrl+key, body)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r) // Use mux for testing, do not use handler directly!
	if want, got := http.StatusCreated, w.Result().StatusCode; want != got {
		t.Errorf("put - expected a %d, instead got: %d", want, got)
	}

	// Retrieve value and compare
	r = httptest.NewRequest(http.MethodGet, apiUrl+key, nil)
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Errorf("get - expected a %d, instead got: %d", want, got)
	}
	res := w.Result()
	var resp jsonResponse
	err := json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
	res.Body.Close()

	val, ok := resp.Data.(string)
	if !ok {
		t.Errorf("wrong type for resp.Data: %T", resp.Data)
	}

	if string(val) != value {
		t.Errorf("val/value mismatch, val: %q, value: %q", val, value)
	}

	// Delete key and value
	r = httptest.NewRequest(http.MethodDelete, apiUrl+key, nil)
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Errorf("delete - expected a %d, instead got: %d", want, got)
	}

	// Try to get deleted key and value
	r = httptest.NewRequest(http.MethodGet, apiUrl+key, nil)
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	if want, got := http.StatusNotFound, w.Result().StatusCode; want != got {
		t.Errorf("get after delete - expected a %d, instead got: %d", want, got)
	}
}

func TestWebserverGetAll(t *testing.T) {
	// Set mux routes for testing, otherwise mux.Vars(r) will
	// return empty map
	m := mux.NewRouter()
	m.HandleFunc(allUrl, testWebserverFrontend.getAllHandler).Methods("GET")

	// Get all
	r := httptest.NewRequest(http.MethodGet, allUrl, nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r) // Use mux for testing, do not use handler directly!
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Errorf("getall - expected a %d, instead got: %d", want, got)
	}
	res := w.Result()
	defer res.Body.Close()
	var resp jsonResponse
	err := json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
	}
}
