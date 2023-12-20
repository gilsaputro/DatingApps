package middleware

import (
	"fmt"
	"gilsaputro/dating-apps/internal/store/user"
	"gilsaputro/dating-apps/pkg/token"
	mock_token "gilsaputro/dating-apps/pkg/token/mock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewMiddleware(t *testing.T) {
	type args struct {
		tokenMethod token.TokenMethod
		userStore   user.UserStore
	}
	tests := []struct {
		name string
		args args
		want Middleware
	}{
		{
			name: "success",
			args: args{
				tokenMethod: token.TokenConfig{},
				userStore:   user.UserStore{},
			},
			want: Middleware{
				tokenMethod: token.TokenConfig{},
				userStore:   &user.UserStore{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMiddleware(tt.args.tokenMethod, &tt.args.userStore); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiddleware_MiddlewareVerifyToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mToken := mock_token.NewMockTokenMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		token string
		path  string
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.Context().Value("id").(int)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
		}
		w.Header().Set("id", fmt.Sprintf("%v", token))
		_, _ = w.Write([]byte{})
	})

	tests := []struct {
		name      string
		args      args
		mockFunc  func()
		wantToken string
	}{
		{
			name: "success flow",
			args: args{
				token: "token_baru",
				path:  "/user",
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_baru").Return(token.TokenBody{
					UserID: 1,
				}, nil)
			},
			wantToken: "1",
		},
		{
			name: "invalid token flow",
			args: args{
				token: "",
				path:  "/user",
			},
			mockFunc:  func() {},
			wantToken: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Middleware{
				tokenMethod: mToken,
			}

			tt.mockFunc()
			middleware := m.MiddlewareVerifyToken(next)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, tt.args.path, nil)
			if len(tt.args.token) > 0 {
				request.Header.Set("Authorization", "Bearer "+tt.args.token)
			}

			middleware(recorder, request)
			new_token := recorder.Header().Get("id")
			if !reflect.DeepEqual(new_token, tt.wantToken) {
				t.Errorf("NewMiddleware() token = %v, want %v", new_token, tt.wantToken)
			}
		})
	}
}
