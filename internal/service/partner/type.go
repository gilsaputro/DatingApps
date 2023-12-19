package partner

import "errors"

var (
	ErrReachedMaxSwipeQuota    = errors.New("the user already reach max quota for swipe")
	ErrCurrentPartnerIsMissing = errors.New("the user partner is missing, please find one partner first")
	ErrUserAlreadyLikePartner  = errors.New("the user already like the partner")
)

// PartnerServiceRequest is list parameter for Partner Partner
type PartnerServiceRequest struct {
	UserID     int
	IsVerified bool
}

// PartnerServiceInfo struct is list parameter info for partner sevice
type PartnerServiceInfo struct {
	PartnerID   int
	Fullname    string
	Status      string
	CreatedDate string
}
