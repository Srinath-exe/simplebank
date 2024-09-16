package api

// import (
// 	"bytes"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	mockdb "github.com/Srinath-exe/simplebank/db/mock"
// 	db "github.com/Srinath-exe/simplebank/db/sqlc"
// 	"github.com/Srinath-exe/simplebank/token"
// 	"github.com/Srinath-exe/simplebank/util"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/require"
// )

// func TestGetAccountApi(t *testing.T) {
// 	user, _ := randomUser(t)
// 	account := randomAccount(user.Username)

// 	testCases := []struct {
// 		name       string
// 		accountID  int64
// 		buildStubs func(store *mockdb.MockStore)
// 		setupAuth  func(t *testing.T, request *http.Request, tokenMaker token.Maker)

// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name:      "OK",
// 			accountID: account.ID,
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
// 					Times(1).
// 					Return(account, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				requireBodyMatchAccount(t, recorder.Body, account)
// 			},
// 		},
// 		{
// 			name:      "Not Found",
// 			accountID: account.ID,
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
// 					Times(1).
// 					Return(db.Account{}, sql.ErrNoRows)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 			},
// 		},
// 		{
// 			name:      "Internal Error",
// 			accountID: account.ID, setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
// 					Times(1).
// 					Return(db.Account{}, sql.ErrConnDone)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 		{
// 			name:      "Invalid ID",
// 			accountID: 0, setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetAccount(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		}, {
// 			name:      "UnAuthorized User",
// 			accountID: account.ID, setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
// 					Times(1).
// 					Return(account, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name:      "No Authorization",
// 			accountID: account.ID,
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				// Do nothing
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetAccount(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := NewTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			url := fmt.Sprintf("/accounts/%d", tc.accountID)

// 			request, err := http.NewRequest(http.MethodGet, url, nil)
// 			require.NoError(t, err)

// 			tc.setupAuth(t, request, server.tokenMaker)

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(t, recorder)
// 		})
// 	}
// }

// func randomAccount(owner string) db.Account {
// 	return db.Account{
// 		ID:       util.RandomInt(1, 1000),
// 		Owner:    owner,
// 		Currency: util.RandomCurrency(),
// 		Balance:  util.RandomMoney(),
// 	}
// }

// func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
// 	data, err := io.ReadAll(body)
// 	require.NoError(t, err)

// 	var gotAccount db.Account
// 	err = json.Unmarshal(data, &gotAccount)
// 	require.NoError(t, err)

// 	require.Equal(t, account, gotAccount)
// }
