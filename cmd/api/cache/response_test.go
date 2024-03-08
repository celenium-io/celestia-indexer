package cache

import (
	_ "bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewResponseRecorder(t *testing.T) {
	w := httptest.NewRecorder()
	recorder := NewResponseRecorder(w)
	if recorder.ResponseWriter != w {
		t.Errorf("ResponseWriter not set correctly")
	}
}

func TestWrite(t *testing.T) {
	w := httptest.NewRecorder()
	recorder := NewResponseRecorder(w)
	data := []byte("test data")
	n, err := recorder.Write(data)
	if err != nil || n != len(data) {
		t.Errorf("Error writing data")
	}
}

func TestCopyHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	recorder := NewResponseRecorder(w)
	w.Header().Set("Test-Header", "Test-Value")
	recorder.copyHeaders()
	if recorder.headers.Get("Test-Header") != "Test-Value" {
		t.Errorf("Headers not copied correctly")
	}
}



