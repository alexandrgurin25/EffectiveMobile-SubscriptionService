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

func (h *Handlers) GetList(w http.ResponseWriter, r *http.Request) {

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

	gotSubs, err := h.service.GetList(ctx, id)
	if err != nil {
		var errStr string
		if errors.Is(err, sql.ErrNoRows) {
			errStr = "Subscriptions not found"
			h.sendError(w, http.StatusNotFound, errStr)
		} else {
			errStr = "Failed to fetch subscriptions"
			h.sendError(w, http.StatusInternalServerError, errStr)
		}

		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.Any("id", idStr),
			zap.Error(err))
		return
	}

	responses := make([]subscription.SubResponse, 0, len(gotSubs))

	for _, sub := range gotSubs {
		res := subscription.SubResponse{
			Id:        sub.Id,
			Name:      sub.Name,
			Price:     sub.Price,
			UserId:    sub.UserId,
			StartDate: sub.StartDate,
			EndData:   sub.EndDate,
		}

		responses = append(responses, res)
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"List subscriptions got successfully!",
		zap.Any("res", responses))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}
