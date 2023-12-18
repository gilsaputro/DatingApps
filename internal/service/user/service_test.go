package user

import (
	"fmt"
	"gilsaputro/dating-apps/internal/store/user"
	mock_store "gilsaputro/dating-apps/internal/store/user/mock"
	"gilsaputro/dating-apps/pkg/hash"
	mock_hash "gilsaputro/dating-apps/pkg/hash/mock"
	"gilsaputro/dating-apps/pkg/token"
	mock_token "gilsaputro/dating-apps/pkg/token/mock"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewUserService(t *testing.T) {
	type args struct {
		store user.UserStoreMethod
		token token.TokenMethod
		hash  hash.HashMethod
	}
	tests := []struct {
		name string
		args args
		want UserServiceMethod
	}{
		{
			name: "success",
			args: args{
				store: &user.UserStore{},
				token: token.TokenConfig{},
				hash:  &hash.HashConfig{},
			},
			want: &UserService{
				store: &user.UserStore{},
				token: token.TokenConfig{},
				hash:  &hash.HashConfig{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(tt.args.store, tt.args.token, tt.args.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_AddUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mToken := mock_token.NewMockTokenMethod(mockCtrl)
	mHash := mock_hash.NewMockHashMethod(mockCtrl)
	mStore := mock_store.NewMockUserStoreMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		request AddUserServiceRequest
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
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username2").Return(user.UserStoreInfo{}, nil)

				bPassword := []byte("hash_password")
				mHash.EXPECT().HashValue("password2").Return(bPassword, nil)

				mStore.EXPECT().CreateUser(user.UserStoreInfo{
					Username: "username2",
					Password: "hash_password",
					Fullname: "fullname2",
					Email:    "email2",
				}).Return(nil)
			},
			args: args{
				request: AddUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username2",
					Password:     "password2",
					Fullname:     "fullname2",
					Email:        "email2",
				},
			},
			wantErr: false,
		},
		{
			name: "failed create to db flow",
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username2").Return(user.UserStoreInfo{}, nil)

				bPassword := []byte("hash_password")
				mHash.EXPECT().HashValue("password2").Return(bPassword, nil)

				mStore.EXPECT().CreateUser(user.UserStoreInfo{
					Username: "username2",
					Password: "hash_password",
					Fullname: "fullname2",
					Email:    "email2",
				}).Return(fmt.Errorf("some error"))
			},
			args: args{
				request: AddUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username2",
					Password:     "password2",
					Fullname:     "fullname2",
					Email:        "email2",
				},
			},
			wantErr: true,
		},
		{
			name: "failed hash flow",
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username2").Return(user.UserStoreInfo{}, nil)

				mHash.EXPECT().HashValue("password2").Return([]byte{}, fmt.Errorf("some error"))
			},
			args: args{
				request: AddUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username2",
					Password:     "password2",
					Fullname:     "fullname2",
					Email:        "email2",
				},
			},
			wantErr: true,
		},
		{
			name: "failed validate username flow",
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username2").Return(user.UserStoreInfo{}, fmt.Errorf("some error"))
			},
			args: args{
				request: AddUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username2",
					Password:     "password2",
					Fullname:     "fullname2",
					Email:        "email2",
				},
			},
			wantErr: true,
		},
		{
			name: "failed duplicate username flow",
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username2").Return(user.UserStoreInfo{
					UserId:   2,
					Username: "username2",
				}, nil)
			},
			args: args{
				request: AddUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username2",
					Password:     "password2",
					Fullname:     "fullname2",
					Email:        "email2",
				},
			},
			wantErr: true,
		},
		{
			name: "failed validate token flow",
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{}, fmt.Errorf("some error"))
			},
			args: args{
				request: AddUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username2",
					Password:     "password2",
					Fullname:     "fullname2",
					Email:        "email2",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserService{
				store: mStore,
				token: mToken,
				hash:  mHash,
			}
			tt.mockFunc()
			if err := service.AddUser(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("UserService.AddUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mToken := mock_token.NewMockTokenMethod(mockCtrl)
	mHash := mock_hash.NewMockHashMethod(mockCtrl)
	mStore := mock_store.NewMockUserStoreMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		request DeleteUserServiceRequest
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "success flow",
			args: args{
				request: DeleteUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username",
					Password:     "password",
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					Username: "username",
				}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username").Return(user.UserStoreInfo{
					UserId:   1,
					Username: "username",
					Password: "hash_password",
				}, nil)

				mHash.EXPECT().CompareValue("hash_password", "password").Return(true)

				mStore.EXPECT().DeleteUser(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "failed delete to db flow",
			args: args{
				request: DeleteUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username",
					Password:     "password",
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					Username: "username",
				}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username").Return(user.UserStoreInfo{
					UserId:   1,
					Username: "username",
					Password: "hash_password",
				}, nil)

				mHash.EXPECT().CompareValue("hash_password", "password").Return(true)

				mStore.EXPECT().DeleteUser(1).Return(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
		{
			name: "failed invalid password flow",
			args: args{
				request: DeleteUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username",
					Password:     "password",
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					Username: "username",
				}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username").Return(user.UserStoreInfo{
					UserId:   1,
					Username: "username",
					Password: "hash_password",
				}, nil)

				mHash.EXPECT().CompareValue("hash_password", "password").Return(false)
			},
			wantErr: true,
		},
		{
			name: "failed get user data flow",
			args: args{
				request: DeleteUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username",
					Password:     "password",
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					Username: "username",
				}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username").Return(user.UserStoreInfo{
					UserId:   1,
					Username: "username",
					Password: "hash_password",
				}, fmt.Errorf("some error"))
			},
			wantErr: true,
		},
		{
			name: "failed user not found flow",
			args: args{
				request: DeleteUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username",
					Password:     "password",
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					Username: "username",
				}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username").Return(user.UserStoreInfo{}, nil)
			},
			wantErr: true,
		},
		{
			name: "failed invalid token flow",
			args: args{
				request: DeleteUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username",
					Password:     "password",
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					Username: "username1",
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "failed error validate token flow",
			args: args{
				request: DeleteUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username",
					Password:     "password",
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{}, fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserService{
				store: mStore,
				token: mToken,
				hash:  mHash,
			}
			tt.mockFunc()
			if err := service.DeleteUser(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("UserService.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mToken := mock_token.NewMockTokenMethod(mockCtrl)
	mHash := mock_hash.NewMockHashMethod(mockCtrl)
	mStore := mock_store.NewMockUserStoreMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		request UpdateUserServiceRequest
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func()
		want     UserServiceInfo
		wantErr  bool
	}{
		{
			name: "success flow",
			args: args{
				request: UpdateUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username",
					Password:     "password",
					Fullname:     "fullname",
					Email:        "email",
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					Username: "username",
				}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username").Return(user.UserStoreInfo{
					UserId:   1,
					Username: "username",
					Password: "hash_password",
				}, nil)

				bPassword := []byte("hash_password")
				mHash.EXPECT().HashValue("password").Return(bPassword, nil)

				mStore.EXPECT().UpdateUser(gomock.Any()).Return(nil)
			},
			want: UserServiceInfo{
				UserId:   1,
				Username: "username",
				Fullname: "fullname",
				Email:    "email",
			},
			wantErr: false,
		},
		{
			name: "failed update flow",
			args: args{
				request: UpdateUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username",
					Password:     "password",
					Fullname:     "fullname",
					Email:        "email",
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					Username: "username",
				}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username").Return(user.UserStoreInfo{
					UserId:   1,
					Username: "username",
					Password: "hash_password",
				}, nil)

				bPassword := []byte("hash_password")
				mHash.EXPECT().HashValue("password").Return(bPassword, nil)

				mStore.EXPECT().UpdateUser(gomock.Any()).Return(fmt.Errorf("some error"))
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
		{
			name: "failed hash flow",
			args: args{
				request: UpdateUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username",
					Password:     "password",
					Fullname:     "fullname",
					Email:        "email",
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					Username: "username",
				}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username").Return(user.UserStoreInfo{
					UserId:   1,
					Username: "username",
					Password: "hash_password",
				}, nil)

				mHash.EXPECT().HashValue("password").Return([]byte{}, fmt.Errorf("some error"))
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
		{
			name: "failed invalid user flow",
			args: args{
				request: UpdateUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username",
					Password:     "password",
					Fullname:     "fullname",
					Email:        "email",
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					Username: "username",
				}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username").Return(user.UserStoreInfo{}, nil)
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
		{
			name: "failed get user flow",
			args: args{
				request: UpdateUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username",
					Password:     "password",
					Fullname:     "fullname",
					Email:        "email",
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					Username: "username",
				}, nil)

				mStore.EXPECT().GetUserInfoByUsername("username").Return(user.UserStoreInfo{}, fmt.Errorf("some error"))
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
		{
			name: "failed invalid token flow",
			args: args{
				request: UpdateUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username",
					Password:     "password",
					Fullname:     "fullname",
					Email:        "email",
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					Username: "abc",
				}, nil)
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
		{
			name: "failed validate token flow",
			args: args{
				request: UpdateUserServiceRequest{
					TokenRequest: "token_request",
					Username:     "username",
					Password:     "password",
					Fullname:     "fullname",
					Email:        "email",
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{}, fmt.Errorf("some error"))
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserService{
				store: mStore,
				token: mToken,
				hash:  mHash,
			}
			tt.mockFunc()
			got, err := service.UpdateUser(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.UpdateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mToken := mock_token.NewMockTokenMethod(mockCtrl)
	mHash := mock_hash.NewMockHashMethod(mockCtrl)
	mStore := mock_store.NewMockUserStoreMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		request GetByIDServiceRequest
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func()
		want     UserServiceInfo
		wantErr  bool
	}{
		{
			name: "success flow",
			args: args{
				request: GetByIDServiceRequest{
					TokenRequest: "token_request",
					UserId:       1,
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					UserID:   1,
					Username: "username",
				}, nil)

				mStore.EXPECT().GetUserInfoByID(1).Return(user.UserStoreInfo{
					UserId:      1,
					Username:    "username",
					Fullname:    "fullname",
					Email:       "email",
					CreatedDate: "2023-18-05",
				}, nil)
			},
			want: UserServiceInfo{
				UserId:      1,
				Username:    "username",
				Fullname:    "fullname",
				Email:       "email",
				CreatedDate: "2023-18-05",
			},
			wantErr: false,
		},
		{
			name: "failed on get flow",
			args: args{
				request: GetByIDServiceRequest{
					TokenRequest: "token_request",
					UserId:       1,
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					UserID:   1,
					Username: "username",
				}, nil)

				mStore.EXPECT().GetUserInfoByID(1).Return(user.UserStoreInfo{}, fmt.Errorf("some error"))
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
		{
			name: "failed invalid user flow",
			args: args{
				request: GetByIDServiceRequest{
					TokenRequest: "token_request",
					UserId:       1,
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					UserID:   1,
					Username: "username",
				}, nil)

				mStore.EXPECT().GetUserInfoByID(1).Return(user.UserStoreInfo{}, nil)
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
		{
			name: "failed invalid token flow",
			args: args{
				request: GetByIDServiceRequest{
					TokenRequest: "token_request",
					UserId:       1,
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{
					UserID:   2,
					Username: "username",
				}, nil)
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
		{
			name: "failed validate token flow",
			args: args{
				request: GetByIDServiceRequest{
					TokenRequest: "token_request",
					UserId:       1,
				},
			},
			mockFunc: func() {
				mToken.EXPECT().ValidateToken("token_request").Return(token.TokenBody{}, fmt.Errorf("some error"))
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserService{
				store: mStore,
				token: mToken,
				hash:  mHash,
			}
			tt.mockFunc()
			got, err := service.GetUserByID(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
