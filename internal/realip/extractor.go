// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package realip

import (
	"net"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Extractor returns an echo.IPExtractor built from a comma-separated list of
// trusted proxy addresses (single IPs or CIDR ranges). X-Forwarded-For is
// honored only for connections coming from these addresses. When the list is
// empty, the remote address of the connection is always used and forwarding
// headers are ignored.
func Extractor(trustedProxies string) (echo.IPExtractor, error) {
	options := []echo.TrustOption{
		echo.TrustLoopback(false),
		echo.TrustLinkLocal(false),
		echo.TrustPrivateNet(false),
	}

	var trusted int
	for _, proxy := range strings.Split(trustedProxies, ",") {
		proxy = strings.TrimSpace(proxy)
		if proxy == "" {
			continue
		}
		if !strings.Contains(proxy, "/") {
			ip := net.ParseIP(proxy)
			if ip == nil {
				return nil, errors.Errorf("invalid trusted proxy address: %s", proxy)
			}
			bits := 32
			if ip.To4() == nil {
				bits = 128
			}
			options = append(options, echo.TrustIPRange(&net.IPNet{
				IP:   ip,
				Mask: net.CIDRMask(bits, bits),
			}))
			trusted++
			continue
		}
		_, ipNet, err := net.ParseCIDR(proxy)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid trusted proxy range: %s", proxy)
		}
		options = append(options, echo.TrustIPRange(ipNet))
		trusted++
	}

	if trusted == 0 {
		return echo.ExtractIPDirect(), nil
	}
	return echo.ExtractIPFromXFFHeader(options...), nil
}
