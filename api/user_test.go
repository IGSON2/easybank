package api

import (
	"bytes"
	mockdb "easybank/db/mock"
	db "easybank/db/sqlc"
	"easybank/util"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestLoginUser(t *testing.T) {
	hashed, err := util.HashPassword("camkpa")
	require.NoError(t, err)

	user := db.User{
		Username:       "camkpa",
		HashedPassword: hashed,
		FullName:       "qfgcoc",
		Email:          "uzfrdi@email.com",
	}

	testCases := []struct {
		name          string
		body          fiber.Map
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(resp *http.Response)
	}{
		{
			name: "OK",
			body: fiber.Map{
				"username": "camkpa",
				"password": "secret",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq("camkpa")).Times(1).Return(user, nil)
			},
			checkResponse: func(resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrler := gomock.NewController(t)
			defer ctrler.Finish()

			store := mockdb.NewMockStore(ctrler)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			// recoder := httptest.NewRecorder() // 테스트를 위해 실제로 http 서버를 실행할 필요가 없도록 도와주는 recoder 생성

			url := "user/login"
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			resp, err := server.router.Test(request)
			require.NoError(t, err)
			tc.checkResponse(resp)
		})

	}
}
