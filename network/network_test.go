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
