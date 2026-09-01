package api

import (
	"testing"
	"time"

	db "github.com/MidNight91119/freelance-marketplace/internal/db/sqlc"
	"github.com/MidNight91119/freelance-marketplace/internal/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}
