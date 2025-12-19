# WebSocket Prometheus Metrics

This document describes the Prometheus metrics exposed by the WebSocket implementation.

## Available Metrics

### Connection Metrics

#### `websocket_active_connections`
- **Type**: Gauge
- **Description**: Current number of active WebSocket connections
- **Use**: Monitor real-time connection count

#### `websocket_connections_total`
- **Type**: Counter
- **Labels**: `status` (accepted, rejected)
- **Description**: Total number of WebSocket connection attempts
- **Use**: Track connection acceptance/rejection rate

#### `websocket_connection_duration_seconds`
- **Type**: Histogram
- **Description**: Duration of WebSocket connections
- **Buckets**: 1s, 10s, 60s, 300s, 600s, 1800s, 3600s, 7200s, 14400s
- **Use**: Analyze connection lifetime distribution

### Message Metrics

#### `websocket_messages_sent_total`
- **Type**: Counter
- **Labels**: `channel` (head, blocks, gas_price)
- **Description**: Total number of messages sent to clients
- **Use**: Track message throughput per channel

#### `websocket_messages_dropped_total`
- **Type**: Counter
- **Labels**: `channel` (head, blocks, gas_price)
- **Description**: Messages dropped due to full client buffer
- **Use**: Identify slow consumers or buffer sizing issues

#### `websocket_message_broadcast_seconds`
- **Type**: Histogram
- **Labels**: `channel` (head, blocks, gas_price)
- **Description**: Time to broadcast message to all subscribed clients
- **Use**: Monitor broadcast performance

### Subscription Metrics

#### `websocket_subscriptions`
- **Type**: Gauge
- **Labels**: `channel` (head, blocks, gas_price)
- **Description**: Current number of active subscriptions per channel
- **Use**: Monitor subscription distribution

#### `websocket_subscribe_requests_total`
- **Type**: Counter
- **Labels**: `channel` (head, blocks, gas_price), `status` (success, error)
- **Description**: Total subscribe requests
- **Use**: Track subscription success/error rate

#### `websocket_unsubscribe_requests_total`
- **Type**: Counter
- **Labels**: `channel` (head, blocks, gas_price)
- **Description**: Total unsubscribe requests
- **Use**: Track unsubscribe activity

### Error Metrics

#### `websocket_errors_total`
- **Type**: Counter
- **Labels**: `type` (read, write, upgrade, unknown_method, unknown_channel)
- **Description**: Total WebSocket errors by type
- **Use**: Monitor error patterns

## Example Queries

### Connection Health
```promql
# Current active connections
websocket_active_connections

# Connection rate (last 5 minutes)
rate(websocket_connections_total[5m])

# Rejection rate
rate(websocket_connections_total{status="rejected"}[5m]) / rate(websocket_connections_total[5m])

# Average connection duration
histogram_quantile(0.5, rate(websocket_connection_duration_seconds_bucket[5m]))
```

### Message Throughput
```promql
# Messages sent per second by channel
rate(websocket_messages_sent_total[1m])

# Message drop rate (should be near 0)
rate(websocket_messages_dropped_total[5m])

# 95th percentile broadcast latency
histogram_quantile(0.95, rate(websocket_message_broadcast_seconds_bucket[5m]))
```

### Subscription Patterns
```promql
# Active subscriptions by channel
websocket_subscriptions

# Subscribe success rate
rate(websocket_subscribe_requests_total{status="success"}[5m]) / rate(websocket_subscribe_requests_total[5m])

# Churn rate (subscribes + unsubscribes)
rate(websocket_subscribe_requests_total[5m]) + rate(websocket_unsubscribe_requests_total[5m])
```

### Error Monitoring
```promql
# Total error rate
rate(websocket_errors_total[5m])

# Errors by type
sum by (type) (rate(websocket_errors_total[5m]))
```

## Alerting Rules

### Critical Alerts

```yaml
# High message drop rate - clients can't keep up
- alert: WebSocketHighMessageDropRate
  expr: rate(websocket_messages_dropped_total[5m]) > 10
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "High WebSocket message drop rate"
    description: "{{ $value }} messages/sec are being dropped"

# High connection rejection rate
- alert: WebSocketHighRejectionRate
  expr: rate(websocket_connections_total{status="rejected"}[5m]) / rate(websocket_connections_total[5m]) > 0.1
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "High WebSocket connection rejection rate"
    description: "{{ $value | humanizePercentage }} of connections are being rejected"

# Slow broadcast latency
- alert: WebSocketSlowBroadcast
  expr: histogram_quantile(0.95, rate(websocket_message_broadcast_seconds_bucket[5m])) > 1
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Slow WebSocket broadcast latency"
    description: "95th percentile broadcast time is {{ $value }}s"

# High error rate
- alert: WebSocketHighErrorRate
  expr: rate(websocket_errors_total[5m]) > 5
  for: 2m
  labels:
    severity: critical
  annotations:
    summary: "High WebSocket error rate"
    description: "{{ $value }} errors/sec"
```

## Grafana Dashboard

Example panel queries:

### Active Connections Panel
```promql
websocket_active_connections
```

### Message Rate Panel (Stacked Area)
```promql
sum by (channel) (rate(websocket_messages_sent_total[1m]))
```

### Broadcast Latency Panel (Heatmap)
```promql
rate(websocket_message_broadcast_seconds_bucket[5m])
```

### Error Rate Panel
```promql
sum by (type) (rate(websocket_errors_total[5m]))
```
