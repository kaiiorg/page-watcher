package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultDb(t *testing.T) {
	// Act
	db := defaultDb()

	// Assert
	require.NotNil(t, db)
	require.Equal(t, DEFAULT_DB_PATH, db.Path)
}
