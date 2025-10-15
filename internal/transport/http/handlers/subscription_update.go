package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"subscriptions/internal/entity"
	"subscriptions/internal/transport/http/dto/subscription"
	"subscriptions/pkg/logger"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Put updates subscription by ID
// @Summary Обновление подписки по ID
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID in UUID format"
// @Param input body subscription.SubRequest true "Subscription update data"
// @Success 200 {object} subscription.SubResponse "Updated subscription details"
// @Failure 400 {object} subscription.ErrorResponse "Invalid format for UUID in `id`, JSON, or validation error"
// @Failure 404 {object} subscription.ErrorResponse "Subscription not found"
// @Failure 500 {object} subscription.ErrorResponse "Internal server error"
// @Router /api/subscriptions/{id} [put]
func (h *Handlers) Put(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")

	UUID, err := uuid.Parse(idStr)

	if err != nil {
		errStr := "Invalid format for UUID in `id`"
		h.sendError(w, http.StatusBadRequest, errStr)
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.Any("id", idStr),
			zap.Error(err))
		return
	}

	id := UUID.String()

	var req subscription.SubRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"Failed to decode JSON request",
			zap.Error(err),
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.Any("headers", r.Header),
		)
		h.sendError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	defer r.Body.Close()

	if req.Price < 0 || req.Name == "" ||
		req.UserId == "" || req.StartDate == "" {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"Empty fields in json or negative <Price>",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.Any("headers", r.Header),
		)
		h.sendError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	updateSubscription := entity.Subscription{
		Id:        id,
		Name:      req.Name,
		Price:     req.Price,
		UserId:    req.UserId,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	putSub, err := h.service.UpdateById(ctx, &updateSubscription)

	if err != nil {
		var errStr string
		if errors.Is(err, sql.ErrNoRows) {
			errStr = "Subscription not found"
			h.sendError(w, http.StatusNotFound, errStr)
		} else {
			errStr = "Couldn't renew subscription"
			h.sendError(w, http.StatusInternalServerError, errStr)
		}

		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.Any("id", idStr),
			zap.Error(err))
		return
	}

	res := subscription.SubResponse{
		Id:        putSub.Id,
		Name:      putSub.Name,
		Price:     putSub.Price,
		UserId:    putSub.UserId,
		StartDate: putSub.StartDate,
		EndDate:   putSub.EndDate,
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"Subscription put successfully!",
		zap.Any("res", res))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
