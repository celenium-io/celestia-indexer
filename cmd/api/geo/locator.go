// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package geo

import (
	"context"
	"io"
	"net"
	"net/netip"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dipdup-io/workerpool"
	"github.com/oschwald/maxminddb-golang/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Files of sapics/ip-location-db (PDDL, no attribution required).
const (
	CountryFile = "server-country.mmdb"
	AsnFile     = "origin-asn.mmdb"
)

const (
	refreshInterval = time.Hour
	resolveTimeout  = 3 * time.Second
)

type Location struct {
	IP          string
	CountryCode string
	ASN         uint32
	ASOrg       string
}

type ILocator interface {
	io.Closer

	Start(ctx context.Context)
	Locate(host string) *Location
}

type HostsSource interface {
	FibreHosts(ctx context.Context) ([]string, error)
}

type lookupFunc func(ctx context.Context, network, host string) ([]netip.Addr, error)

// Locator resolves fibre hosts in the background; Locate never touches DNS.
type Locator struct {
	country *maxminddb.Reader
	asn     *maxminddb.Reader
	hosts   HostsSource
	lookup  lookupFunc

	data map[string]*Location
	mx   *sync.RWMutex
	g    workerpool.Group
	log  zerolog.Logger
}

func NewLocator(dir string, hosts HostsSource) (*Locator, error) {
	country, err := maxminddb.Open(filepath.Join(dir, CountryFile))
	if err != nil {
		return nil, errors.Wrap(err, "open country database")
	}
	asn, err := maxminddb.Open(filepath.Join(dir, AsnFile))
	if err != nil {
		_ = country.Close()
		return nil, errors.Wrap(err, "open asn database")
	}
	return newLocator(country, asn, hosts, net.DefaultResolver.LookupNetIP), nil
}

func newLocator(country, asn *maxminddb.Reader, hosts HostsSource, lookup lookupFunc) *Locator {
	return &Locator{
		country: country,
		asn:     asn,
		hosts:   hosts,
		lookup:  lookup,
		data:    make(map[string]*Location),
		mx:      new(sync.RWMutex),
		g:       workerpool.NewGroup(),
		log:     log.With().Str("module", "geo_locator").Logger(),
	}
}

func (l *Locator) Start(ctx context.Context) {
	l.g.GoCtx(ctx, l.sync)
}

// Locate returns nil for unknown hosts: IP literals are looked up at once, DNS names wait for the next refresh.
func (l *Locator) Locate(host string) *Location {
	l.mx.RLock()
	loc, ok := l.data[host]
	l.mx.RUnlock()
	if ok {
		return loc
	}

	if addr, err := netip.ParseAddr(hostname(host)); err == nil {
		return l.locateAddr(addr)
	}
	return nil
}

func (l *Locator) sync(ctx context.Context) {
	l.refresh(ctx)

	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.refresh(ctx)
		}
	}
}

func (l *Locator) refresh(ctx context.Context) {
	hosts, err := l.hosts.FibreHosts(ctx)
	if err != nil {
		l.log.Error().Err(err).Msg("receiving fibre hosts")
		return
	}

	l.mx.RLock()
	prev := l.data
	l.mx.RUnlock()

	data := make(map[string]*Location, len(hosts))
	for _, host := range hosts {
		addr, err := l.resolve(ctx, host)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			l.log.Warn().Err(err).Str("host", host).Msg("resolving fibre host")
			if loc, ok := prev[host]; ok {
				data[host] = loc
			}
			continue
		}
		data[host] = l.locateAddr(addr)
	}

	l.mx.Lock()
	l.data = data
	l.mx.Unlock()

	l.log.Info().Int("hosts", len(data)).Msg("fibre hosts located")
}

func (l *Locator) resolve(ctx context.Context, host string) (netip.Addr, error) {
	name := hostname(host)
	if addr, err := netip.ParseAddr(name); err == nil {
		return addr, nil
	}

	resolveCtx, cancel := context.WithTimeout(ctx, resolveTimeout)
	defer cancel()

	addrs, err := l.lookup(resolveCtx, "ip", name)
	if err != nil {
		return netip.Addr{}, err
	}
	if len(addrs) == 0 {
		return netip.Addr{}, errors.Errorf("no addresses for %s", name)
	}
	return addrs[0], nil
}

func (l *Locator) locateAddr(addr netip.Addr) *Location {
	addr = addr.Unmap()

	var country struct {
		CountryCode string `maxminddb:"country_code"`
	}
	if err := l.country.Lookup(addr).Decode(&country); err != nil {
		l.log.Error().Err(err).Stringer("ip", addr).Msg("country lookup")
	}

	var asn struct {
		Number       uint32 `maxminddb:"autonomous_system_number"`
		Organization string `maxminddb:"autonomous_system_organization"`
	}
	if err := l.asn.Lookup(addr).Decode(&asn); err != nil {
		l.log.Error().Err(err).Stringer("ip", addr).Msg("asn lookup")
	}

	if country.CountryCode == "" && asn.Number == 0 {
		return nil
	}
	return &Location{
		IP:          addr.String(),
		CountryCode: country.CountryCode,
		ASN:         asn.Number,
		ASOrg:       asn.Organization,
	}
}

func (l *Locator) Close() error {
	l.g.Wait()
	countryErr := l.country.Close()
	if err := l.asn.Close(); err != nil {
		return err
	}
	return countryErr
}

// hostname strips an optional scheme, port and IPv6 brackets.
func hostname(host string) string {
	if _, rest, ok := strings.Cut(host, "://"); ok {
		host = strings.TrimLeft(rest, "/")
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		return h
	}
	return strings.Trim(host, "[]")
}

// Noop is used when the databases are not available.
type Noop struct{}

func (Noop) Start(context.Context)   {}
func (Noop) Locate(string) *Location { return nil }
func (Noop) Close() error            { return nil }
