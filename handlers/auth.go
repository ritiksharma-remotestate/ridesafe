package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"ridesafe/error_custom"
	"ridesafe/models"

	// "ridesafe/repository"
	"ridesafe/services"
	"ridesafe/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid requestbody", http.StatusBadRequest)
		return
	}
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(req)nnknknkn

	if errs := utils.CheckValidation(user); errs != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid request", errs.Error())

		return
	}

	_, err = services.RegisterUser(user)

	if err != nil {
		switch {
		case errors.Is(err, error_custom.ErrUserAlreadyExists):
			utils.RespondError(w, http.StatusConflict, err, "user already exists")

		default:
			utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user")
		}
		return
	}
	utils.RespondJSON(w, http.StatusCreated, map[string]string{
		"message": "user created successfully",
	})

}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid requestbody", http.StatusBadRequest)
		return
	}

	if errs := utils.CheckValidation(req); errs != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "invalid request", errs.Error())

		return
	}
	//  iknow this users struct is unnecessary but i already wrote code using it so i am not removing for now

	token, err := services.Login(req.Email, req.Password)

	if err != nil {
    switch {
    case errors.Is(err, error_custom.ErrInvalidCredentials):
        utils.RespondError(w, http.StatusUnauthorized, err, "invalid email or password")

    default:
        utils.RespondError(w, http.StatusInternalServerError, err, "failed to login")
    }
    return
}
	utils.RespondJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})

}
