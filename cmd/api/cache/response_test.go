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

