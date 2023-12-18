package authentication

import (
	"gilsaputro/dating-apps/internal/store/user"
	"gilsaputro/dating-apps/pkg/hash"
	"gilsaputro/dating-apps/pkg/token"
	"strings"
)

// AuthenticationServiceMethod is list method for Authentication Service
type AuthenticationServiceMethod interface {
	Login(LoginServiceRequest) (string, error)
	Register(RegisterServiceRequest) error
}

// AuthenticationService is list dependencies for Authentication service
type AuthenticationService struct {
	store user.UserStoreMethod
	token token.TokenMethod
	hash  hash.HashMethod
}

// NewAuthenticationService is func to generate AuthenticationServiceMethod interface
func NewAuthenticationService(store user.UserStoreMethod, token token.TokenMethod, hash hash.HashMethod) AuthenticationServiceMethod {
	return &AuthenticationService{
		hash:  hash,
		token: token,
		store: store,
	}
}

// Login is service layer func to validate and generate token if the Authentication is exists
func (u *AuthenticationService) Login(request LoginServiceRequest) (string, error) {
	AuthenticationInfo, err := u.store.GetUserInfoByUsername(request.Username)
	if err != nil {
		return "", err
	}

	if AuthenticationInfo.UserId <= 0 {
		return "", ErrUserNameNotExists
	}

	if !u.hash.CompareValue(AuthenticationInfo.Password, request.Password) {
		return "", ErrPasswordIsIncorrect
	}

	return u.token.GenerateToken(token.TokenBody{
		UserID: AuthenticationInfo.UserId,
	})
}

// Register is service layer func to validate and creating Authentication to database if the Authentication is not exists
func (u *AuthenticationService) Register(request RegisterServiceRequest) error {
	AuthenticationInfo, err := u.store.GetUserInfoByUsername(request.Username)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return err
	}

	if AuthenticationInfo.UserId > 0 {
		return ErrUserNameAlreadyExists
	}

	hashPassword, err := u.hash.HashValue(request.Password)
	if err != nil {
		return err
	}

	return u.store.CreateUser(user.UserStoreInfo{
		Username: request.Username,
		Password: string(hashPassword),
		Fullname: request.Fullname,
		Email:    request.Email,
	})
}
