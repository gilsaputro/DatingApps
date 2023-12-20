package user

import (
	"context"
	"encoding/json"
	"fmt"
	"gilsaputro/dating-apps/internal/handler/utilhttp"
	"gilsaputro/dating-apps/internal/service/user"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// UpgradeUserRequest is list request parameter for Upgrade Api
type UpgradeUserRequest struct {
	Password string `json:"password"`
}

// UpgradeUserHandler is func handler for Upgrade user
func (h *UserHandler) UpgradeUserHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[UpgradeUserHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body UpgradeUserRequest
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		code = http.StatusBadRequest
		err = fmt.Errorf("Bad Request")
		return
	}

	err = json.Unmarshal(data, &body)
	if err != nil {
		code = http.StatusBadRequest
		err = fmt.Errorf("Bad Request")
		return
	}

	// checking valid body
	if len(body.Password) < 1 {
		code = http.StatusBadRequest
		err = fmt.Errorf("Invalid Parameter Request")
		return
	}

	var userID int
	var ok bool
	userID, ok = r.Context().Value("id").(int)
	if !ok {
		code = http.StatusInternalServerError
		err = fmt.Errorf("Internal Server Error")
		return
	}

	errChan := make(chan error, 1)
	go func(ctx context.Context) {
		err = h.service.UpgradeUser(user.UpgradeServiceRequest{
			UserId:   userID,
			Password: body.Password,
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
			if err == user.ErrUserNameNotExists || err == user.ErrPasswordIsIncorrect {
				code = http.StatusBadRequest
			} else {
				code = http.StatusInternalServerError
			}
			return
		}
	}

	response = mapResponseUpgrade()
}

func mapResponseUpgrade() utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	return res
}
