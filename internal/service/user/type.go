package user

import "errors"

// list Service error
var (
	ErrUserNameNotExists     = errors.New("username is not exists")
	ErrUserNameAlreadyExists = errors.New("username already exists")
	ErrPasswordIsIncorrect   = errors.New("password is incorrect")
	ErrUserIsVerified        = errors.New("user already verified")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrDataNotFound          = errors.New("data not found")
)

// UserServiceInfo struct is list parameter info for user sevice
type UserServiceInfo struct {
	UserId      int
	Username    string
	Fullname    string
	Email       string
	CreatedDate string
}

// DeleteUserServiceRequest is list parameter for add user by user
type DeleteUserServiceRequest struct {
	UserId   int
	Password string
}

// UpdateUserServiceRequest is list parameter for update user
type UpdateUserServiceRequest struct {
	UserId   int
	Username string
	Password string
	Fullname string
	Email    string
}

// GetByIDServiceRequest is list parameter for get user by id
type GetByIDServiceRequest struct {
	UserId int
}

// UpgradeServiceRequest is list parameter for upgrade user
type UpgradeServiceRequest struct {
	UserId   int
	Password string
}
