package partnercache

import (
	"fmt"
	"gilsaputro/dating-apps/pkg/redis"
	mock_redis "gilsaputro/dating-apps/pkg/redis/mock"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestNewPartnerCacheStore(t *testing.T) {
	type args struct {
		rd redis.RedisMethod
	}
	tests := []struct {
		name string
		args args
		want PartnerCacheStoreMethod
	}{
		{
			name: "success flow",
			args: args{
				rd: &redis.RedisClient{},
			},
			want: &PartnerCacheStore{
				rd: &redis.RedisClient{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPartnerCacheStore(tt.args.rd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPartnerCacheStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPartnerCacheStore_SetCurentPartnerState(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	rd := mock_redis.NewMockRedisMethod(mockCtrl)
	type args struct {
		userID    int
		partnerID int
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
				rd.EXPECT().Set("CPS:1", 1, 24*time.Hour).Return(nil)
			},
			args: args{
				userID:    1,
				partnerID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PartnerCacheStore{
				rd: rd,
			}
			tt.mockFunc()
			if err := s.SetCurentPartnerState(tt.args.userID, tt.args.partnerID); (err != nil) != tt.wantErr {
				t.Errorf("PartnerCacheStore.SetCurentPartnerState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPartnerCacheStore_GetCurentPartnerState(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	rd := mock_redis.NewMockRedisMethod(mockCtrl)
	type args struct {
		userID string
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
				rd.EXPECT().Get("CPS:1").Return("1", nil)
			},
			args: args{
				userID: "1",
			},
			want:    "1",
			wantErr: false,
		},
		{
			name: "nil data flow",
			mockFunc: func() {
				rd.EXPECT().Get("CPS:1").Return("", fmt.Errorf("redis: nil"))
			},
			args: args{
				userID: "1",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PartnerCacheStore{
				rd: rd,
			}
			tt.mockFunc()
			got, err := s.GetCurentPartnerState(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PartnerCacheStore.GetCurentPartnerState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PartnerCacheStore.GetCurentPartnerState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPartnerCacheStore_SetViewedPartnerHistory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	rd := mock_redis.NewMockRedisMethod(mockCtrl)
	type args struct {
		userID string
		value  string
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
				rd.EXPECT().Set("VPH:1", "1,2,3,4", 24*time.Hour).Return(nil)
			},
			args: args{
				userID: "1",
				value:  "1,2,3,4",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PartnerCacheStore{
				rd: rd,
			}
			tt.mockFunc()
			if err := s.SetViewedPartnerHistory(tt.args.userID, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("PartnerCacheStore.SetViewedPartnerHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPartnerCacheStore_GetViewedPartnerHistory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	rd := mock_redis.NewMockRedisMethod(mockCtrl)
	type args struct {
		userID string
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
				rd.EXPECT().Get("VPH:1").Return("1,2,3", nil)
			},
			args: args{
				userID: "1",
			},
			want:    "1,2,3",
			wantErr: false,
		},
		{
			name: "nil data flow",
			mockFunc: func() {
				rd.EXPECT().Get("VPH:1").Return("", fmt.Errorf("redis: nil"))
			},
			args: args{
				userID: "1",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PartnerCacheStore{
				rd: rd,
			}
			tt.mockFunc()
			got, err := s.GetViewedPartnerHistory(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PartnerCacheStore.GetViewedPartnerHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PartnerCacheStore.GetViewedPartnerHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPartnerCacheStore_SetViewedUserCounter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	rd := mock_redis.NewMockRedisMethod(mockCtrl)
	type args struct {
		userID string
		value  string
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
				rd.EXPECT().Set(gomock.Any(), "10", 24*time.Hour).Return(nil)
			},
			args: args{
				userID: "1",
				value:  "10",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PartnerCacheStore{
				rd: rd,
			}
			tt.mockFunc()
			if err := s.SetViewedUserCounter(tt.args.userID, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("PartnerCacheStore.SetViewedUserCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPartnerCacheStore_GetViewedUserCounter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	rd := mock_redis.NewMockRedisMethod(mockCtrl)
	type args struct {
		userID string
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
				rd.EXPECT().Get(gomock.Any()).Return("10", nil)
			},
			args: args{
				userID: "1",
			},
			want:    "10",
			wantErr: false,
		},
		{
			name: "nil data flow",
			mockFunc: func() {
				rd.EXPECT().Get(gomock.Any()).Return("", fmt.Errorf("redis: nil"))
			},
			args: args{
				userID: "1",
			},
			want:    "0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PartnerCacheStore{
				rd: rd,
			}
			tt.mockFunc()
			got, err := s.GetViewedUserCounter(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PartnerCacheStore.GetViewedUserCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PartnerCacheStore.GetViewedUserCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}
