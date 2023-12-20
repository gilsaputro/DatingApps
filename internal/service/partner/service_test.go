package partner

import (
	"fmt"
	"gilsaputro/dating-apps/internal/store/partnercache"
	mock_partner "gilsaputro/dating-apps/internal/store/partnercache/mock"
	"gilsaputro/dating-apps/internal/store/user"
	mock_user "gilsaputro/dating-apps/internal/store/user/mock"
	"gilsaputro/dating-apps/internal/store/userhistory"
	mock_userhist "gilsaputro/dating-apps/internal/store/userhistory/mock"
	"gilsaputro/dating-apps/models"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
)

func TestNewPartnerService(t *testing.T) {
	type args struct {
		storeUser  user.UserStoreMethod
		storeHist  userhistory.UserHistoryStoreMethod
		cache      partnercache.PartnerCacheStoreMethod
		maxCounter int
	}
	tests := []struct {
		name string
		args args
		want PartnerServiceMethod
	}{
		{
			name: "success flow",
			args: args{
				storeUser: &user.UserStore{},
				storeHist: &userhistory.UserHistoryStore{},
				cache:     &partnercache.PartnerCacheStore{},
			},
			want: &PartnerService{
				storeUser:  &user.UserStore{},
				storeHist:  &userhistory.UserHistoryStore{},
				cache:      &partnercache.PartnerCacheStore{},
				maxCounter: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPartnerService(tt.args.storeUser, tt.args.storeHist, tt.args.cache, tt.args.maxCounter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPartnerService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPartnerService_PassPartner(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	uStore := mock_user.NewMockUserStoreMethod(mockCtrl)
	hStore := mock_userhist.NewMockUserHistoryStoreMethod(mockCtrl)
	pStore := mock_partner.NewMockPartnerCacheStoreMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		request PartnerServiceRequest
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		want     PartnerServiceInfo
		wantErr  bool
	}{
		{
			name: "success flow",
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			mockFunc: func() {
				pStore.EXPECT().GetViewedUserCounter("1").Return("1", nil)
				uStore.EXPECT().Count().Return(4, nil)
				pStore.EXPECT().GetViewedPartnerHistory("1").Return("2,3", nil)
				pStore.EXPECT().SetViewedPartnerHistory("1", "2,3,4").Return(nil)
				pStore.EXPECT().SetCurentPartnerState(1, 4).Return(nil)

				uStore.EXPECT().GetUserInfoByID(4).Return(models.User{
					Model: gorm.Model{
						ID: 4,
					},
					Username:   "U4",
					Fullname:   "F4",
					Email:      "E4",
					IsVerified: true,
				}, nil)

				hStore.EXPECT().CountByUserIDAndPartnerID(1, 4).Return(0, nil)

				pStore.EXPECT().SetViewedUserCounter("1", "2").Return(nil)
			},
			want: PartnerServiceInfo{
				PartnerID:   4,
				Fullname:    "F4",
				IsVerified:  true,
				Status:      "PENDING",
				CreatedDate: "0001-01-01 00:00:00 +0000 UTC",
			},
			wantErr: false,
		},
		{
			name: "error on get partner info flow",
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			mockFunc: func() {
				pStore.EXPECT().GetViewedUserCounter("1").Return("1", nil)
				uStore.EXPECT().Count().Return(4, nil)
				pStore.EXPECT().GetViewedPartnerHistory("1").Return("2,3", nil)
				pStore.EXPECT().SetViewedPartnerHistory("1", "2,3,4").Return(nil)
				pStore.EXPECT().SetCurentPartnerState(1, 4).Return(nil)

				uStore.EXPECT().GetUserInfoByID(4).Return(models.User{
					Model: gorm.Model{
						ID: 4,
					},
					Username:   "U4",
					Fullname:   "F4",
					Email:      "E4",
					IsVerified: true,
				}, fmt.Errorf("some error"))
			},
			want:    PartnerServiceInfo{},
			wantErr: true,
		},
		{
			name: "error on SetCurentPartnerState flow",
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			mockFunc: func() {
				pStore.EXPECT().GetViewedUserCounter("1").Return("1", nil)
				uStore.EXPECT().Count().Return(4, nil)
				pStore.EXPECT().GetViewedPartnerHistory("1").Return("2,3", nil)
				pStore.EXPECT().SetViewedPartnerHistory("1", "2,3,4").Return(nil)
				pStore.EXPECT().SetCurentPartnerState(1, 4).Return(fmt.Errorf("some error"))
			},
			want:    PartnerServiceInfo{},
			wantErr: true,
		},
		{
			name: "error on SetViewedPartnerHistory flow",
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			mockFunc: func() {
				pStore.EXPECT().GetViewedUserCounter("1").Return("1", nil)
				uStore.EXPECT().Count().Return(4, nil)
				pStore.EXPECT().GetViewedPartnerHistory("1").Return("2,3", nil)
				pStore.EXPECT().SetViewedPartnerHistory("1", "2,3,4").Return(fmt.Errorf("some error"))
			},
			want:    PartnerServiceInfo{},
			wantErr: true,
		},
		{
			name: "error on SetViewedPartnerHistory flow",
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			mockFunc: func() {
				pStore.EXPECT().GetViewedUserCounter("1").Return("1", nil)
				uStore.EXPECT().Count().Return(4, nil)
				pStore.EXPECT().GetViewedPartnerHistory("1").Return("2,3", fmt.Errorf("some error"))
			},
			want:    PartnerServiceInfo{},
			wantErr: true,
		},
		{
			name: "error on Count flow",
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			mockFunc: func() {
				pStore.EXPECT().GetViewedUserCounter("1").Return("1", nil)
				uStore.EXPECT().Count().Return(4, fmt.Errorf("some error"))
			},
			want:    PartnerServiceInfo{},
			wantErr: true,
		},
		{
			name: "error on Count flow",
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			mockFunc: func() {
				pStore.EXPECT().GetViewedUserCounter("1").Return("1", fmt.Errorf("some error"))
			},
			want:    PartnerServiceInfo{},
			wantErr: true,
		},
		{
			name: "error on max counter flow",
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			mockFunc: func() {
				pStore.EXPECT().GetViewedUserCounter("1").Return("100", nil)
			},
			want:    PartnerServiceInfo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPartnerService(uStore, hStore, pStore, 10)
			tt.mockFunc()
			got, err := s.PassPartner(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("PartnerService.PassPartner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PartnerService.PassPartner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPartnerService_GetCurrentPartner(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	uStore := mock_user.NewMockUserStoreMethod(mockCtrl)
	hStore := mock_userhist.NewMockUserHistoryStoreMethod(mockCtrl)
	pStore := mock_partner.NewMockPartnerCacheStoreMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		request PartnerServiceRequest
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		want     PartnerServiceInfo
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pStore.EXPECT().GetViewedUserCounter("1").Return("1", nil)
				pStore.EXPECT().GetCurentPartnerState("1").Return("4", nil)
				uStore.EXPECT().GetUserInfoByID(4).Return(models.User{
					Model: gorm.Model{
						ID: 4,
					},
					Username:   "U4",
					Fullname:   "F4",
					Email:      "E4",
					IsVerified: true,
				}, nil)

				hStore.EXPECT().CountByUserIDAndPartnerID(1, 4).Return(0, nil)
			},
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			want: PartnerServiceInfo{
				PartnerID:   4,
				Fullname:    "F4",
				IsVerified:  true,
				Status:      "PENDING",
				CreatedDate: "0001-01-01 00:00:00 +0000 UTC",
			},
			wantErr: false,
		},
		{
			name: "error get profile partner",
			mockFunc: func() {
				pStore.EXPECT().GetViewedUserCounter("1").Return("1", nil)
				pStore.EXPECT().GetCurentPartnerState("1").Return("4", nil)
				uStore.EXPECT().GetUserInfoByID(4).Return(models.User{}, fmt.Errorf("some error"))
			},
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			want:    PartnerServiceInfo{},
			wantErr: true,
		},
		{
			name: "error get state partner",
			mockFunc: func() {
				pStore.EXPECT().GetViewedUserCounter("1").Return("1", nil)
				pStore.EXPECT().GetCurentPartnerState("1").Return("", fmt.Errorf("some error"))
			},
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			want:    PartnerServiceInfo{},
			wantErr: true,
		},
		{
			name: "error get counter partner",
			mockFunc: func() {
				pStore.EXPECT().GetViewedUserCounter("1").Return("", fmt.Errorf("some error"))
			},
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			want:    PartnerServiceInfo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPartnerService(uStore, hStore, pStore, 10)
			tt.mockFunc()
			got, err := s.GetCurrentPartner(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("PartnerService.GetCurrentPartner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PartnerService.GetCurrentPartner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPartnerService_LikePartner(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	uStore := mock_user.NewMockUserStoreMethod(mockCtrl)
	hStore := mock_userhist.NewMockUserHistoryStoreMethod(mockCtrl)
	pStore := mock_partner.NewMockPartnerCacheStoreMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		request PartnerServiceRequest
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pStore.EXPECT().GetCurentPartnerState("1").Return("4", nil)
				hStore.EXPECT().CountByUserIDAndPartnerID(1, 4).Return(0, nil)
				hStore.EXPECT().CountByUserIDAndPartnerID(4, 1).Return(1, nil)
				uStore.EXPECT().GetUserInfoByID(4).Return(models.User{
					Model: gorm.Model{
						ID: 4,
					},
					Fullname: "P4",
				}, nil)

				hStore.EXPECT().CreateUserHistory(models.UserMatchHistory{
					UserID:      1,
					PartnerID:   4,
					PartnerName: "P4",
					Status:      models.MatchStatusApproved,
				}).Return(nil)

				hStore.EXPECT().UpdatePartnerStatus(models.UserMatchHistory{
					UserID:    uint(4),
					PartnerID: uint(1),
					Status:    models.MatchStatusApproved,
				}).Return(nil)
			},
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			wantErr: false,
		},
		{
			name: "error on get detail",
			mockFunc: func() {
				pStore.EXPECT().GetCurentPartnerState("1").Return("4", nil)
				hStore.EXPECT().CountByUserIDAndPartnerID(1, 4).Return(0, nil)
				hStore.EXPECT().CountByUserIDAndPartnerID(4, 1).Return(1, nil)
				uStore.EXPECT().GetUserInfoByID(4).Return(models.User{
					Model: gorm.Model{
						ID: 4,
					},
					Fullname: "P4",
				}, fmt.Errorf("some error"))
			},
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			wantErr: true,
		},
		{
			name: "error on user already like",
			mockFunc: func() {
				pStore.EXPECT().GetCurentPartnerState("1").Return("4", nil)
				hStore.EXPECT().CountByUserIDAndPartnerID(1, 4).Return(1, nil)
			},
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			wantErr: true,
		},
		{
			name: "error on check partner detail",
			mockFunc: func() {
				pStore.EXPECT().GetCurentPartnerState("1").Return("4", nil)
				hStore.EXPECT().CountByUserIDAndPartnerID(1, 4).Return(0, nil)
				hStore.EXPECT().CountByUserIDAndPartnerID(4, 1).Return(1, fmt.Errorf("some error"))
			},
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			wantErr: true,
		},
		{
			name: "error on check user detail",
			mockFunc: func() {
				pStore.EXPECT().GetCurentPartnerState("1").Return("4", nil)
				hStore.EXPECT().CountByUserIDAndPartnerID(1, 4).Return(0, fmt.Errorf("some error"))
			},
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			wantErr: true,
		},
		{
			name: "error on check user detail",
			mockFunc: func() {
				pStore.EXPECT().GetCurentPartnerState("1").Return("4", fmt.Errorf("some error"))
			},
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: false,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPartnerService(uStore, hStore, pStore, 10)
			tt.mockFunc()
			if err := s.LikePartner(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("PartnerService.LikePartner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPartnerService_GetListLikedPartner(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	uStore := mock_user.NewMockUserStoreMethod(mockCtrl)
	hStore := mock_userhist.NewMockUserHistoryStoreMethod(mockCtrl)
	pStore := mock_partner.NewMockPartnerCacheStoreMethod(mockCtrl)
	type args struct {
		request PartnerServiceRequest
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		want     []PartnerServiceInfo
		wantErr  bool
	}{
		{
			name: "success flow",
			mockFunc: func() {
				hStore.EXPECT().GetUserHistoryListByUserID(models.UserMatchHistory{
					UserID: 1,
				}).Return([]models.UserMatchHistory{
					{
						Model: gorm.Model{
							ID: 1,
						},
						UserID:      1,
						PartnerID:   4,
						PartnerName: "P4",
						Status:      models.MatchStatusApproved,
					},
				}, nil)
			},
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: true,
				},
			},
			want: []PartnerServiceInfo{
				{
					PartnerID:   4,
					Fullname:    "P4",
					Status:      models.MatchStatusApproved.String(),
					CreatedDate: "0001-01-01 00:00:00 +0000 UTC",
				},
			},
			wantErr: false,
		},
		{
			name: "error flow",
			mockFunc: func() {
				hStore.EXPECT().GetUserHistoryListByUserID(models.UserMatchHistory{
					UserID: 1,
				}).Return([]models.UserMatchHistory{}, fmt.Errorf("some error"))
			},
			args: args{
				request: PartnerServiceRequest{
					UserID:     1,
					IsVerified: true,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPartnerService(uStore, hStore, pStore, 10)
			tt.mockFunc()
			got, err := s.GetListLikedPartner(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("PartnerService.GetListLikedPartner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PartnerService.GetListLikedPartner() = %v, want %v", got, tt.want)
			}
		})
	}
}
