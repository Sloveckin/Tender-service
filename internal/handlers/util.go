package handlers

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strings"
	"tender-service/internal/data"
	"tender-service/internal/service"
)

func OkResponse(w http.ResponseWriter, r *http.Request, a interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, a)
}

func BadRequestWithServiceError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ServiceErrorToStatusCode(err))
	json := data.ErrorResponse{Message: err.Error()}
	render.JSON(w, r, json)
}

func RequestWithMessage(w http.ResponseWriter, r *http.Request, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json := data.ErrorResponse{Message: message}
	render.JSON(w, r, json)
}

func Create(w http.ResponseWriter, r *http.Request, a interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, a)
}

func ValidationError(w http.ResponseWriter, r *http.Request, errors validator.ValidationErrors) {
	var errMsg []string
	for _, err := range errors {
		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("Field %s is required", err.Field()))
		case "max":
			errMsg = append(errMsg, fmt.Sprintf("Field %s must be less than %s", err.Field(), err.Param()))
		case "oneof":
			errMsg = append(errMsg, fmt.Sprintf("Field %s must be one of %s", err.Field(), err.Param()))
		case "uuid":
			errMsg = append(errMsg, fmt.Sprintf("Field %s must be a valid UUID", err.Field()))
		}
	}
	RequestWithMessage(w, r, strings.Join(errMsg, ", "), http.StatusBadRequest)
}

func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

func ValidateJSON(w http.ResponseWriter, r *http.Request, request interface{}, logger *slog.Logger) bool {
	if err := validator.New().Struct(request); err != nil {
		logger.Info(err.Error())
		ValidationError(w, r, err.(validator.ValidationErrors))
		return false
	}
	return true
}

func CheckJSON(w http.ResponseWriter, r *http.Request, request interface{}, logger *slog.Logger) bool {
	err := render.DecodeJSON(r.Body, &request)
	if err != nil {
		logger.Info(err.Error())
		RequestWithMessage(w, r, "Invalid JSON request", http.StatusBadRequest)
		return false
	}
	return true
}

func ServiceErrorToStatusCode(err error) int {
	mp := map[error]int{
		service.TenderNotFound:              404,
		service.OrganizationNotFound:        404,
		service.BidNotFound:                 404,
		service.UserNotFound:                401,
		service.TenderAlreadyExists:         400,
		service.NotCorrectParams:            400,
		service.NotExpectedErrorFromStorage: 500,
	}
	res, ok := mp[err]
	if ok {
		return res
	}
	return mp[errors.Unwrap(err)]
}
