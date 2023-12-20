package user

import (
	"fmt"
	"gilsaputro/dating-apps/internal/store/user"
	"gilsaputro/dating-apps/internal/store/user/mock"
	"gilsaputro/dating-apps/models"
	"gilsaputro/dating-apps/pkg/hash"
	mock_hash "gilsaputro/dating-apps/pkg/hash/mock"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
)

func TestNewUserService(t *testing.T) {
	type args struct {
		store user.UserStoreMethod
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
				hash:  &hash.HashConfig{},
			},
			want: &UserService{
				store: &user.UserStore{},
				hash:  &hash.HashConfig{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(tt.args.store, tt.args.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mHash := mock_hash.NewMockHashMethod(mockCtrl)
	mStore := mock.NewMockUserStoreMethod(mockCtrl)
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
					UserId:   1,
					Password: "password",
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "username",
					Password: "hash_password",
				}, nil)

				mHash.EXPECT().CompareValue("hash_password", "password").Return(true)

				mStore.EXPECT().DeleteUser(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "failed delete flow",
			args: args{
				request: DeleteUserServiceRequest{
					UserId:   1,
					Password: "password",
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "username",
					Password: "hash_password",
				}, nil)

				mHash.EXPECT().CompareValue("hash_password", "password").Return(true)

				mStore.EXPECT().DeleteUser(1).Return(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
		{
			name: "invalid password flow",
			args: args{
				request: DeleteUserServiceRequest{
					UserId:   1,
					Password: "password",
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "username",
					Password: "hash_password",
				}, nil)

				mHash.EXPECT().CompareValue("hash_password", "password").Return(false)
			},
			wantErr: true,
		},
		{
			name: "error get info",
			args: args{
				request: DeleteUserServiceRequest{
					UserId:   1,
					Password: "password",
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "username",
					Password: "hash_password",
				}, fmt.Errorf("some error"))
			},
			wantErr: true,
		},
		{
			name: "invalid userid",
			args: args{
				request: DeleteUserServiceRequest{
					UserId:   0,
					Password: "password",
				},
			},
			mockFunc: func() {
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserService{
				store: mStore,
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
	mHash := mock_hash.NewMockHashMethod(mockCtrl)
	mStore := mock.NewMockUserStoreMethod(mockCtrl)
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
					UserId:   1,
					Username: "username",
					Password: "pass",
					Fullname: "full",
					Email:    "email",
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "username",
					Password: "hash_password",
				}, nil)

				mHash.EXPECT().HashValue("pass").Return([]byte(`hash_password`), nil)

				mStore.EXPECT().UpdateUser(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "username",
					Password: "hash_password",
					Fullname: "full",
					Email:    "email",
				}).Return(nil)
			},
			want: UserServiceInfo{
				UserId:      1,
				Username:    "username",
				Fullname:    "full",
				Email:       "email",
				CreatedDate: "0001-01-01 00:00:00 +0000 UTC",
			},
			wantErr: false,
		},
		{
			name: "error update flow",
			args: args{
				request: UpdateUserServiceRequest{
					UserId:   1,
					Username: "username",
					Password: "pass",
					Fullname: "full",
					Email:    "email",
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "username",
					Password: "hash_password",
				}, nil)

				mHash.EXPECT().HashValue("pass").Return([]byte(`hash_password`), nil)

				mStore.EXPECT().UpdateUser(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "username",
					Password: "hash_password",
					Fullname: "full",
					Email:    "email",
				}).Return(fmt.Errorf("some error"))
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
		{
			name: "error invalid passsword value flow",
			args: args{
				request: UpdateUserServiceRequest{
					UserId:   1,
					Username: "username",
					Password: "pass",
					Fullname: "full",
					Email:    "email",
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "username",
					Password: "hash_password",
				}, nil)

				mHash.EXPECT().HashValue("pass").Return([]byte(``), fmt.Errorf("some error"))
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
		{
			name: "error get info flow",
			args: args{
				request: UpdateUserServiceRequest{
					UserId:   1,
					Username: "username",
					Password: "pass",
					Fullname: "full",
					Email:    "email",
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{}, fmt.Errorf("some error"))
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
		{
			name: "error invalid userid flow",
			args: args{
				request: UpdateUserServiceRequest{},
			},
			mockFunc: func() {
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserService{
				store: mStore,
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
	mHash := mock_hash.NewMockHashMethod(mockCtrl)
	mStore := mock.NewMockUserStoreMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		request GetByIDServiceRequest
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		want     UserServiceInfo
		wantErr  bool
	}{
		{
			name: "success flow",
			args: args{
				request: GetByIDServiceRequest{
					UserId: 1,
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "username",
					Fullname: "full",
					Email:    "email",
				}, nil)
			},
			want: UserServiceInfo{
				UserId:      1,
				Username:    "username",
				Fullname:    "full",
				Email:       "email",
				CreatedDate: "0001-01-01 00:00:00 +0000 UTC",
			},
			wantErr: false,
		},
		{
			name: "error flow",
			args: args{
				request: GetByIDServiceRequest{
					UserId: 1,
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{}, fmt.Errorf("some error"))
			},
			want:    UserServiceInfo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserService{
				store: mStore,
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

func TestUserService_UpgradeUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mHash := mock_hash.NewMockHashMethod(mockCtrl)
	mStore := mock.NewMockUserStoreMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		request UpgradeServiceRequest
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		wantErr  bool
	}{
		{
			name: "success flow",
			args: args{
				request: UpgradeServiceRequest{
					UserId:   1,
					Password: "password",
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "username",
					Password: "hash_password",
				}, nil)

				mHash.EXPECT().CompareValue("hash_password", "password").Return(true)

				mStore.EXPECT().UpdateUser(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error flow",
			args: args{
				request: UpgradeServiceRequest{
					UserId:   1,
					Password: "password",
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "username",
					Password: "hash_password",
				}, nil)

				mHash.EXPECT().CompareValue("hash_password", "password").Return(true)

				mStore.EXPECT().UpdateUser(gomock.Any()).Return(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
		{
			name: "error invalid password flow",
			args: args{
				request: UpgradeServiceRequest{
					UserId:   1,
					Password: "password",
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "username",
					Password: "hash_password",
				}, nil)

				mHash.EXPECT().CompareValue("hash_password", "password").Return(false)
			},
			wantErr: true,
		},
		{
			name: "error user is verified flow",
			args: args{
				request: UpgradeServiceRequest{
					UserId:   1,
					Password: "password",
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username:   "username",
					Password:   "hash_password",
					IsVerified: true,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "error on get data flow",
			args: args{
				request: UpgradeServiceRequest{
					UserId:   1,
					Password: "password",
				},
			},
			mockFunc: func() {
				mStore.EXPECT().GetUserInfoByID(int(1)).Return(models.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "username",
					Password: "hash_password",
				}, fmt.Errorf("some error"))
			},
			wantErr: true,
		},
		{
			name: "error on get data flow",
			args: args{
				request: UpgradeServiceRequest{},
			},
			mockFunc: func() {
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserService{
				store: mStore,
				hash:  mHash,
			}
			tt.mockFunc()
			if err := service.UpgradeUser(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("UserService.UpgradeUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
