package network

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartAndStop(t *testing.T) {
	host, err := NewHost(&Config{})
	require.NoError(t, err)

	err = host.Start()
	require.NoError(t, err)

	err = host.Stop()
	require.NoError(t, err)
}

func TestPubSub(t *testing.T) {
	hostA, err := NewHost(&Config{})
	require.NoError(t, err)

	err = hostA.Start()
	require.NoError(t, err)

	hostB, err := NewHost(&Config{})
	require.NoError(t, err)

	err = hostB.Start()
	require.NoError(t, err)

	defer func() {
		err = hostA.Stop()
		require.NoError(t, err)

		err = hostB.Stop()
		require.NoError(t, err)
	}()

	addrB := hostB.AddrInfo()

	err = hostA.Connect(addrB)
	require.NoError(t, err)

	data := []byte("nootwashere")
	err = hostB.Publish(data)
	require.NoError(t, err)

	msg, err := hostA.Next()
	require.NoError(t, err)
	require.Equal(t, data, msg.From)
}
