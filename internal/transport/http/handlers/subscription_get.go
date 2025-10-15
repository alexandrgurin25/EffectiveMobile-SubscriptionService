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

	if err != nil {
		errStr := "Invalid UUID"
		h.sendError(w, http.StatusBadRequest, errStr)
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.Any("id", idStr),
			zap.Error(err))
		return
	}

	id := UUID.String()

	gotSub, err := h.service.GetById(ctx, id)
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
		Id:        gotSub.Id,
		Name:      gotSub.Name,
		Price:     gotSub.Price,
		UserId:    gotSub.UserId,
		StartDate: gotSub.StartDate,
		EndDate:   gotSub.EndDate,
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"Subscription got successfully!",
		zap.Any("res", res))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
