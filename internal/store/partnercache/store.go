package partnercache

import (
	"fmt"
	"gilsaputro/dating-apps/pkg/redis"
	"strings"
	"time"
)

// PartnerCacheStoreMethod is set of methods for interacting with a partner cache storage system
type PartnerCacheStoreMethod interface {
	SetCurentPartnerState(userID, partnerID int) error
	GetCurentPartnerState(userID string) (string, error)
	SetViewedPartnerHistory(userID, value string) error
	GetViewedPartnerHistory(userID string) (string, error)
	SetViewedUserCounter(userID, value string) error
	GetViewedUserCounter(userID string) (string, error)
}

// PartnerCacheStore is list dependencies partner cache store
type PartnerCacheStore struct {
	rd redis.RedisMethod
}

// NewPartnerCacheStore is func to generate PartnerCacheStoreMethod interface
func NewPartnerCacheStore(rd redis.RedisMethod) PartnerCacheStoreMethod {
	return &PartnerCacheStore{
		rd: rd,
	}
}

const currentpartnerState string = `CPS:%v` // format CPS:<userid>

// SetCurentPartnerState is func to store current partner state of user id
func (f *PartnerCacheStore) SetCurentPartnerState(userID, partnerID int) error {
	key := fmt.Sprintf(currentpartnerState, userID)
	f.rd.Set(key, partnerID, 24*time.Hour)
	return nil
}

// GetCurentPartnerState is func to store current partner state of user id
func (f *PartnerCacheStore) GetCurentPartnerState(userID string) (string, error) {
	key := fmt.Sprintf(currentpartnerState, userID)
	c, err := f.rd.Get(key)
	if err != nil && strings.Contains(err.Error(), "redis: nil") {
		return "", nil
	}
	return c, err
}

const viewedPartnerHistory string = `VPH:%v` // format VPH:<userid>

// SetViewedPartnerHistory is func to store viewed partner history of user id
func (f *PartnerCacheStore) SetViewedPartnerHistory(userID, value string) error {
	key := fmt.Sprintf(viewedPartnerHistory, userID)
	f.rd.Set(key, value, 24*time.Hour)
	return nil
}

// GetViewedPartnerHistory is func to get viewed partner history of user id
func (f *PartnerCacheStore) GetViewedPartnerHistory(userID string) (string, error) {
	key := fmt.Sprintf(viewedPartnerHistory, userID)
	c, err := f.rd.Get(key)
	if err != nil && strings.Contains(err.Error(), "redis: nil") {
		return "", nil
	}

	return c, err
}

const viewedUserCounter string = `VUC:%v:%v` // format VUC:<datetime>:<userid>
const dateFormat string = "20060102"         // YYYYMMDD format

// SetViewedUserCounter is func to store daily viewed user counter of user id
func (f *PartnerCacheStore) SetViewedUserCounter(userID, value string) error {
	currentTime := time.Now().Format(dateFormat)
	key := fmt.Sprintf(viewedUserCounter, currentTime, userID)
	f.rd.Set(key, value, 24*time.Hour)
	return nil
}

// GetViewedUserCounter is func to get daily viewed user counter of user id
func (f *PartnerCacheStore) GetViewedUserCounter(userID string) (string, error) {
	currentTime := time.Now().Format(dateFormat)
	key := fmt.Sprintf(viewedUserCounter, currentTime, userID)
	c, err := f.rd.Get(key)
	if err != nil && strings.Contains(err.Error(), "redis: nil") {
		return "0", nil
	}

	return c, err
}
