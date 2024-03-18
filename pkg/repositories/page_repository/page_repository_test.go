package page_repository

import (
	"fmt"
	"github.com/kaiiorg/page-watcher/pkg/config"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewSqlitePageRepository_NoError(t *testing.T) {
	// Arrange
	dbPath := fmt.Sprintf("%s.db", uuid.NewString())
	defer os.Remove(dbPath)
	c := &config.DB{
		Path:   dbPath,
		Retain: 2,
	}

	// Act
	spr, err := NewSqlitePageRepository(c)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, spr)
	require.NotNil(t, spr.db)
}
