package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/MidNight91119/freelance-marketplace/internal/token"
	"github.com/MidNight91119/freelance-marketplace/internal/util"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	req *http.Request,
	tokenMaker token.Maker,
	userID int64,
	role string,
) {
	accessToken, payload, err := tokenMaker.CreateToken(userID, util.RandomEmail(), role, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authHeader := fmt.Sprintf("%s %s", authTypeBearer, accessToken)

	req.Header.Set(authHeaderKey, authHeader)
}
