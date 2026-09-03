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
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		ID:             util.RandomInt(1, 1000),
		Name:           util.RandomName(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
		Role:           db.RolesClient,
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
		"name":     user.Name,
		"email":    user.Email,
		"password": password,
		"role":     user.Role,
	})
	request, _ := http.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestLoginAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          map[string]any
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]any{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), user.Email).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "WRONG_PASSWORD",
			body: map[string]any{
				"email":    user.Email,
				"password": "wrongPassword",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), user.Email).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UNKNOWN_EMAIL",
			body: map[string]any{
				"email":    "wrong@email.com",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			body, err := json.Marshal(tc.body)
			require.NoError(t, err)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}
