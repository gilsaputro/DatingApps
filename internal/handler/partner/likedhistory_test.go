package partner

import (
	"context"
	"fmt"
	"gilsaputro/dating-apps/internal/service/partner"
	"gilsaputro/dating-apps/internal/service/partner/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestPartnerHandler_LikedHistoryHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	m := mock.NewMockPartnerServiceMethod(mockCtrl)
	defer mockCtrl.Finish()
	type args struct {
		userID     int
		isVerified bool
		timeout    int
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
				userID:     1,
				isVerified: true,
				timeout:    5,
			},
			mockFunc: func() {
				m.EXPECT().GetListLikedPartner(partner.PartnerServiceRequest{
					UserID:     1,
					IsVerified: true,
				}).Return([]partner.PartnerServiceInfo{
					{
						PartnerID:  1,
						Fullname:   "full",
						IsVerified: true,
						Status:     "PENDING",
					},
				}, nil)
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 200,
				body: `{"data":[{"id":1,"fullname":"full","status":"full","is_verified":false,"created_date":""}],"code":200,"message":"success"}`,
			},
		},
		{
			name: "error on service flow",
			args: args{
				userID:     1,
				isVerified: true,
				timeout:    5,
			},
			mockFunc: func() {
				m.EXPECT().GetListLikedPartner(partner.PartnerServiceRequest{
					UserID:     1,
					IsVerified: true,
				}).Return([]partner.PartnerServiceInfo{}, fmt.Errorf("some error"))
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
			name: "error on service flow",
			args: args{
				userID:     1,
				isVerified: true,
				timeout:    5,
			},
			mockFunc: func() {
				m.EXPECT().GetListLikedPartner(partner.PartnerServiceRequest{
					UserID:     1,
					IsVerified: true,
				}).Return([]partner.PartnerServiceInfo{}, fmt.Errorf("some error"))
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
			name: "error on verified value flow",
			args: args{
				userID:  1,
				timeout: 5,
			},
			mockFunc: func() {
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 500,
				body: `{"code":500,"message":"Internal Server Error"}`,
			},
		},
		{
			name: "error on userid value flow",
			args: args{
				timeout: 5,
			},
			mockFunc: func() {
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 500,
				body: `{"code":500,"message":"Internal Server Error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			defer mockCtrl.Finish()
			handler := NewPartnerHandler(m, WithTimeoutOptions(tt.args.timeout))
			r := httptest.NewRequest(http.MethodGet, "/user", strings.NewReader(``))
			ctx, cancel := tt.mockContext()
			defer cancel()
			r = r.WithContext(ctx)
			if tt.args.userID > 0 {
				r = r.WithContext(context.WithValue(r.Context(), "id", tt.args.userID))
			}

			if tt.args.isVerified {
				r = r.WithContext(context.WithValue(r.Context(), "isverified", tt.args.isVerified))
			}
			w := httptest.NewRecorder()
			handler.LikedHistoryHandler(w, r)
			result := w.Result()
			resBody, err := ioutil.ReadAll(result.Body)

			if err != nil {
				t.Fatalf("Error read body err = %v\n", err)
			}

			if string(resBody) != tt.want.body {
				t.Fatalf("LikedHistoryHandler body got =%s, want %s \n", string(resBody), tt.want.body)
			}

			if result.StatusCode != tt.want.code {
				t.Fatalf("LikedHistoryHandler status code got =%d, want %d \n", result.StatusCode, tt.want.code)
			}
		})
	}
}
