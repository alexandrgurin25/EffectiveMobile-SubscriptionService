package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"subscriptions/internal/transport/http/dto/subscription"
	"subscriptions/pkg/logger"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	idStr := chi.URLParam(r, "id")

	UUID, err := uuid.Parse(idStr)

	id := UUID.String()

	if err != nil {
		errStr := "Invalid UUID"
		h.sendError(w, http.StatusBadRequest, errStr)
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.Any("id", idStr),
			zap.Error(err))
		return
	}

	gettedSub, err := h.service.GetById(ctx, id)
	if err != nil {
		var errStr string
		if errors.Is(err, sql.ErrNoRows) {
			errStr = "Subscription not found"
			h.sendError(w, http.StatusNotFound, errStr)
		} else {
			errStr = "Failed to fetch subscription"
			h.sendError(w, http.StatusInternalServerError, errStr)
		}

		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.Any("id", idStr),
			zap.Error(err))
		return
	}

	res := subscription.SubResponse{
		Id:        gettedSub.Id,
		Name:      gettedSub.Name,
		Price:     gettedSub.Price,
		UserId:    gettedSub.UserId,
		StartDate: gettedSub.StartDate,
		EndData:   gettedSub.EndDate,
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"Subscription get successfully!",
		zap.Any("res", idStr),
		zap.Error(err))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
