package api

import (
	"bytes"
	mockdb "easybank/db/mock"
	db "easybank/db/sqlc"
	"easybank/token"
	"easybank/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashed, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashed,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}

func TestCreateUser(t *testing.T) {
	rUser, password := randomUser(t)

	testcase := []struct {
		name          string
		body          fiber.Map
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(res *http.Response)
	}{
		{
			name: "OK",
			body: fiber.Map{
				"username":  rUser.Username,
				"password":  password,
				"full_name": rUser.FullName,
				"email":     rUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore) { // Store가 실행할 메소드와 전달될 파라미터 타입, 반환 타입, 실행 횟수 등을 미리 지정한다
				arg := db.CreateUserParams{
					Username: rUser.Username,
					FullName: rUser.FullName,
					Email:    rUser.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1) // sqlc - mysql 지원 상 return value type check 불가
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
				requireBodyMatcher(t, res.Body, rUser)
			},
		},
		{
			name: "InvalidUser",
			body: fiber.Map{
				"username":  "1234567890123",
				"password":  password,
				"full_name": rUser.FullName,
				"email":     rUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0) // CreateUser 호출 전 함수 종료됨
			},
			checkResponse: func(res *http.Response) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
				requireBodyMatcher(t, res.Body, rUser)
			},
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			testServer := newTestServer(t, store)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/user"
			request, err := http.NewRequest(fiber.MethodPost, url, bytes.NewReader(data)) // newBuffer와, newReader의 차이는 무엇인가?
			require.NoError(t, err)
			request.Header.Set("content-type", "application/json")

			res, err := testServer.router.Test(request)
			require.NoError(t, err)
			tc.checkResponse(res)
		})
	}
}

func TestLoginUser(t *testing.T) {
	owner, password := "ginber", "password"
	hashed, err := util.HashPassword(password)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		body          fiber.Map
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(resp *http.Response)
	}{
		{
			name: "OK",
			body: fiber.Map{
				"username": "camkpa",
				"password": "secret",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, owner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq("camkpa")).Times(1).Return(EqCreateUserParams(db.CreateUserParams{owner, hashed}), nil)
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

			request.Header.Set("content-type", "application/json")
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			resp, err := server.router.Test(request)
			require.NoError(t, err)
			tc.checkResponse(resp)
		})

	}
}

func requireBodyMatcher(t *testing.T, body io.ReadCloser, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	if err != nil {
		var errResponses []util.ErrorResponse
		err = json.Unmarshal(data, &errResponses) //"[{\"failedfield\":\"createUserRequest.Username\",\"tag\":\"max\",\"value\":\"12\"}]"
		require.NoError(t, err)
		require.NotNil(t, errResponses)
		return
	}

	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
