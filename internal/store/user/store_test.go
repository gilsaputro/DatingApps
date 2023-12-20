package user

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

func TestNewUserStore(t *testing.T) {
	type args struct {
		pg postgres.PostgresMethod
	}
	tests := []struct {
		name string
		args args
		want UserStoreMethod
	}{
		{
			name: "success flow",
			args: args{
				pg: &postgres.Client{},
			},
			want: &UserStore{
				pg: &postgres.Client{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserStore(tt.args.pg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserStore() = %v, want %v", got, tt.want)
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

func TestUserStore_CreateUser(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func()
		args     models.User
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","username","password","fullname","email","is_verified") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "users"."id"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mockDB.ExpectCommit()
			},
			args: models.User{
				Username:   "abc",
				Password:   "password_hashed",
				Fullname:   "abc_a",
				Email:      "abc@dev.com",
				IsVerified: false,
			},
			wantErr: false,
		},
		{
			name: "failed insert",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","username","password","fullname","email","is_verified") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "users"."id"`)).WillReturnError(fmt.Errorf("some error"))
				mockDB.ExpectCommit()
			},
			args: models.User{
				Username:   "abc",
				Password:   "password_hashed",
				Fullname:   "abc_a",
				Email:      "abc@dev.com",
				IsVerified: false,
			},
			wantErr: true,
		},
		{
			name: "nil database",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			args: models.User{
				Username:   "abc",
				Password:   "password_hashed",
				Fullname:   "abc_a",
				Email:      "abc@dev.com",
				IsVerified: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserStore{
				pg: pg,
			}
			tt.mockFunc()
			if err := service.CreateUser(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("UserStore.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserStore_UpdateUser(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	var userDataMock = &models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Username: "abc",
		Password: "pas1",
		Fullname: "full",
		Email:    "email",
	}
	var expectedRows = sqlmock.NewRows([]string{"id", "username", "password", "fullname", "email"}).
		AddRow(userDataMock.ID, userDataMock.Username, userDataMock.Password, userDataMock.Fullname, userDataMock.Email)

	tests := []struct {
		name     string
		mockFunc func()
		args     models.User
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND ((username = $1 AND id = $2)) ORDER BY "users"."id" ASC LIMIT 1`)).WillReturnRows(expectedRows)
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "updated_at" = $1, "deleted_at" = $2, "username" = $3, "password" = $4, "fullname" = $5, "email" = $6, "is_verified" = $7 WHERE "users"."deleted_at" IS NULL AND "users"."id" = $8`)).WillReturnResult(sqlmock.NewResult(1, 1))
				mockDB.ExpectCommit()
			},
			args: models.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username:   "abc",
				Password:   "password_hashed",
				Fullname:   "abc_a",
				Email:      "abc@dev.com",
				IsVerified: true,
			},
			wantErr: false,
		},
		{
			name: "failed update",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND ((username = $1 AND id = $2)) ORDER BY "users"."id" ASC LIMIT 1`)).WillReturnRows(expectedRows)
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "updated_at" = $1, "deleted_at" = $2, "username" = $3, "password" = $4, "fullname" = $5, "email" = $6, "is_verified" = $7 WHERE "users"."deleted_at" IS NULL AND "users"."id" = $8`)).WillReturnError(fmt.Errorf("some error"))
			},
			args: models.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username: "abc",
				Password: "password_hashed",
				Fullname: "abc_a",
				Email:    "abc@dev.com",
			},
			wantErr: true,
		},
		{
			name: "failed get data",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND ((username = $1 AND id = $2)) ORDER BY "users"."id" ASC LIMIT 1`)).WillReturnError(fmt.Errorf("some error"))
			},
			args: models.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username: "abc",
				Password: "password_hashed",
				Fullname: "abc_a",
				Email:    "abc@dev.com",
			},
			wantErr: true,
		},
		{
			name: "nil database",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			args: models.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username: "abc",
				Password: "password_hashed",
				Fullname: "abc_a",
				Email:    "abc@dev.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserStore{
				pg: pg,
			}
			tt.mockFunc()
			if err := service.UpdateUser(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("UserStore.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserStore_GetUserInfoByUsername(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	var userDataMock = &models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Username: "abc",
		Password: "pas1",
		Fullname: "full",
		Email:    "email",
	}
	var expectedRows = sqlmock.NewRows([]string{"id", "username", "password", "fullname", "email"}).
		AddRow(userDataMock.ID, userDataMock.Username, userDataMock.Password, userDataMock.Fullname, userDataMock.Email)

	tests := []struct {
		name     string
		username string
		mockFunc func()
		want     models.User
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND ((username = $1)) ORDER BY "users"."id" ASC LIMIT 1`)).WillReturnRows(expectedRows)
			},
			username: "abc",
			want: models.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username: "abc",
				Password: "pas1",
				Fullname: "full",
				Email:    "email",
			},
			wantErr: false,
		},
		{
			name: "failed get data",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND ((username = $1)) ORDER BY "users"."id" ASC LIMIT 1`)).WillReturnError(fmt.Errorf("some error"))
			},
			username: "abc",
			want:     models.User{},
			wantErr:  true,
		},
		{
			name: "nil database",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			username: "abc",
			want:     models.User{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserStore{
				pg: pg,
			}
			tt.mockFunc()
			got, err := service.GetUserInfoByUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.GetUserInfoByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserStore.GetUserInfoByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserStore_DeleteUser(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func()
		userid   int
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=$1 WHERE "users"."deleted_at" IS NULL AND "users"."id" = $2`)).WillReturnResult(sqlmock.NewResult(1, 1))
				mockDB.ExpectCommit()
			},
			userid:  1,
			wantErr: false,
		},
		{
			name: "error delete",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=$1 WHERE "users"."deleted_at" IS NULL AND "users"."id" = $2`)).WillReturnError(fmt.Errorf("some error"))
			},
			userid:  1,
			wantErr: true,
		},
		{
			name: "nil database",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			userid:  1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserStore{
				pg: pg,
			}
			tt.mockFunc()
			if err := service.DeleteUser(tt.userid); (err != nil) != tt.wantErr {
				t.Errorf("UserStore.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserStore_GetUserInfoByID(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	var userDataMock = &models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Username: "abc",
		Password: "pas1",
		Fullname: "full",
		Email:    "email",
	}
	var expectedRows = sqlmock.NewRows([]string{"id", "username", "password", "fullname", "email"}).
		AddRow(userDataMock.ID, userDataMock.Username, userDataMock.Password, userDataMock.Fullname, userDataMock.Email)

	tests := []struct {
		name     string
		userid   int
		mockFunc func()
		want     models.User
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND (("users"."id" = 1)) ORDER BY "users"."id" ASC LIMIT 1`)).WillReturnRows(expectedRows)
			},
			userid: 1,
			want: models.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username: "abc",
				Password: "pas1",
				Fullname: "full",
				Email:    "email",
			},
			wantErr: false,
		},
		{
			name: "error get data",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND (("users"."id" = 1)) ORDER BY "users"."id" ASC LIMIT 1`)).WillReturnError(fmt.Errorf("some error"))
			},
			userid:  1,
			want:    models.User{},
			wantErr: true,
		},
		{
			name: "nil database",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			want:    models.User{},
			userid:  1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserStore{
				pg: pg,
			}
			tt.mockFunc()
			got, err := service.GetUserInfoByID(tt.userid)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.GetUserInfoByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserStore.GetUserInfoByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserStore_Count(t *testing.T) {
	db, mockDB, gormDB := InitDBsMockupStat()
	defer db.Close()
	defer gormDB.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pg := mock_postgres.NewMockPostgresMethod(mockCtrl)
	tests := []struct {
		name     string
		mockFunc func()
		want     int
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users"  WHERE "users"."deleted_at" IS NULL`)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
			},
			want:    10,
			wantErr: false,
		},
		{
			name: "error get data",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(gormDB)
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users"  WHERE "users"."deleted_at" IS NULL`)).WillReturnError(fmt.Errorf("some error"))
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "nil database",
			mockFunc: func() {
				pg.EXPECT().GetDB().Return(nil)
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := UserStore{
				pg: pg,
			}
			tt.mockFunc()
			got, err := service.Count()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserStore.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}
