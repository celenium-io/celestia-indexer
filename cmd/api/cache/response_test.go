// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
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

func TestWriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	recorder := NewResponseRecorder(w)
	recorder.WriteHeader(http.StatusOK)
	if recorder.status != http.StatusOK {
		t.Errorf("Status code not set correctly")
	}
}

func TestResult(t *testing.T) {
	w := httptest.NewRecorder()
	recorder := NewResponseRecorder(w)
	recorder.WriteHeader(http.StatusOK)
	const data = "test data"
	_, err := recorder.Write([]byte(data))
	if err != nil {
		t.Errorf("Error on write: %v", err)
	}
	result := recorder.Result()
	if result.StatusCode != http.StatusOK || string(result.Body) != data {
		t.Errorf("Incorrect CacheEntry returned")
	}
}

func TestEncodeDecode(t *testing.T) {
	entry := &CacheEntry{
		Header:     make(http.Header),
		StatusCode: http.StatusOK,
		Body:       []byte("test data"),
	}
	data, err := entry.Encode()
	if err != nil {
		t.Errorf("Error encoding CacheEntry")
	}
	newEntry := &CacheEntry{}
	err = newEntry.Decode(data)
	if err != nil || newEntry.StatusCode != http.StatusOK || string(newEntry.Body) != "test data" {
		t.Errorf("Error decoding CacheEntry")
	}
}

func TestReplay(t *testing.T) {
	entry := &CacheEntry{
		Header:     make(http.Header),
		StatusCode: http.StatusOK,
		Body:       []byte("test data"),
	}
	w := httptest.NewRecorder()
	err := entry.Replay(w)
	if err != nil || w.Code != http.StatusOK || w.Body.String() != "test data" {
		t.Errorf("Error replaying CacheEntry")
	}
}
