package postgres

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPostgresIntegration(t *testing.T) {
	cfg := &Config{
		"localhost",
		"5432",
		"postgres",
		"12345",
		"postgres",
	}
	db, err := New(cfg)
	require.NoError(t, err, "Could not connect to postgres")
	require.NotNil(t, db, "Could not connect to postgres")
	defer db.Close()

	err = db.Ping()
	assert.NoError(t, err, "Could not ping postgres")
}

func TestPostgresIncorrectConnection(t *testing.T) {
	cfg := &Config{
		Host:     "invalidhost",
		Port:     "5432",
		Username: "postgres",
		Password: "1234",
		Database: "postgres",
	}
	db, err := New(cfg)
	require.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "failed to connect to postgres")
}
