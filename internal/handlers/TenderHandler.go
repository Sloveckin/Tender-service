package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"log/slog"
	"net/http"
	"strconv"
	"tender-service/internal/requests"
	"tender-service/internal/service"
)

type TenderHandler struct {
	tenderService       *service.TenderService
	userService         *service.UserService
	organizationService *service.OrganizationService
	logger              *slog.Logger
}

func InitTenderHandler(t *service.TenderService, u *service.UserService, o *service.OrganizationService, l *slog.Logger) *TenderHandler {
	return &TenderHandler{tenderService: t, userService: u, organizationService: o, logger: l}
}

func (t *TenderHandler) CreateTender() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request requests.CreateTenderRequest
		if !CheckJSON(w, r, &request, t.logger) {
			return
		}

		if !ValidateJSON(w, r, request, t.logger) {
			return
		}

		result, err := t.tenderService.CreateTender(&request)
		if err != nil {
			t.logger.Error(err.Error())
			BadRequestWithServiceError(w, r, err)
			return
		}

		t.logger.Info(fmt.Sprintf("Tender with name %s created. id=%s", result.Name, result.Id))
		Create(w, r, result)
	}
}

func (t *TenderHandler) Edit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			t.logger.Info(fmt.Sprintf("Request with url=%s denied. No tender id", r.URL.String()))
			RequestWithMessage(w, r, "Missing tender id", http.StatusBadRequest)
			return
		}

		if !IsValidUUID(id) {
			t.logger.Info(fmt.Sprintf("Request with url=%s denied. Invalid id format", r.URL.String()))
			RequestWithMessage(w, r, fmt.Sprintf("%s not valid uuid", id), http.StatusBadRequest)
			return
		}

		var request requests.EditTenderRequest
		if !CheckJSON(w, r, &request, t.logger) {
			return
		}

		if !ValidateJSON(w, r, request, t.logger) {
			return
		}

		tender, err := t.tenderService.EditTender(id, &request)
		if err != nil {
			t.logger.Error(err.Error())
			BadRequestWithServiceError(w, r, err)
			return
		}
		t.logger.Info(fmt.Sprintf("Tender with id=%s edit successfully", id))
		Create(w, r, tender)
	}
}

func (t *TenderHandler) Tenders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenders, err := t.tenderService.GetTenders()
		if err != nil {
			t.logger.Error(err.Error())
			BadRequestWithServiceError(w, r, err)
			return
		}
		t.logger.Info("Tender get successfully")
		OkResponse(w, r, tenders)
	}
}

func (t *TenderHandler) MyTenders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query()
		if !param.Has("username") {
			t.logger.Info(fmt.Sprintf("BidCreateRequest with url=%s denied. Username not provided", r.URL.String()))
			RequestWithMessage(w, r, "Username not provided", http.StatusBadRequest)
			return
		}

		username := param.Get("username")
		tenders, err := t.userService.GetUserTendersByUsername(username)
		if err != nil {
			t.logger.Error(err.Error())
			BadRequestWithServiceError(w, r, err)
			return
		}

		if tenders == nil {
			t.logger.Info(fmt.Sprintf("Tenders by postgres %s not found", username))
			OkResponse(w, r, nil)
		} else {
			t.logger.Info(fmt.Sprintf("Tender by postgres %s get successfully", username))
			OkResponse(w, r, tenders)
		}

	}
}

func (t *TenderHandler) TenderRollback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			t.logger.Info(fmt.Sprintf("Request with url=%s denied. No tender id", r.URL.String()))
			RequestWithMessage(w, r, "Missing tender id", http.StatusBadRequest)
			return
		}

		if !IsValidUUID(id) {
			t.logger.Info(fmt.Sprintf("Request with url=%s denied. Invalid id format", r.URL.String()))
			RequestWithMessage(w, r, fmt.Sprintf("%s not valid uuid", id), http.StatusBadRequest)
			return
		}

		versionStr := chi.URLParam(r, "version")
		if versionStr == "" {
			t.logger.Info(fmt.Sprintf("BidCreateRequest with url=%s denied. Version not found", r.URL.String()))
			RequestWithMessage(w, r, "Missing version", http.StatusBadRequest)
			return
		}

		version, err := strconv.Atoi(versionStr)
		if err != nil {
			t.logger.Info(fmt.Sprintf("BidCreateRequest with url=%s denied. Version not positive number", r.URL.String()))
			RequestWithMessage(w, r, "Version not positive number", http.StatusBadRequest)
			return
		}

		tender, err := t.tenderService.RollbackTender(id, int64(version))
		if err != nil {
			t.logger.Error(err.Error())
			BadRequestWithServiceError(w, r, err)
			return
		}

		OkResponse(w, r, tender)
	}
}
