package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"subscriptions/pkg/logger"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Delete removes subscription by ID
// @Summary Удаление подписки по ID
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID in UUID format"
// @Success 204 "Subscription deleted successfully"
// @Failure 400 {object} subscription.ErrorResponse "Invalid format for UUID in `id`"
// @Failure 404 {object} subscription.ErrorResponse "Subscription not found"
// @Failure 500 {object} subscription.ErrorResponse "Internal server error"
// @Router /api/subscriptions/{id} [delete]
func (h Handlers) Delete(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.DeleteById(ctx, id)

	if err != nil {
		var errStr string
		if errors.Is(err, sql.ErrNoRows) {
			errStr = "Subscription not found"
			h.sendError(w, http.StatusNotFound, errStr)
		} else {
			errStr = "Couldn't delete subscription"
			h.sendError(w, http.StatusInternalServerError, errStr)
		}

		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.Any("id", idStr),
			zap.Error(err))
		return
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"Subscription delete successfully!",
		zap.Any("id", id))
	w.WriteHeader(http.StatusNoContent)
}
