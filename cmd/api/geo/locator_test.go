// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package geo

import (
	"context"
	"net"
	"net/netip"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type hostsSource []string

func (h hostsSource) FibreHosts(context.Context) ([]string, error) { return h, nil }

func writeDB(t *testing.T, dir, name, dbType string, records map[string]mmdbtype.Map) {
	t.Helper()
	tree, err := mmdbwriter.New(mmdbwriter.Options{DatabaseType: dbType, RecordSize: 24})
	require.NoError(t, err)
	for cidr, value := range records {
		_, network, err := net.ParseCIDR(cidr)
		require.NoError(t, err)
		require.NoError(t, tree.Insert(network, value))
	}
	f, err := os.Create(filepath.Join(dir, name))
	require.NoError(t, err)
	defer f.Close()
	_, err = tree.WriteTo(f)
	require.NoError(t, err)
}

func newTestLocator(t *testing.T, hosts HostsSource, lookup lookupFunc) *Locator {
	t.Helper()
	dir := t.TempDir()
	writeDB(t, dir, CountryFile, "country", map[string]mmdbtype.Map{
		"195.154.212.0/24": {"country_code": mmdbtype.String("FR")},
		"65.21.0.0/16":     {"country_code": mmdbtype.String("FI")},
		"2a01:4f9::/32":    {"country_code": mmdbtype.String("FI")},
	})
	writeDB(t, dir, AsnFile, "asn", map[string]mmdbtype.Map{
		"195.154.212.0/24": {
			"autonomous_system_number":       mmdbtype.Uint32(12876),
			"autonomous_system_organization": mmdbtype.String("Scaleway SAS"),
		},
		"65.21.0.0/16": {
			"autonomous_system_number":       mmdbtype.Uint32(24940),
			"autonomous_system_organization": mmdbtype.String("Hetzner Online GmbH"),
		},
	})

	l, err := NewLocator(dir, hosts)
	require.NoError(t, err)
	if lookup != nil {
		l.lookup = lookup
	}
	t.Cleanup(func() { _ = l.Close() })
	return l
}

func TestHostname(t *testing.T) {
	tests := []struct {
		host string
		want string
	}{
		{"195.154.212.12:7980", "195.154.212.12"},
		{"fibre.mocha.ppbcn.com:7980", "fibre.mocha.ppbcn.com"},
		{"fibre.mocha.ppbcn.com", "fibre.mocha.ppbcn.com"},
		{"195.154.212.12", "195.154.212.12"},
		{"[2a01:4f9::1]:7980", "2a01:4f9::1"},
		{"[2a01:4f9::1]", "2a01:4f9::1"},
		{"dns:///fibre.example.com:443", "fibre.example.com"},
		{"https://fibre.example.com:443", "fibre.example.com"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			require.Equal(t, tt.want, hostname(tt.host))
		})
	}
}

func TestPickAddr(t *testing.T) {
	parse := func(addrs ...string) []netip.Addr {
		res := make([]netip.Addr, len(addrs))
		for i := range addrs {
			res[i] = netip.MustParseAddr(addrs[i])
		}
		return res
	}
	tests := []struct {
		name  string
		addrs []netip.Addr
		want  string
	}{
		{"single v4", parse("65.21.15.43"), "65.21.15.43"},
		{"single v6", parse("2a01:4f9::1"), "2a01:4f9::1"},
		{"v6 first", parse("2a01:4f9::1", "65.21.15.43"), "65.21.15.43"},
		{"v4 first", parse("65.21.15.43", "2a01:4f9::1"), "65.21.15.43"},
		{"lowest v4", parse("195.154.212.12", "65.21.15.43", "2a01:4f9::1"), "65.21.15.43"},
		{"lowest v6", parse("2a01:4f9::2", "2a01:4f9::1"), "2a01:4f9::1"},
		{"v4-mapped", parse("2a01:4f9::1", "::ffff:65.21.15.43"), "65.21.15.43"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, pickAddr(tt.addrs).String())
		})
	}
}

func TestLocateIP(t *testing.T) {
	l := newTestLocator(t, hostsSource{}, nil)

	require.Equal(t, &Location{
		IP:          "195.154.212.12",
		CountryCode: "FR",
		ASN:         12876,
		ASOrg:       "Scaleway SAS",
	}, l.Locate("195.154.212.12:7980"))

	require.Equal(t, &Location{
		IP:          "2a01:4f9::1",
		CountryCode: "FI",
	}, l.Locate("[2a01:4f9::1]:7980"))

	require.Nil(t, l.Locate("10.0.0.1:7980"))
	require.Nil(t, l.Locate("garbage"))
}

func TestLocateDNS(t *testing.T) {
	const host = "fibre.example.com:7980"
	var fail bool
	lookup := func(_ context.Context, _, name string) ([]netip.Addr, error) {
		require.Equal(t, "fibre.example.com", name)
		if fail {
			return nil, errors.New("no such host")
		}
		return []netip.Addr{netip.MustParseAddr("65.21.15.43")}, nil
	}
	l := newTestLocator(t, hostsSource{host, "195.154.212.12:7980"}, lookup)

	require.Nil(t, l.Locate(host), "DNS names are not resolved on request")

	l.refresh(t.Context())
	want := &Location{
		IP:          "65.21.15.43",
		CountryCode: "FI",
		ASN:         24940,
		ASOrg:       "Hetzner Online GmbH",
	}
	require.Equal(t, want, l.Locate(host))
	require.Equal(t, "FR", l.Locate("195.154.212.12:7980").CountryCode)

	fail = true
	l.refresh(t.Context())
	require.Equal(t, want, l.Locate(host), "keeps the last known location")
}

func TestLocateResolvesOnMiss(t *testing.T) {
	const host = "new.example.com:7980"
	var calls atomic.Int32
	lookup := func(_ context.Context, _, name string) ([]netip.Addr, error) {
		calls.Add(1)
		return []netip.Addr{netip.MustParseAddr("65.21.15.43")}, nil
	}
	// the host is not in the hourly snapshot yet, as after a fresh MsgSetFibreProviderInfo
	l := newTestLocator(t, hostsSource{}, lookup)
	l.Start(t.Context())

	require.Nil(t, l.Locate(host))
	require.Eventually(t, func() bool {
		return l.Locate(host) != nil
	}, time.Second, 10*time.Millisecond)
	require.Equal(t, "FI", l.Locate(host).CountryCode)
	require.EqualValues(t, 1, calls.Load())
}

func TestLocateMissRetryInterval(t *testing.T) {
	const host = "broken.example.com:7980"
	var calls atomic.Int32
	lookup := func(context.Context, string, string) ([]netip.Addr, error) {
		calls.Add(1)
		return nil, errors.New("no such host")
	}
	l := newTestLocator(t, hostsSource{}, lookup)
	l.Start(t.Context())

	for range 10 {
		require.Nil(t, l.Locate(host))
	}
	require.Eventually(t, func() bool { return calls.Load() == 1 }, time.Second, 10*time.Millisecond)
	for range 10 {
		require.Nil(t, l.Locate(host))
	}
	time.Sleep(50 * time.Millisecond)
	require.EqualValues(t, 1, calls.Load(), "failed host is not retried before retryInterval")
}

func TestNewLocatorMissingFiles(t *testing.T) {
	_, err := NewLocator(t.TempDir(), hostsSource{})
	require.Error(t, err)
}
