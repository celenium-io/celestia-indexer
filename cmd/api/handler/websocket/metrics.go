// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

package websocket

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// WebSocket connection metrics
	wsActiveConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "websocket_active_connections",
		Help: "Current number of active WebSocket connections",
	})

	wsConnectionsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "websocket_connections_total",
		Help: "Total number of WebSocket connections",
	}, []string{"status"}) // status: accepted, rejected

	// Message metrics
	wsMessagesSent = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "websocket_messages_sent_total",
		Help: "Total number of messages sent to clients",
	}, []string{"channel"}) // channel: head, blocks, gas_price

	wsMessagesDropped = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "websocket_messages_dropped_total",
		Help: "Total number of messages dropped due to full client buffer",
	}, []string{"channel"})

	wsMessageLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "websocket_message_broadcast_seconds",
		Help:    "Time taken to broadcast message to all clients",
		Buckets: prometheus.DefBuckets,
	}, []string{"channel"})

	// Subscription metrics
	wsSubscriptions = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "websocket_subscriptions",
		Help: "Current number of active subscriptions per channel",
	}, []string{"channel"})

	wsSubscribeRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "websocket_subscribe_requests_total",
		Help: "Total number of subscribe requests",
	}, []string{"channel", "status"}) // status: success, error

	wsUnsubscribeRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "websocket_unsubscribe_requests_total",
		Help: "Total number of unsubscribe requests",
	}, []string{"channel"})

	// Error metrics
	wsErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "websocket_errors_total",
		Help: "Total number of WebSocket errors",
	}, []string{"type"}) // type: read, write, upgrade, unknown_method, unknown_channel

	// Connection duration
	wsConnectionDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "websocket_connection_duration_seconds",
		Help:    "Duration of WebSocket connections",
		Buckets: []float64{1, 10, 60, 300, 600, 1800, 3600, 7200, 14400},
	})
)
