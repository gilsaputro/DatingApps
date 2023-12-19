package userhistory

import (
	"errors"
	"gilsaputro/dating-apps/models"
	"gilsaputro/dating-apps/pkg/postgres"

	"github.com/jinzhu/gorm"
)

// UserHistoryStoreMethod is set of methods for interacting with a user storage system
type UserHistoryStoreMethod interface {
	CreateUserHistory(hist models.UserMatchHistory) error
	GetUserHistoryListByUserID(hist models.UserMatchHistory) ([]models.UserMatchHistory, error)
	CountByUserIDAndPartnerID(userID, partnerID int) (int, error)
	UpdatePartnerStatus(history models.UserMatchHistory) error
}

// UserHistoryStore is list dependencies user store
type UserHistoryStore struct {
	pg postgres.PostgresMethod
}

// NewUserStore is func to generate UserStoreMethod interface
func NewUserHistoryStore(pg postgres.PostgresMethod) UserHistoryStoreMethod {
	return &UserHistoryStore{
		pg: pg,
	}
}

func (u *UserHistoryStore) getDB() (*gorm.DB, error) {
	db := u.pg.GetDB()
	if db == nil {
		return nil, errors.New("Database Client is not init")
	}

	return db, nil
}

func (u UserHistoryStore) CreateUserHistory(history models.UserMatchHistory) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	return db.Create(&history).Error
}

func (u UserHistoryStore) GetUserHistoryDetailByID(history models.UserMatchHistory) (models.UserMatchHistory, error) {
	db, err := u.getDB()
	if err != nil {
		return models.UserMatchHistory{}, err
	}

	result := models.UserMatchHistory{}
	err = db.Model(models.UserMatchHistory{}).First(&result, history).Error
	if err != nil {
		return models.UserMatchHistory{}, err
	}

	return result, err
}

func (u UserHistoryStore) GetUserHistoryListByUserID(history models.UserMatchHistory) ([]models.UserMatchHistory, error) {
	db, err := u.getDB()
	if err != nil {
		return nil, err
	}

	result := []models.UserMatchHistory{}
	err = db.Model(models.UserMatchHistory{}).Find(&result, history).Error
	if err != nil {
		return nil, err
	}

	return result, err
}

func (u UserHistoryStore) CountByUserIDAndPartnerID(userID, partnerID int) (int, error) {
	db, err := u.getDB()
	if err != nil {
		return 0, err
	}

	var count int
	err = db.Model(models.UserMatchHistory{}).Where("user_id = ? AND partner_id = ?", userID, partnerID).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, err
}

func (u UserHistoryStore) UpdatePartnerStatus(history models.UserMatchHistory) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	var user models.UserMatchHistory
	err = db.Where("user_id = ? AND partner_id = ?", history.UserID, history.PartnerID).First(&user).Error
	if err != nil {
		return err
	}

	user.Status = history.Status
	return db.Save(&user).Error
}
