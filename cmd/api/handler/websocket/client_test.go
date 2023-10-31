package websocket

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotifyClosedClient(t *testing.T) {
	client := newClient(10, nil)
	err := client.Close()
	require.NoError(t, err, "closing client")
	client.Notify("test")
}
