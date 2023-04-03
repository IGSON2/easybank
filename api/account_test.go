package api

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	mockdb "easybank/db/mock"
	db "easybank/db/sqlc"
	"easybank/util"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	account := randomAccount()
	testCases := []struct {
		name         string
		accountID    int64
		buildStub    func(store *mockdb.MockStore)
		expectStatus int
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(randomAccount(), nil)
			},
			expectStatus: http.StatusOK,
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)). // gomock의 Any와 Eq의 역할은..?
					Times(1).                                        // 함수가 1번 호출된다.
					Return(db.Account{}, sql.ErrNoRows)
			},
			expectStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrler := gomock.NewController(t)
			defer ctrler.Finish()

			store := mockdb.NewMockStore(ctrler)
			tc.buildStub(store)

			server := newTestServer(t, store)
			// recoder := httptest.NewRecorder() // 테스트를 위해 실제로 http 서버를 실행할 필요가 없도록 도와주는 recoder 생성

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			resp, err := server.router.Test(request)
			require.NoError(t, err)
			require.Equal(t, tc.expectStatus, resp.StatusCode, tc.name)
		})

	}
}

func TestCreateAccount(t *testing.T) {
	user, password := randomUser(t)

	testcases := []struct {
		name           string
		body           fiber.Map
		buildstub      func(store *mockdb.MockStore)
		expectedStatus int
	}{
		{
			name: "OK",
			body: fiber.Map{
				"owner":    user.Owner,
				"currency": user.Currency,
			},
			buildstub: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    user.Owner,
					Balance:  int64(0),
					Currency: user.Currency,
				}

				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(new(driver.RowsAffected), sql.ErrNoRows)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "InvalidCurrency",
			body: fiber.Map{
				"owner":   user.Owner,
				"curreny": "invalid",
			},
			buildstub: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(new(driver.RowsAffected), sql.ErrNoRows)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildstub(store)

			server := newTestServer(t, store)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data)) // bytes.Reader로써 전달
			require.NoError(t, err)

			res, err := server.router.Test(request)
			require.NoError(t, err)
			require.Equal(t, tc.expectedStatus, res.StatusCode, tc.name)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:        util.RandomInt(1, 1000),
		Owner:     util.RandomOwner(),
		Balance:   util.RandomBalance(),
		Currency:  util.RandomCurrency(),
		CreatedAt: time.Now(),
	}
}
