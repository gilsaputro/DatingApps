package find

import "errors"

var (
	ErrReachedMaxSwipeQuota = errors.New("the user already reach max quota for swipe")
)

// FindPartnerServiceRequest is list parameter for Find Partner
type FindPartnerServiceRequest struct {
	UserID     int
	IsVerified bool
}

// PartnerServiceInfo struct is list parameter info for partner sevice
type PartnerServiceInfo struct {
	PartnerID   int
	Fullname    string
	CreatedDate string
}
