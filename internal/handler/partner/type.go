package partner

import (
	"gilsaputro/dating-apps/internal/handler/utilhttp"
	"gilsaputro/dating-apps/internal/service/partner"
)

// PartnerPartnerResponse is list response parameter for Login Api
type PartnerResponse struct {
	PartnerID   int    `json:"id"`
	Fullname    string `json:"fullname"`
	Status      string `json:"status"`
	IsVerified  bool   `json:"is_verified"`
	CreatedDate string `json:"created_date"`
}

func mapResponse(result partner.PartnerServiceInfo) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	data := PartnerResponse{
		PartnerID:   result.PartnerID,
		Fullname:    result.Fullname,
		Status:      result.Status,
		IsVerified:  result.IsVerified,
		CreatedDate: result.CreatedDate,
	}
	res.Data = data
	return res
}
