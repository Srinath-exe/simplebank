package api

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	db "github.com/Srinath-exe/simplebank/db/sqlc"
	"github.com/Srinath-exe/simplebank/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
