package partner

import (
	"fmt"
	"gilsaputro/dating-apps/internal/store/partnercache"
	"gilsaputro/dating-apps/internal/store/user"
	"gilsaputro/dating-apps/internal/store/userhistory"
	"gilsaputro/dating-apps/models"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// PartnerServiceMethod is list method for Partner Service
type PartnerServiceMethod interface {
	LikePartner(request PartnerServiceRequest) error
	PassPartner(request PartnerServiceRequest) (PartnerServiceInfo, error)
	GetCurrentPartner(request PartnerServiceRequest) (PartnerServiceInfo, error)
	GetListLikedPartner(request PartnerServiceRequest) ([]PartnerServiceInfo, error)
}

// PartnerService is list dependencies for Partner service
type PartnerService struct {
	storeUser  user.UserStoreMethod
	storeHist  userhistory.UserHistoryStoreMethod
	cache      partnercache.PartnerCacheStoreMethod
	maxCounter int
}

// NewPartnerService is func to generate PartnerServiceMethod interface
func NewPartnerService(storeUser user.UserStoreMethod, storeHist userhistory.UserHistoryStoreMethod, cache partnercache.PartnerCacheStoreMethod, maxCounter int) PartnerServiceMethod {
	if maxCounter <= 0 {
		maxCounter = 10
	}
	return &PartnerService{
		storeHist:  storeHist,
		storeUser:  storeUser,
		cache:      cache,
		maxCounter: maxCounter,
	}
}

func (f PartnerService) PassPartner(request PartnerServiceRequest) (PartnerServiceInfo, error) {
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

	newPartnerID, err := f.generateNewPartner(request)
	if err != nil {
		return PartnerServiceInfo{}, err
	}

	PartnerInfo, err := f.storeUser.GetUserInfoByID(newPartnerID)
	if err != nil {
		return PartnerServiceInfo{}, err
	}

	status := f.getPartnerStatus(request.UserID, int(PartnerInfo.ID))

	// Set NewPartnerID to History and Add Counter
	if !request.IsVerified && status == "PENDING" {
		numCounter++
		f.cache.SetViewedUserCounter(userID, fmt.Sprintf("%d", numCounter))
	}

	return PartnerServiceInfo{
		PartnerID:   int(PartnerInfo.ID),
		Fullname:    PartnerInfo.Fullname,
		IsVerified:  PartnerInfo.IsVerified,
		Status:      status,
		CreatedDate: PartnerInfo.CreatedAt.String(),
	}, nil
}

func (f PartnerService) generateNewPartner(request PartnerServiceRequest) (int, error) {
	userID := fmt.Sprintf("%v", request.UserID)
	// Get Total User
	totalUser, err := f.storeUser.Count()
	if err != nil {
		return 0, err
	}

	// Get User Partner History to de-duplicate generate same partner
	partnerHistory, err := f.cache.GetViewedPartnerHistory(userID)
	if err != nil {
		return 0, err
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
		return 0, err
	}
	err = f.cache.SetCurentPartnerState(request.UserID, newPartnerID)
	if err != nil {
		return 0, err
	}

	return newPartnerID, err
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

func (f PartnerService) getPartnerStatus(userID int, partnerID int) string {
	var status = "PENDING"
	count, err := f.storeHist.CountByUserIDAndPartnerID(userID, partnerID)
	if err != nil {
		return status
	}

	if count > 0 {
		status = "LIKED"
	}
	return status
}

func (f PartnerService) GetCurrentPartner(request PartnerServiceRequest) (PartnerServiceInfo, error) {
	userID := fmt.Sprintf("%v", request.UserID)

	if !request.IsVerified {
		c, err := f.cache.GetViewedUserCounter(userID)
		if err != nil {
			return PartnerServiceInfo{}, err
		}

		count, err := strconv.Atoi(c)
		if err != nil {
			count = 0
		}

		if count >= f.maxCounter {
			return PartnerServiceInfo{}, ErrReachedMaxSwipeQuota
		}
	}

	newPartnerID, err := f.cache.GetCurentPartnerState(userID)
	if err != nil {
		return PartnerServiceInfo{}, err
	}

	partnerID, err := strconv.Atoi(newPartnerID)
	if err != nil {
		partnerID = 0
	}

	// generate if the current is not state (should be for first time user)
	if partnerID <= 0 {
		partnerID, err = f.generateNewPartner(request)
		if err != nil {
			return PartnerServiceInfo{}, err
		}
	}

	PartnerInfo, err := f.storeUser.GetUserInfoByID(partnerID)
	if err != nil {
		return PartnerServiceInfo{}, err
	}

	status := f.getPartnerStatus(request.UserID, int(PartnerInfo.ID))

	return PartnerServiceInfo{
		PartnerID:   int(PartnerInfo.ID),
		Fullname:    PartnerInfo.Fullname,
		IsVerified:  PartnerInfo.IsVerified,
		Status:      status,
		CreatedDate: PartnerInfo.CreatedAt.String(),
	}, nil
}

func (f PartnerService) LikePartner(request PartnerServiceRequest) error {
	userID := fmt.Sprintf("%v", request.UserID)
	partnerID, err := f.cache.GetCurentPartnerState(userID)
	if err != nil {
		return err
	}
	intPartnerID, err := strconv.Atoi(partnerID)
	if err != nil {
		intPartnerID = 0
	}
	// if partner id is not set return error
	if intPartnerID <= 0 {
		return ErrCurrentPartnerIsMissing
	}

	count, err := f.storeHist.CountByUserIDAndPartnerID(request.UserID, intPartnerID)
	if err != nil {
		return err
	}

	// user already like the partner
	if count > 0 {
		return ErrUserAlreadyLikePartner
	}

	// check if the partner like user
	count, err = f.storeHist.CountByUserIDAndPartnerID(intPartnerID, request.UserID)
	if err != nil {
		return err
	}

	status := models.MatchStatusPending
	if count > 0 {
		status = models.MatchStatusApproved
	}

	partnerInfo, err := f.storeUser.GetUserInfoByID(intPartnerID)
	if err != nil {
		return err
	}

	f.storeHist.CreateUserHistory(models.UserMatchHistory{
		UserID:      uint(request.UserID),
		PartnerID:   uint(partnerInfo.ID),
		PartnerName: partnerInfo.Fullname,
		Status:      status,
	})

	// update partner status if the partner already like the user
	if count > 0 {
		f.storeHist.UpdatePartnerStatus(models.UserMatchHistory{
			UserID:    uint(intPartnerID),
			PartnerID: uint(request.UserID),
			Status:    status,
		})
	}

	return nil
}

func (f PartnerService) GetListLikedPartner(request PartnerServiceRequest) ([]PartnerServiceInfo, error) {
	hist, err := f.storeHist.GetUserHistoryListByUserID(models.UserMatchHistory{
		UserID: uint(request.UserID),
	})

	if err != nil {
		return nil, err
	}

	var result []PartnerServiceInfo
	for _, data := range hist {
		result = append(result, PartnerServiceInfo{
			PartnerID:   int(data.PartnerID),
			Fullname:    data.PartnerName,
			Status:      data.Status.String(),
			CreatedDate: data.CreatedAt.String(),
		})
	}

	return result, nil
}
