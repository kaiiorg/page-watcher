package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultWeb(t *testing.T) {
	// Act
	web := defaultWeb()

	// Assert
	require.NotNil(t, web)
	require.Equal(t, DefaultHostPort, web.Port)
}
