package handlers

import (
	"errors"
	"fmt"
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
	fmt.Println(rideID)
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

	rideID := r.PathValue("id")

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
	otp, err := services.MarkRideArrived(rideID, user.UserID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err, "failed to mark ride arrived")
		return
	}

	utils.RespondJSON(w, http.StatusOK, otp)

}

func StartRide(w http.ResponseWriter, r *http.Request) {

	rideID := r.PathValue("id")

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
	rideID := r.PathValue("id")
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
	ride, err := services.GetRideByID(rideID, user.UserID, user.Role)
	if err != nil {
		utils.RespondError(
			w,
			http.StatusNotFound,
			nil,
			"ride not found",
		)
		return
	}
	utils.RespondJSON(
		w,
		http.StatusOK,
		ride,
	)

}
func VerifyRideOTP(w http.ResponseWriter, r *http.Request) {

	var req models.VerifyOTPRequest

	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondError(
			w,
			http.StatusBadRequest,
			err,
			"failed to parse request",
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

	rideID := r.PathValue("id")

	err := services.VerifyRideOTP(rideID, user.UserID, req.OTP)

	if err != nil {

		switch {

		case errors.Is(err, error_custom.ErrRideNotFound):
			utils.RespondError(
				w,
				http.StatusNotFound,
				err,
				"ride not found",
			)

		case errors.Is(err, error_custom.ErrUnauthorizedDriver):
			utils.RespondError(
				w,
				http.StatusForbidden,
				err,
				"unauthorized driver",
			)

		case errors.Is(err, error_custom.ErrInvalidOTP):
			utils.RespondError(
				w,
				http.StatusBadRequest,
				err,
				"invalid otp",
			)
		case errors.Is(err, error_custom.ErrTooManyOTPAttempts):
			utils.RespondError(
				w,
				http.StatusTooManyRequests,
				err,
				"too many incorrect otp attempts, please request a new ride or contact support",
			)
		default:
			utils.RespondError(
				w,
				http.StatusInternalServerError,
				err,
				"internal server error",
			)
		}

		return
	}

	utils.RespondJSON(
		w,
		http.StatusOK,
		map[string]string{
			"message": "OTP verified successfully",
		},
	)
}
