package user

import (
	"context"
	"fmt"
	"gilsaputro/dating-apps/internal/service/user"
	"gilsaputro/dating-apps/internal/service/user/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUserHandler_UpdateUserHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	m := mock.NewMockUserServiceMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		userID  int
		body    string
		timeout int
	}
	type want struct {
		body string
		code int
	}
	tests := []struct {
		name        string
		args        args
		mockFunc    func()
		mockContext func() (context.Context, func())
		want        want
	}{
		{
			name: "success flow",
			args: args{
				userID: 1,
				body: `{
					"password": "pas1"
				}`,
				timeout: 5,
			},
			mockFunc: func() {
				m.EXPECT().UpdateUser(user.UpdateUserServiceRequest{
					UserId:   1,
					Password: "pas1",
				}).Return(user.UserServiceInfo{
					UserId:     1,
					Username:   "username",
					Fullname:   "full name",
					Email:      "email.com",
					IsVerified: true,
				}, nil)
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 200,
				body: `{"data":{"username":"username","email":"email.com","fullname":"full name"},"code":200,"message":"success"}`,
			},
		},
		{
			name: "error on service flow",
			args: args{
				userID: 1,
				body: `{
					"password": "pas1"
				}`,
				timeout: 5,
			},
			mockFunc: func() {
				m.EXPECT().UpdateUser(user.UpdateUserServiceRequest{
					UserId:   1,
					Password: "pas1",
				}).Return(user.UserServiceInfo{}, fmt.Errorf("some error"))
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 500,
				body: `{"code":500,"message":"some error"}`,
			},
		},
		{
			name: "error on service flow user not exists",
			args: args{
				userID: 1,
				body: `{
					"password": "pas1"
				}`,
				timeout: 5,
			},
			mockFunc: func() {
				m.EXPECT().UpdateUser(user.UpdateUserServiceRequest{
					UserId:   1,
					Password: "pas1",
				}).Return(user.UserServiceInfo{}, user.ErrUserNameNotExists)
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 400,
				body: `{"code":400,"message":"username is not exists"}`,
			},
		},
		{
			name: "error on invalid body value#2",
			args: args{
				userID: 1,
				body: `{
				}`,
				timeout: 5,
			},
			mockFunc: func() {
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 400,
				body: `{"code":400,"message":"Invalid Parameter Request"}`,
			},
		},
		{
			name: "error on invalid body value",
			args: args{
				userID:  1,
				body:    `{`,
				timeout: 5,
			},
			mockFunc: func() {
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 400,
				body: `{"code":400,"message":"Bad Request"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			defer mockCtrl.Finish()
			handler := UserHandler{
				service:      m,
				timeoutInSec: tt.args.timeout,
			}
			r := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(tt.args.body))
			ctx, cancel := tt.mockContext()
			defer cancel()
			r = r.WithContext(ctx)
			if tt.args.userID > 0 {
				r = r.WithContext(context.WithValue(r.Context(), "id", tt.args.userID))
			}
			w := httptest.NewRecorder()
			handler.EditUserHandler(w, r)
			result := w.Result()
			resBody, err := ioutil.ReadAll(result.Body)

			if err != nil {
				t.Fatalf("Error read body err = %v\n", err)
			}

			if string(resBody) != tt.want.body {
				t.Fatalf("GetStatHandler body got =%s, want %s \n", string(resBody), tt.want.body)
			}

			if result.StatusCode != tt.want.code {
				t.Fatalf("GetStatHandler status code got =%d, want %d \n", result.StatusCode, tt.want.code)
			}
		})
	}
}
