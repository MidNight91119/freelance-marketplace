package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/MidNight91119/freelance-marketplace/internal/db/mock"
	db "github.com/MidNight91119/freelance-marketplace/internal/db/sqlc"
	"github.com/MidNight91119/freelance-marketplace/internal/token"
	"github.com/MidNight91119/freelance-marketplace/internal/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListMyProposalsAPI(t *testing.T) {
	userID := util.RandomInt(1, 1000)

	testcases := []struct {
		name          string
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		buildstubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, userID, roleFreelancer)
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProposalsByFreelancer(gomock.Any(), userID).
					Times(1).
					Return([]db.ListProposalsByFreelancerRow{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "FORBIDDEN",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, userID, roleClient)
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProposalsByFreelancer(gomock.Any(), userID).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "UNAUTHORIZED",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProposalsByFreelancer(gomock.Any(), userID).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildstubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/api/proposals/mine", nil)
			tc.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
