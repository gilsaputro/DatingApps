package authentication

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

// LoginUserServiceRequest is list parameter for login user
type LoginServiceRequest struct {
	Username string
	Password string
}

// RegisterUserServiceRequest is list parameter for register user
type RegisterServiceRequest struct {
	TokenRequest string
	Username     string
	Password     string
	Fullname     string
	Email        string
}
