// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package websocket

type ManagerOption func(*Manager)

func WithWebsocketClientsPerIp(limit int) ManagerOption {
	return func(m *Manager) {
		m.websocketClientsPerIp = limit
	}
}
