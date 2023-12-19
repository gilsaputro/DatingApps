package partner

import (
	"context"
	"encoding/json"
	"fmt"
	"gilsaputro/dating-apps/internal/handler/utilhttp"
	"gilsaputro/dating-apps/internal/service/partner"
	"log"
	"net/http"
	"time"
)

// PartnerPartnerResponse is list response parameter for Login Api
type PartnerResponse struct {
	PartnerID   int    `json:"id"`
	Fullname    string `json:"fullname"`
	Status      string `json:"status"`
	CreatedDate string `json:"created_date"`
}

// PassPartnerHandler is func handler for generate new partner
func (h *PartnerHandler) PassPartnerHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.timeoutInSec)*time.Second)
	defer cancel()

	var err error
	var response utilhttp.StandardResponse
	var code int = http.StatusOK

	defer func() {
		response.Code = code
		if err == nil {
			response.Message = "success"
		} else {
			response.Message = err.Error()
		}

		data, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			log.Println("[PartnerPartnerHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var userID int
	var ok bool
	userID, ok = r.Context().Value("id").(int)
	if !ok {
		code = http.StatusInternalServerError
		err = fmt.Errorf("Internal Server Error")
		return
	}

	var isVerified bool
	isVerified, ok = r.Context().Value("isverified").(bool)
	if !ok {
		code = http.StatusInternalServerError
		err = fmt.Errorf("Internal Server Error")
		return
	}

	errChan := make(chan error, 1)
	var partnerInfo partner.PartnerServiceInfo
	go func(ctx context.Context) {
		partnerInfo, err = h.service.PassPartner(partner.PartnerServiceRequest{
			UserID:     userID,
			IsVerified: isVerified,
		})
		errChan <- err
	}(ctx)

	select {
	case <-ctx.Done():
		code = http.StatusGatewayTimeout
		err = fmt.Errorf("Timeout")
		return
	case err = <-errChan:
		if err != nil {
			code = http.StatusInternalServerError
			return
		}
	}

	response = mapResponse(partnerInfo)
}
