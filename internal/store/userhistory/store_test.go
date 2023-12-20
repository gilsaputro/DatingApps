package userhistory

import (
	"database/sql"
	"fmt"
	"gilsaputro/dating-apps/models"
	"gilsaputro/dating-apps/pkg/postgres"
	mock_postgres "gilsaputro/dating-apps/pkg/postgres/mock"
	"log"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewUserHistoryStore(t *testing.T) {
	type args struct {
		pg postgres.PostgresMethod
	}
	tests := []struct {
		name string
		args args
		want UserHistoryStoreMethod
	}{
		{
			name: "success flow",
			args: args{
				pg: &postgres.Client{},
			},
			want: &UserHistoryStore{
				pg: &postgres.Client{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserHistoryStore(tt.args.pg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserHistoryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func InitDBsMockupStat() (*sql.DB, sqlmock.Sqlmock, *gorm.DB) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open("postgres", db)
	gormDB.LogMode(true)
	gormDB.SetLogger(log.New(os.Stdout, "\n", 0))
	gormDB.Debug()
	return db, mock, gormDB
}

func TestUserHistoryStore_CreateUserHistory(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	type args struct {
		history models.UserMatchHistory
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
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_match_histories" ("created_at","updated_at","deleted_at","user_id","partner_id","partner_name","status") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "user_match_histories"."id"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mockDB.ExpectCommit()
			},
			args: args{
				history: models.UserMatchHistory{
					UserID:      1,
					PartnerID:   1,
					PartnerName: "A",
					Status:      1,
				},
			},
			wantErr: false,
		},
		{
			name: "error on db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_match_histories" ("created_at","updated_at","deleted_at","user_id","partner_id","partner_name","status") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "user_match_histories"."id"`)).WillReturnError(fmt.Errorf("some error"))
			},
			args: args{
				history: models.UserMatchHistory{
					UserID:      1,
					PartnerID:   1,
					PartnerName: "A",
					Status:      1,
				},
			},
			wantErr: true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			args: args{
				history: models.UserMatchHistory{
					UserID:      1,
					PartnerID:   1,
					PartnerName: "A",
					Status:      1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserHistoryStore{
				pg: pg,
			}
			tt.mockFunc()
			if err := service.CreateUserHistory(tt.args.history); (err != nil) != tt.wantErr {
				t.Errorf("UserHistoryStore.CreateUserHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHistoryStore_GetUserHistoryListByUserID(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	var userDataMock = &models.UserMatchHistory{
		Model: gorm.Model{
			ID: 1,
		},
		UserID:      1,
		PartnerID:   1,
		PartnerName: "A",
		Status:      1,
	}
	var expectedRows = sqlmock.NewRows([]string{"id", "user_id", "partner_id", "partner_name", "status"}).
		AddRow(userDataMock.ID, userDataMock.UserID, userDataMock.PartnerID, userDataMock.PartnerName, userDataMock.Status)

	type args struct {
		history models.UserMatchHistory
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		want     []models.UserMatchHistory
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_match_histories" WHERE "user_match_histories"."deleted_at" IS NULL AND (("user_match_histories"."user_id" = $1))`)).WillReturnRows(expectedRows)
			},
			args: args{
				history: models.UserMatchHistory{
					UserID: 1,
				},
			},
			want: []models.UserMatchHistory{
				{
					Model: gorm.Model{
						ID: 1,
					},
					UserID:      1,
					PartnerID:   1,
					PartnerName: "A",
					Status:      1,
				},
			},
			wantErr: false,
		},
		{
			name: "error on db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_match_histories" WHERE "user_match_histories"."deleted_at" IS NULL AND (("user_match_histories"."user_id" = $1))`)).WillReturnError(fmt.Errorf("some error"))
			},
			args: args{
				history: models.UserMatchHistory{
					UserID: 1,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "db is nil",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			args: args{
				history: models.UserMatchHistory{
					UserID: 1,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserHistoryStore{
				pg: pg,
			}
			tt.mockFunc()
			got, err := service.GetUserHistoryListByUserID(tt.args.history)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserHistoryStore.GetUserHistoryListByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserHistoryStore.GetUserHistoryListByUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserHistoryStore_CountByUserIDAndPartnerID(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	type args struct {
		userID    int
		partnerID int
	}
	tests := []struct {
		name     string
		mockFunc func()
		args     args
		want     int
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "user_match_histories"  WHERE "user_match_histories"."deleted_at" IS NULL AND ((user_id = $1 AND partner_id = $2))`)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			args: args{
				userID:    1,
				partnerID: 1,
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "error on db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "user_match_histories"  WHERE "user_match_histories"."deleted_at" IS NULL AND ((user_id = $1 AND partner_id = $2))`)).WillReturnError(fmt.Errorf("some error"))
			},
			args: args{
				userID:    1,
				partnerID: 1,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			args: args{
				userID:    1,
				partnerID: 1,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserHistoryStore{
				pg: pg,
			}
			tt.mockFunc()
			got, err := service.CountByUserIDAndPartnerID(tt.args.userID, tt.args.partnerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserHistoryStore.CountByUserIDAndPartnerID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserHistoryStore.CountByUserIDAndPartnerID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserHistoryStore_UpdatePartnerStatus(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	var userDataMock = &models.UserMatchHistory{
		Model: gorm.Model{
			ID: 1,
		},
		UserID:      1,
		PartnerID:   1,
		PartnerName: "A",
		Status:      0,
	}
	var expectedRows = sqlmock.NewRows([]string{"id", "user_id", "partner_id", "partner_name", "status"}).
		AddRow(userDataMock.ID, userDataMock.UserID, userDataMock.PartnerID, userDataMock.PartnerName, userDataMock.Status)

	type args struct {
		history models.UserMatchHistory
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
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_match_histories" WHERE "user_match_histories"."deleted_at" IS NULL AND ((user_id = $1 AND partner_id = $2)) ORDER BY "user_match_histories"."id" ASC LIMIT 1`)).WillReturnRows(expectedRows)
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "user_match_histories" SET "updated_at" = $1, "deleted_at" = $2, "user_id" = $3, "partner_id" = $4, "partner_name" = $5, "status" = $6 WHERE "user_match_histories"."deleted_at" IS NULL AND "user_match_histories"."id" = $7`)).WillReturnResult(sqlmock.NewResult(1, 1))
				mockDB.ExpectCommit()
			},
			args: args{
				history: models.UserMatchHistory{
					Model: gorm.Model{
						ID: 1,
					},
					UserID:      1,
					PartnerID:   1,
					PartnerName: "A",
					Status:      1,
				},
			},
			wantErr: false,
		},
		{
			name: "error on db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_match_histories" WHERE "user_match_histories"."deleted_at" IS NULL AND ((user_id = $1 AND partner_id = $2)) ORDER BY "user_match_histories"."id" ASC LIMIT 1`)).WillReturnRows(expectedRows)
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "user_match_histories" SET "updated_at" = $1, "deleted_at" = $2, "user_id" = $3, "partner_id" = $4, "partner_name" = $5, "status" = $6 WHERE "user_match_histories"."deleted_at" IS NULL AND "user_match_histories"."id" = $7`)).WillReturnError(fmt.Errorf("some error"))
			},
			args: args{
				history: models.UserMatchHistory{
					Model: gorm.Model{
						ID: 1,
					},
					UserID:      1,
					PartnerID:   1,
					PartnerName: "A",
					Status:      1,
				},
			},
			wantErr: true,
		},
		{
			name: "error on db select",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_match_histories" WHERE "user_match_histories"."deleted_at" IS NULL AND ((user_id = $1 AND partner_id = $2)) ORDER BY "user_match_histories"."id" ASC LIMIT 1`)).WillReturnError(fmt.Errorf("some error"))
			},
			args: args{
				history: models.UserMatchHistory{
					Model: gorm.Model{
						ID: 1,
					},
					UserID:      1,
					PartnerID:   1,
					PartnerName: "A",
					Status:      1,
				},
			},
			wantErr: true,
		},
		{
			name: "nil db",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			args: args{
				history: models.UserMatchHistory{
					Model: gorm.Model{
						ID: 1,
					},
					UserID:      1,
					PartnerID:   1,
					PartnerName: "A",
					Status:      1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserHistoryStore{
				pg: pg,
			}
			tt.mockFunc()
			if err := service.UpdatePartnerStatus(tt.args.history); (err != nil) != tt.wantErr {
				t.Errorf("UserHistoryStore.UpdatePartnerStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
