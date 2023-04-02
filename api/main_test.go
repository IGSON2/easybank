package api

import (
	db "easybank/db/sqlc"
	"easybank/token"
	"easybank/util"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

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

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func createAndSetAuthToken(t *testing.T, req *http.Request, tokenMaker token.Maker, username string) {
	if len(username) == 0 {
		return
	}

	token, _, err := tokenMaker.CreateToken(username, time.Minute)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
	req.Header.Set(authorizationHeaderKey, authorizationHeader)
}
