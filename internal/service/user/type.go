package user

import "errors"

// list Service error
var (
	ErrNotGuest              = errors.New("register feature only available for guest")
	ErrUserNameNotExists     = errors.New("username is not exists")
	ErrUserNameAlreadyExists = errors.New("username already exists")
	ErrPasswordIsIncorrect   = errors.New("password is incorrect")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrCannotDeleteOtherUser = errors.New("cannot delete other user, please login first")
	ErrDataNotFound          = errors.New("data not found")
	ErrCannotUpdateOtherUser = errors.New("cannot edit other user, please login first")
	ErrCannotGetOtherUser    = errors.New("cannot get other user data")
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
