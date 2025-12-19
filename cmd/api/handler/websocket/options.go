package websocket

type ManagerOption func(*Manager)

func WithWebsocketClientsPerIp(limit int) ManagerOption {
	return func(m *Manager) {
		m.websocketClientsPerIp = limit
	}
}
