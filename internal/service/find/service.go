package find

import (
	"fmt"
	"gilsaputro/dating-apps/internal/store/findcache"
	"gilsaputro/dating-apps/internal/store/user"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// FindServiceMethod is list method for Find Service
type FindServiceMethod interface {
	FindPartner(request FindPartnerServiceRequest) (PartnerServiceInfo, error)
}

// FindService is list dependencies for Find service
type FindService struct {
	store      user.UserStoreMethod
	cache      findcache.FindCacheStoreMethod
	maxCounter int
}

// NewFindService is func to generate FindServiceMethod interface
func NewFindService(store user.UserStoreMethod, cache findcache.FindCacheStoreMethod, maxCounter int) FindServiceMethod {
	if maxCounter <= 0 {
		maxCounter = 10
	}
	return &FindService{
		store:      store,
		cache:      cache,
		maxCounter: maxCounter,
	}
}

func (f FindService) FindPartner(request FindPartnerServiceRequest) (PartnerServiceInfo, error) {
	userID := fmt.Sprintf("%v", request.UserID)
	var numCounter = 0
	// check max counter swipe if the user is not verified
	if !request.IsVerified {
		counter, err := f.cache.GetViewedUserCounter(userID)
		if err != nil {
			return PartnerServiceInfo{}, err
		}

		numCounter, _ = strconv.Atoi(counter)
		if numCounter >= f.maxCounter {
			return PartnerServiceInfo{}, ErrReachedMaxSwipeQuota
		}
	}

	// Get Total User
	totalUser, err := f.store.Count()
	if err != nil {
		return PartnerServiceInfo{}, err
	}

	// Get User Partner History to duplicate generate same partner
	partnerHistory, err := f.cache.GetViewedPartnerHistory(userID)
	if err != nil {
		return PartnerServiceInfo{}, err
	}

	var listPartnerHistoryInt []int
	listPartnerHistory := strings.Split(partnerHistory, ",")
	if len(listPartnerHistory) < totalUser {
		for _, str := range listPartnerHistory {
			// Konversi string ke integer
			num, _ := strconv.Atoi(str)
			// Menambahkan integer ke slice integer
			listPartnerHistoryInt = append(listPartnerHistoryInt, num)
		}
	}
	excludePartnerID := append(listPartnerHistoryInt, request.UserID)
	newPartnerID := generateRandomNumber(totalUser, excludePartnerID)

	PartnerInfo, err := f.store.GetUserInfoByID(newPartnerID)
	if err != nil {
		return PartnerServiceInfo{}, err
	}

	// Set NewPartnerID to History and Add Counter
	if !request.IsVerified {
		numCounter++
		f.cache.SetViewedUserCounter(userID, fmt.Sprintf("%d", numCounter))
	}

	if len(listPartnerHistoryInt) > 10 {
		listPartnerHistoryInt = listPartnerHistoryInt[1:]
	}
	listPartnerHistoryInt = append(listPartnerHistoryInt, newPartnerID)

	newPartnerHistory := make([]string, len(listPartnerHistoryInt))
	for i, v := range listPartnerHistoryInt {
		newPartnerHistory[i] = fmt.Sprint(v)
	}

	err = f.cache.SetViewedPartnerHistory(userID, strings.Join(newPartnerHistory, ","))
	if err != nil {
		return PartnerServiceInfo{}, err
	}
	err = f.cache.SetCurentPartnerState(request.UserID, newPartnerID)
	if err != nil {
		return PartnerServiceInfo{}, err
	}

	return PartnerServiceInfo{
		PartnerID:   PartnerInfo.UserId,
		Fullname:    PartnerInfo.Fullname,
		CreatedDate: PartnerInfo.CreatedDate,
	}, nil
}

func generateRandomNumber(max int, exclude []int) int {
	rand.Seed(time.Now().UnixNano())

	excludeMap := make(map[int]bool)
	for _, num := range exclude {
		excludeMap[num] = true
	}

	var randomNumber int
	for {
		randomNumber = rand.Intn(max) + 1
		if !excludeMap[randomNumber] {
			break
		}
	}

	return randomNumber
}
