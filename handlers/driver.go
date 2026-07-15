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

func CreateDriver(w http.ResponseWriter, r *http.Request) {
	claims := middlewares.ClaimsContext(r)
	if claims == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized user")
		return
	}

	driverID := claims.UserID
	// driverCTX:=middlewares.ClaimsContext(r)
	err := services.CreateDriver(driverID)
	if err != nil {
		switch {
		case errors.Is(err, error_custom.ErrUserNotFound):
			utils.RespondError(w, http.StatusNotFound, err, "user not found")

		case errors.Is(err, error_custom.ErrNotADriver):
			utils.RespondError(w, http.StatusForbidden, err, "user is not a driver")

		case errors.Is(err, error_custom.ErrDriverAlreadyExists):
			utils.RespondError(w, http.StatusConflict, err, "driver already exists")

		default:
			utils.RespondError(w, http.StatusInternalServerError, err, "internal server error")
		}
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{
		"message": "driver created successfully",
	})
}

func GetMyDriverProfile(w http.ResponseWriter, r *http.Request) {
	claims := middlewares.ClaimsContext(r)
	if claims == nil {
		// unauthorized
		return
	}

	driverID := claims.UserID
	driver, err := services.GetDriverByUserID(driverID)
	if err != nil {
		switch {
		case errors.Is(err, error_custom.ErrUserNotFound):
			utils.RespondError(w, http.StatusNotFound, err, "driver not found")

		default:
			utils.RespondError(w, http.StatusInternalServerError, err, "internal server error")
		}
		return
	}

	utils.RespondJSON(w, http.StatusOK, driver)
}

func UpdateDriverLocation(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateDriverLocationRequest

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
	err := services.UpdateDriverLocation(user.UserID, req.Latitude, req.Longitude)
	if err != nil {
		switch {
		case errors.Is(err, error_custom.ErrUserNotFound):
			utils.RespondError(w, http.StatusNotFound, err, "driver not found")

		default:
			utils.RespondError(w, http.StatusInternalServerError, err, "failed to update location")
		}
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "location updated successfully",
	})
}

func SetDriverOnline(w http.ResponseWriter, r *http.Request) {
	var status models.DriverOnlineRequest
	if err := utils.ParseBody(r.Body, &status); err != nil {
		utils.RespondError(
			w,
			http.StatusBadRequest,
			err,
			"failed to parse request body",
		)
		return
	}

	if errs := utils.CheckValidation(status); errs != nil {
		utils.RespondError(
			w,
			http.StatusBadRequest,
			nil,
			"invalid request",
			errs.Error(),
		)
		return
	}
	claims := middlewares.ClaimsContext(r)
	if claims == nil {
		// unauthorized
		return
	}

	driverID := claims.UserID
	// driverCTX:=middlewares.ClaimsContext(r)
	err := services.DriverOnline(driverID, status.Online)
	if err != nil {
		switch {
		case errors.Is(err, error_custom.ErrUserNotFound):
			utils.RespondError(w, http.StatusNotFound, err, "driver not found")

		default:
			utils.RespondError(w, http.StatusInternalServerError, err, "failed to update online status")
		}
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "driver online status updated",
	})
}

func SetDriverAvailable(w http.ResponseWriter, r *http.Request) {
	var status models.DriverAvailableRequest
	claims := middlewares.ClaimsContext(r)
	if claims == nil {
		utils.RespondError(w, http.StatusForbidden, nil, "claims not found")
		return
	}

	driverID := claims.UserID
	if err := utils.ParseBody(r.Body, &status); err != nil {
		utils.RespondError(
			w,
			http.StatusBadRequest,
			err,
			"failed to parse request body",
		)
		return
	}

	if errs := utils.CheckValidation(status); errs != nil {
		utils.RespondError(
			w,
			http.StatusBadRequest,
			nil,
			"invalid request",
			errs.Error(),
		)
		return
	}
	// driverCTX:=middlewares.ClaimsContext(r)
	err := services.DriverAvailable(driverID, status.Available)
	if err != nil {
		switch {
		case errors.Is(err, error_custom.ErrUserNotFound):
			utils.RespondError(w, http.StatusNotFound, err, "driver not found")

		default:
			utils.RespondError(w, http.StatusInternalServerError, err, "failed to update availability")
		}
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "driver availability updated",
	})
}

func GetAvailableDrivers(w http.ResponseWriter, r *http.Request) {
	rideID:= r.PathValue("id")
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

	drivers, err := services.GetAvailableDrivers(rideID,user.UserID)
	if err != nil {
		utils.RespondError(
			w,
			http.StatusInternalServerError,
			err,
			"failed to fetch available drivers",
		)
		return
	}

	utils.RespondJSON(
		w,
		http.StatusOK,
		drivers,
	)
}

func GetDriverLocation(w http.ResponseWriter, r *http.Request) {
	driverID := r.PathValue("id")

	lat, lon, err := services.GetDriverLocation(driverID)
	if err != nil {
		utils.RespondError(
			w,
			http.StatusInternalServerError,
			err,
			"failed to fetch driver location",
		)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]float64{
		"latitude":  lat,
		"longitude": lon,
	})
}
