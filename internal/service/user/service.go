package user

import (
	"gilsaputro/dating-apps/internal/store/user"
	"gilsaputro/dating-apps/pkg/hash"
)

// UserServiceMethod is list method for User Service
type UserServiceMethod interface {
	DeleteUser(DeleteUserServiceRequest) error
	UpdateUser(UpdateUserServiceRequest) (UserServiceInfo, error)
	GetUserByID(GetByIDServiceRequest) (UserServiceInfo, error)
	UpgradeUser(UpgradeServiceRequest) error
}

// UserService is list dependencies for user service
type UserService struct {
	store user.UserStoreMethod
	hash  hash.HashMethod
}

// NewUserService is func to generate UserServiceMethod interface
func NewUserService(store user.UserStoreMethod, hash hash.HashMethod) UserServiceMethod {
	return &UserService{
		hash:  hash,
		store: store,
	}
}

// DeleteUser is service level func to validate and delete user info in database
func (u *UserService) DeleteUser(request DeleteUserServiceRequest) error {
	if request.UserId <= 0 {
		return ErrDataNotFound
	}

	userInfo, err := u.store.GetUserInfoByID(request.UserId)
	if err != nil || userInfo.ID <= 0 {
		return err
	}

	if !u.hash.CompareValue(userInfo.Password, request.Password) {
		return ErrPasswordIsIncorrect
	}

	return u.store.DeleteUser(int(userInfo.ID))
}

// UpdateUser is service level func to validate and update user info in database
func (u *UserService) UpdateUser(request UpdateUserServiceRequest) (UserServiceInfo, error) {
	if request.UserId <= 0 {
		return UserServiceInfo{}, ErrDataNotFound
	}

	userInfo, err := u.store.GetUserInfoByID(request.UserId)
	if err != nil || userInfo.ID <= 0 {
		return UserServiceInfo{}, err
	}

	if len(request.Password) > 0 {
		hashPassword, err := u.hash.HashValue(request.Password)
		if err != nil {
			return UserServiceInfo{}, err
		}

		userInfo.Password = string(hashPassword)
	}

	if len(request.Email) > 0 {
		userInfo.Email = request.Email
	}

	if len(request.Fullname) > 0 {
		userInfo.Fullname = request.Fullname
	}

	err = u.store.UpdateUser(userInfo)
	if err != nil {
		return UserServiceInfo{}, err
	}

	return UserServiceInfo{
		UserId:      int(userInfo.ID),
		Username:    userInfo.Username,
		Fullname:    userInfo.Fullname,
		Email:       userInfo.Email,
		CreatedDate: userInfo.CreatedAt.String(),
	}, nil
}

// GetUserByID is service level func to validate and get all user based id
func (u *UserService) GetUserByID(request GetByIDServiceRequest) (UserServiceInfo, error) {
	userInfo, err := u.store.GetUserInfoByID(int(request.UserId))
	if err != nil || userInfo.ID <= 0 {
		return UserServiceInfo{}, err
	}

	return UserServiceInfo{
		UserId:      int(userInfo.ID),
		Username:    userInfo.Username,
		Fullname:    userInfo.Fullname,
		Email:       userInfo.Email,
		IsVerified:  userInfo.IsVerified,
		CreatedDate: userInfo.CreatedAt.String(),
	}, nil
}

func (u *UserService) UpgradeUser(request UpgradeServiceRequest) error {
	if request.UserId <= 0 {
		return ErrDataNotFound
	}

	userInfo, err := u.store.GetUserInfoByID(request.UserId)
	if err != nil || userInfo.ID <= 0 {
		return err
	}

	if userInfo.IsVerified {
		return ErrUserIsVerified
	}

	if !u.hash.CompareValue(userInfo.Password, request.Password) {
		return ErrPasswordIsIncorrect
	}

	userInfo.IsVerified = true

	return u.store.UpdateUser(userInfo)
}
