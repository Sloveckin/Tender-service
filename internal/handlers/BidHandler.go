package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"log/slog"
	"net/http"
	"tender-service/internal/requests"
	s "tender-service/internal/service"
)

type BidHandler struct {
	tenderService       *s.TenderService
	userService         *s.UserService
	organizationService *s.OrganizationService
	bidService          *s.BidService
	logger              *slog.Logger
}

func InitBidHandler(t *s.TenderService, u *s.UserService, o *s.OrganizationService, b *s.BidService, l *slog.Logger) *BidHandler {
	return &BidHandler{tenderService: t, userService: u, organizationService: o, bidService: b, logger: l}
}

func (h *BidHandler) CreateBid() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request requests.BidCreateRequest
		if !CheckJSON(w, r, &request, h.logger) {
			return
		}

		if !ValidateJSON(w, r, request, h.logger) {
			return
		}

		bid, err := h.bidService.CreateBid(&request)
		if err != nil {
			h.logger.Error(err.Error())
			BadRequestWithServiceError(w, r, err)
			return
		}

		Create(w, r, bid)
	}
}

func (h *BidHandler) MyBid() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query()
		if !param.Has("username") {
			h.logger.Info(fmt.Sprintf("BidCreateRequest with url=%s denied. Username not provided", r.URL.String()))
			RequestWithMessage(w, r, "Username not provided", http.StatusBadRequest)
			return
		}

		username := param.Get("username")

		bids, err := h.bidService.GetUserBids(username)
		if err != nil {
			h.logger.Error(err.Error())
			BadRequestWithServiceError(w, r, err)
			return
		}
		h.logger.Info(fmt.Sprintf("Bids by postgres %s get successfully", username))
		OkResponse(w, r, bids)

	}
}

func (h *BidHandler) BidStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			h.logger.Info(fmt.Sprintf("Request with url=%s denied. No tender id", r.URL.String()))
			RequestWithMessage(w, r, "Missing tender id", http.StatusBadRequest)
			return
		}

		if !IsValidUUID(id) {
			h.logger.Info(fmt.Sprintf("Request with url=%s denied. Invalid id format", r.URL.String()))
			RequestWithMessage(w, r, fmt.Sprintf("%s not valid uuid", id), http.StatusBadRequest)
			return
		}

		bid, err := h.bidService.GetBidById(id)
		if err != nil {
			h.logger.Error(err.Error())
			BadRequestWithServiceError(w, r, err)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(bid.Status))
	}
}

func (h *BidHandler) ChangeBidStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			h.logger.Info(fmt.Sprintf("Request with url=%s denied. No tender id", r.URL.String()))
			RequestWithMessage(w, r, "Missing tender id", http.StatusBadRequest)
			return
		}

		if !IsValidUUID(id) {
			h.logger.Info(fmt.Sprintf("Request with url=%s denied. Invalid id format", r.URL.String()))
			RequestWithMessage(w, r, fmt.Sprintf("%s not valid uuid", id), http.StatusBadRequest)
			return
		}

		param := r.URL.Query()
		if !param.Has("status") {
			h.logger.Info(fmt.Sprintf("ChangeBidStatus with url=%s denied. staus not provided", r.URL.String()))
			RequestWithMessage(w, r, "status not provided", http.StatusBadRequest)
			return
		}

		status := param.Get("status")

		bid, err := h.bidService.ChangeBidStatusById(id, status)
		if err != nil {
			h.logger.Error(err.Error())
			BadRequestWithServiceError(w, r, err)
			return
		}
		OkResponse(w, r, bid)
	}
}

func (h *BidHandler) TenderBids() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			h.logger.Info(fmt.Sprintf("Request with url=%s denied. No tender id", r.URL.String()))
			RequestWithMessage(w, r, "Missing tender id", http.StatusBadRequest)
			return
		}

		if !IsValidUUID(id) {
			h.logger.Info(fmt.Sprintf("Request with url=%s denied. Invalid id format", r.URL.String()))
			RequestWithMessage(w, r, fmt.Sprintf("%s not valid uuid", id), http.StatusBadRequest)
			return
		}

		bids, err := h.bidService.GetTenderBids(id)
		if err != nil {
			h.logger.Error(err.Error())
			BadRequestWithServiceError(w, r, err)
			return
		}
		h.logger.Info(fmt.Sprintf("Bids by tender with id=%s get successfully", id))
		OkResponse(w, r, bids)
	}
}
