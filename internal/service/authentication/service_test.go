package authentication

import (
	"fmt"
	"gilsaputro/dating-apps/internal/store/user"
	mock_user "gilsaputro/dating-apps/internal/store/user/mock"
	"gilsaputro/dating-apps/models"
	"gilsaputro/dating-apps/pkg/hash"
	mock_hash "gilsaputro/dating-apps/pkg/hash/mock"
	"gilsaputro/dating-apps/pkg/token"
	mock_token "gilsaputro/dating-apps/pkg/token/mock"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
)

func TestNewAuthenticationService(t *testing.T) {
	type args struct {
		store user.UserStoreMethod
		token token.TokenMethod
		hash  hash.HashMethod
	}
	tests := []struct {
		name string
		args args
		want AuthenticationServiceMethod
	}{
		{
			name: "success flow",
			args: args{
				store: &user.UserStore{},
				token: &token.TokenConfig{},
				hash:  &hash.HashConfig{},
			},
			want: &AuthenticationService{
				store: &user.UserStore{},
				token: &token.TokenConfig{},
				hash:  &hash.HashConfig{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthenticationService(tt.args.store, tt.args.token, tt.args.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthenticationService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthenticationService_Login(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	uStore := mock_user.NewMockUserStoreMethod(mockCtrl)
	mToken := mock_token.NewMockTokenMethod(mockCtrl)
	mHash := mock_hash.NewMockHashMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		request LoginServiceRequest
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		want     string
		wantErr  bool
	}{
		{
			name: "success flow",
			mockFunc: func() {
				uStore.EXPECT().GetUserInfoByUsername("username").Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username:   "username",
					Password:   "password",
					Fullname:   "fullname",
					Email:      "email",
					IsVerified: true,
				}, nil)

				mHash.EXPECT().CompareValue("password", "password").Return(true)

				mToken.EXPECT().GenerateToken(token.TokenBody{
					UserID: int(1),
				}).Return("token", nil)
			},
			args: args{
				request: LoginServiceRequest{
					Username: "username",
					Password: "password",
				},
			},
			want:    "token",
			wantErr: false,
		},
		{
			name: "error password flow",
			mockFunc: func() {
				uStore.EXPECT().GetUserInfoByUsername("username").Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username:   "username",
					Password:   "password",
					Fullname:   "fullname",
					Email:      "email",
					IsVerified: true,
				}, nil)

				mHash.EXPECT().CompareValue("password", "password").Return(false)
			},
			args: args{
				request: LoginServiceRequest{
					Username: "username",
					Password: "password",
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "error get info flow",
			mockFunc: func() {
				uStore.EXPECT().GetUserInfoByUsername("username").Return(models.User{}, fmt.Errorf("some error"))
			},
			args: args{
				request: LoginServiceRequest{
					Username: "username",
					Password: "password",
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "error invalid user flow",
			mockFunc: func() {
				uStore.EXPECT().GetUserInfoByUsername("username").Return(models.User{}, nil)
			},
			args: args{
				request: LoginServiceRequest{
					Username: "username",
					Password: "password",
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAuthenticationService(uStore, mToken, mHash)
			tt.mockFunc()
			got, err := s.Login(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthenticationService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthenticationService.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthenticationService_Register(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	uStore := mock_user.NewMockUserStoreMethod(mockCtrl)
	mToken := mock_token.NewMockTokenMethod(mockCtrl)
	mHash := mock_hash.NewMockHashMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		request RegisterServiceRequest
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		wantErr  bool
	}{
		{
			name: "success flow",
			mockFunc: func() {
				uStore.EXPECT().GetUserInfoByUsername("username").Return(models.User{}, nil)

				mHash.EXPECT().HashValue("password").Return([]byte("hash"), nil)
				uStore.EXPECT().CreateUser(gomock.Any()).Return(nil)
			},
			args: args{
				request: RegisterServiceRequest{
					Username: "username",
					Password: "password",
					Fullname: "fullname",
					Email:    "email",
				},
			},
			wantErr: false,
		},
		{
			name: "error hash password flow",
			mockFunc: func() {
				uStore.EXPECT().GetUserInfoByUsername("username").Return(models.User{}, nil)

				mHash.EXPECT().HashValue("password").Return([]byte("hash"), fmt.Errorf("some error"))
			},
			args: args{
				request: RegisterServiceRequest{
					Username: "username",
					Password: "password",
					Fullname: "fullname",
					Email:    "email",
				},
			},
			wantErr: true,
		},
		{
			name: "error check user flow",
			mockFunc: func() {
				uStore.EXPECT().GetUserInfoByUsername("username").Return(models.User{}, fmt.Errorf("some error"))
			},
			args: args{
				request: RegisterServiceRequest{
					Username: "username",
					Password: "password",
					Fullname: "fullname",
					Email:    "email",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAuthenticationService(uStore, mToken, mHash)
			tt.mockFunc()
			if err := s.Register(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("AuthenticationService.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
