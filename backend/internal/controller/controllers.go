// หมวย
package controller

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"kencatexpress/backend/internal/dto"
	"kencatexpress/backend/internal/middleware"
	"kencatexpress/backend/internal/service"
	"kencatexpress/backend/internal/util"
)

type API struct {
	auth      *service.AuthService
	users     *service.UserService
	shipping  *service.ShippingService
	parcels   *service.ParcelService
	tracking  *service.TrackingService
	messenger *service.MessengerService
	vehicles  *service.VehicleService
	reports   *service.ReportService
	tokenTTL  time.Duration
}

func NewAPI(auth *service.AuthService, users *service.UserService, shipping *service.ShippingService, parcels *service.ParcelService, tracking *service.TrackingService, messenger *service.MessengerService, vehicles *service.VehicleService, reports *service.ReportService, tokenTTL time.Duration) *API {
	return &API{
		auth:      auth,
		users:     users,
		shipping:  shipping,
		parcels:   parcels,
		tracking:  tracking,
		messenger: messenger,
		vehicles:  vehicles,
		reports:   reports,
		tokenTTL:  tokenTTL,
	}
}

// หมวย 7) GET /api/trackings/{trackId}
func (a *API) GetTracking(w http.ResponseWriter, r *http.Request) {
	tracking, err := a.tracking.GetTracking(r.Context(), r.PathValue("trackId"))
	if err != nil {
		util.ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, tracking)
}

// หมวย 8) PUT /api/trackings/{trackId}/status
func (a *API) UpdateTrackingStatus(w http.ResponseWriter, r *http.Request) {
	var req dto.TrackingStatusUpdateRequest
	if err := util.DecodeJSON(r, &req); err != nil {
		util.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if strings.TrimSpace(req.Status) == "" {
		util.ErrorJSON(w, http.StatusBadRequest, "status is required")
		return
	}

	employeeID := employeeIDFromClaims(r)
	if err := a.tracking.UpdateStatus(
		r.Context(),
		r.PathValue("trackId"),
		req.Status,
		req.Description,
		req.Location,
		employeeID,
	); err != nil {
		util.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	util.SuccessJSON(w, http.StatusOK, map[string]interface{}{
		"message": "status updated",
		"trackID": r.PathValue("trackId"),
		"status":  util.NormalizeStatus(req.Status),
	})
}

