package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/MidNight91119/freelance-marketplace/internal/db/mock"
	db "github.com/MidNight91119/freelance-marketplace/internal/db/sqlc"
	"github.com/MidNight91119/freelance-marketplace/internal/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		ID: util.RandomInt(1, 1000),
		Name: util.RandomName(),
		Email: util.RandomEmail(),
		HashedPassword: hashedPassword,
		Role: db.RolesClient,
	}

	return user, password
}

func TestSignupAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mockdb.NewMockStore(ctrl)
	user, password := randomUser(t)

	store.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Times(1).
		Return(user, nil)

	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	body, _ := json.Marshal(map[string]any{
		"name": user.Name,
		"email": user.Email,
		"password": password,
		"role": user.Role,
	})
	request, _ := http.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusCreated, recorder.Code)
}
