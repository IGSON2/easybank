package api

import (
	"easybank/token"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func addAuthorization(t *testing.T, req *http.Request, tokenMaker token.Maker, authorizationType, username string, duration time.Duration) {
	token, payload, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	autorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	req.Header.Set(authorizationHeaderKey, autorizationHeader)
}
