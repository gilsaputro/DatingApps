package user

import (
	"gilsaputro/dating-apps/internal/handler/utilhttp"
	"gilsaputro/dating-apps/internal/service/user"
)

// UpgradeUserRequest is list request parameter for Upgrade Api
type UpgradeUserRequest struct {
	Password string `json:"password"`
}

type UserProfile struct {
	UserID      int    `json:"id"`
	Username    string `json:"username"`
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	IsVerified  bool   `json:"is_verified"`
	CreatedDate string `json:"created_date"`
}

func mapResponseUserProfile(profile user.UserServiceInfo) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	res.Data = UserProfile{
		UserID:      profile.UserId,
		Username:    profile.Username,
		Fullname:    profile.Fullname,
		Email:       profile.Email,
		IsVerified:  profile.IsVerified,
		CreatedDate: profile.CreatedDate,
	}
	return res
}

// EditUserRequest is list request parameter for Edit Api
type EditUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

// EditUserResponse is list response parameter for Edit Api
type EditUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

func mapResponseEdit(result user.UserServiceInfo) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse

	data := EditUserResponse{
		Username: result.Username,
		Email:    result.Email,
		Fullname: result.Fullname,
	}

	res.Data = data
	return res
}
