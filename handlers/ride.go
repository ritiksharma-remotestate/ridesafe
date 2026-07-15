package handlers

import (
	"errors"
	"net/http"
	"ridesafe/error_custom"
	middlewares "ridesafe/middleware"
	"ridesafe/models"

	"ridesafe/services"
	"ridesafe/utils"
)

func CreateRide(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRideRequest
	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondError(
			w,
			http.StatusBadRequest,
			err,
			"failed to parse request body",
		)
		return
	}

	if errs := utils.CheckValidation(req); errs != nil {
		utils.RespondError(
			w,
			http.StatusBadRequest,
			nil,
			"invalid request",
			errs.Error(),
		)
		return
	}

	user := middlewares.ClaimsContext(r)
	if user == nil {
		utils.RespondError(
			w,
			http.StatusUnauthorized,
			nil,
			"user not authenticated",
		)
		return
	}
	ride, err := services.CreateRide(user.UserID, req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create ride")
		return
	}
	utils.RespondJSON(w, http.StatusOK, ride)

}

func AcceptRide(w http.ResponseWriter, r *http.Request) {
	rideID := r.PathValue("id")
	if err := utils.ParseBody(r.Body, &rideID); err != nil {
		utils.RespondError(
			w,
			http.StatusBadRequest,
			err,
			"failed to parse request body",
		)
		return
	}

	user := middlewares.ClaimsContext(r)
	if user == nil {
		utils.RespondError(
			w,
			http.StatusUnauthorized,
			nil,
			"user not authenticated",
		)
		return
	}
	err := services.AcceptRide(rideID, user.UserID)
	if err != nil {
		switch {
		case errors.Is(err, error_custom.ErrRideNotFound):
			utils.RespondError(w, http.StatusNotFound, err, "ride not found")

		case errors.Is(err, error_custom.ErrDriverOffline):
			utils.RespondError(w, http.StatusBadRequest, err, "driver is offline")

		case errors.Is(err, error_custom.ErrDriverUnavailable):
			utils.RespondError(w, http.StatusConflict, err, "driver unavailable")

		case errors.Is(err, error_custom.ErrUnauthorizedDriver):
			utils.RespondError(w, http.StatusForbidden, err, "unauthorized")

		default:
			utils.RespondError(w, http.StatusInternalServerError, err, "failed to accept ride")
		}
		return
	}

	utils.RespondJSON(w, http.StatusOK, "ride accepted")

}

func ArriveRide(w http.ResponseWriter, r *http.Request) {

	// rideID:=r.PathValue("id")
	rideID := r.PathValue("id")
	if err := utils.ParseBody(r.Body, &rideID); err != nil {
		utils.RespondError(
			w,
			http.StatusBadRequest,
			err,
			"failed to parse request body",
		)
		return
	}

	user := middlewares.ClaimsContext(r)
	if user == nil {
		utils.RespondError(
			w,
			http.StatusUnauthorized,
			nil,
			"user not authenticated",
		)
		return
	}
	err := services.MarkRideArrived(rideID, user.UserID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err, "failed to update")
		return
	}

	utils.RespondJSON(w, http.StatusOK, "ride arrived")

}

func StartRide(w http.ResponseWriter, r *http.Request) {

	rideID := r.PathValue("id")
	if err := utils.ParseBody(r.Body, &rideID); err != nil {
		utils.RespondError(
			w,
			http.StatusBadRequest,
			err,
			"failed to parse request body",
		)
		return
	}

	user := middlewares.ClaimsContext(r)
	if user == nil {
		utils.RespondError(
			w,
			http.StatusUnauthorized,
			nil,
			"user not authenticated",
		)
		return
	}
	err := services.StartRide(rideID, user.UserID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err, "failed to update")
		return
	}

	utils.RespondJSON(w, http.StatusOK, "ride started")

}

func CompleteRide(w http.ResponseWriter, r *http.Request) {

	rideID := r.PathValue("id")
	if err := utils.ParseBody(r.Body, &rideID); err != nil {
		utils.RespondError(
			w,
			http.StatusBadRequest,
			err,
			"failed to parse request body",
		)
		return
	}

	user := middlewares.ClaimsContext(r)
	if user == nil {
		utils.RespondError(
			w,
			http.StatusUnauthorized,
			nil,
			"user not authenticated",
		)
		return
	}
	err := services.CompleteRide(rideID, user.UserID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err, "failed to update")
		return
	}

	utils.RespondJSON(w, http.StatusOK, "ride completed")

}

func CancelRide(w http.ResponseWriter, r *http.Request) {

	rideID := r.PathValue("id")
	if err := utils.ParseBody(r.Body, &rideID); err != nil {
		utils.RespondError(
			w,
			http.StatusBadRequest,
			err,
			"failed to parse request body",
		)
		return
	}

	user := middlewares.ClaimsContext(r)
	if user == nil {
		utils.RespondError(
			w,
			http.StatusUnauthorized,
			nil,
			"user not authenticated",
		)
		return
	}
	err := services.CancelRide(rideID, user.UserID, user.Role)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err, "failed to update")
		return
	}

	utils.RespondJSON(w, http.StatusOK, "ride cancelled")

}

func GetRideByID(w http.ResponseWriter, r *http.Request) {

}

// func GetMyRides(w http.ResponseWriter, r *http.Request)
