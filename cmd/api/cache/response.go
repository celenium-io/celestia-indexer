// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"strings"
)

type ResponseRecorder struct {
	http.ResponseWriter

	status       int
	body         bytes.Buffer
	headers      http.Header
	headerCopied bool
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		headers:        make(http.Header),
	}
}

func (w *ResponseRecorder) Write(b []byte) (int, error) {
	w.copyHeaders()
	i, err := w.ResponseWriter.Write(b)
	if err != nil {
		return i, err
	}

	return w.body.Write(b[:i])
}

func (r *ResponseRecorder) copyHeaders() {
	if r.headerCopied {
		return
	}

	r.headerCopied = true
	copyHeaders(r.ResponseWriter.Header(), r.headers)
}

func (w *ResponseRecorder) WriteHeader(statusCode int) {
	w.copyHeaders()

	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Result() *CacheEntry {
	r.copyHeaders()

	return &CacheEntry{
		Header:     r.headers,
		StatusCode: r.status,
		Body:       r.body.Bytes(),
	}
}

func copyHeaders(src, dst http.Header) {
	for k, v := range src {
		for _, val := range v {
			dst.Set(k, val)
		}
	}
}

type CacheEntry struct {
	Header     http.Header
	StatusCode int
	Body       []byte
}

func (c *CacheEntry) Encode() (string, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(c); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (c *CacheEntry) Decode(b string) error {
	dec := gob.NewDecoder(strings.NewReader(b))
	return dec.Decode(c)
}

func (c *CacheEntry) Replay(w http.ResponseWriter) error {
	copyHeaders(c.Header, w.Header())
	if c.StatusCode != 0 {
		w.WriteHeader(c.StatusCode)
	}

	if len(c.Body) == 0 {
		return nil
	}

	_, err := w.Write(c.Body)
	return err
}
