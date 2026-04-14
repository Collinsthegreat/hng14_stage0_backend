package handler

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/Collinsthegreat/hng14_stage0_backend/internal/service"
	"github.com/Collinsthegreat/hng14_stage0_backend/pkg/response"
	"github.com/go-chi/chi/v5"
)

var validNameRegex = regexp.MustCompile(`^[a-zA-Z\-\s]+$`)

type ClassifyHandler struct {
	svc service.ClassifyService
}

func NewClassifyHandler(svc service.ClassifyService) *ClassifyHandler {
	return &ClassifyHandler{svc: svc}
}

func RegisterRoutes(r *chi.Mux, h *ClassifyHandler) {
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.Error(w, http.StatusNotFound, "route not found")
	})

	r.Get("/api/classify", h.HandleClassify)
}

func (h *ClassifyHandler) HandleClassify(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.URL.Query().Get("name"))
	if name == "" {
		response.Error(w, http.StatusBadRequest, "name query parameter is required")
		return
	}

	if !validNameRegex.MatchString(name) {
		response.Error(w, http.StatusUnprocessableEntity, "name must be a valid string")
		return
	}

	name = strings.ToLower(name)

	data, err := h.svc.Classify(r.Context(), name)
	if err != nil {
		if errors.Is(err, service.ErrNoPrediction) {
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "timeout") || strings.Contains(msg, "context deadline exceeded") || strings.Contains(msg, "unexpected status") || strings.Contains(msg, "http request failed") {
			response.Error(w, http.StatusBadGateway, "upstream service unavailable")
			return
		}
		
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, response.SuccessResponse{
		Status: "success",
		Data:   data,
	})
}
