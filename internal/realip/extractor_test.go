// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package realip

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func newRequest(remoteAddr, xff string) *http.Request {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	req.RemoteAddr = remoteAddr
	if xff != "" {
		req.Header.Set("X-Forwarded-For", xff)
	}
	return req
}

func TestExtractor(t *testing.T) {
	tests := []struct {
		name           string
		trustedProxies string
		remoteAddr     string
		xff            string
		want           string
	}{
		{
			name:           "empty config ignores forwarding headers",
			trustedProxies: "",
			remoteAddr:     "203.0.113.10:51000",
			xff:            "1.2.3.4",
			want:           "203.0.113.10",
		},
		{
			name:           "direct client cannot spoof via header",
			trustedProxies: "198.51.100.7",
			remoteAddr:     "203.0.113.10:51000",
			xff:            "1.2.3.4",
			want:           "203.0.113.10",
		},
		{
			name:           "trusted proxy header is honored",
			trustedProxies: "198.51.100.7",
			remoteAddr:     "198.51.100.7:43000",
			xff:            "203.0.113.10",
			want:           "203.0.113.10",
		},
		{
			name:           "client-supplied entry behind trusted proxy is skipped",
			trustedProxies: "198.51.100.7",
			remoteAddr:     "198.51.100.7:43000",
			xff:            "1.2.3.4, 203.0.113.10",
			want:           "203.0.113.10",
		},
		{
			name:           "cidr range",
			trustedProxies: "198.51.100.0/24",
			remoteAddr:     "198.51.100.42:43000",
			xff:            "203.0.113.10",
			want:           "203.0.113.10",
		},
		{
			name:           "multiple proxies comma separated",
			trustedProxies: "192.0.2.1, 198.51.100.7",
			remoteAddr:     "192.0.2.1:43000",
			xff:            "203.0.113.10",
			want:           "203.0.113.10",
		},
		{
			name:           "private remote addr is not trusted by default",
			trustedProxies: "198.51.100.7",
			remoteAddr:     "10.0.0.2:43000",
			xff:            "1.2.3.4",
			want:           "10.0.0.2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor, err := Extractor(tt.trustedProxies)
			require.NoError(t, err)
			require.Equal(t, tt.want, extractor(newRequest(tt.remoteAddr, tt.xff)))
		})
	}
}

func TestExtractorInvalid(t *testing.T) {
	for _, value := range []string{"not-an-ip", "10.0.0.0/99", "1.2.3.4;5.6.7.8"} {
		_, err := Extractor(value)
		require.Error(t, err, value)
	}
}
