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


// กัส 5) GET /api/parcels/{parcelId}
func (a *API) GetParcel(w http.ResponseWriter, r *http.Request) {
	detail, err := a.parcels.GetParcelDetail(r.Context(), r.PathValue("parcelId"))
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
	util.WriteJSON(w, http.StatusOK, detail)
}

// กัส 6) GET /api/shipping/calculate
func (a *API) CalculateShipping(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	deliveryType := q.Get("deliveryType")
	if deliveryType == "" {
		deliveryType = q.Get("type")
	}

	weight, err := strconv.ParseFloat(q.Get("weight"), 64)
	if err != nil {
		util.ErrorJSON(w, http.StatusBadRequest, "invalid weight")
		return
	}

	quote, err := a.shipping.Calculate(
		r.Context(),
		q.Get("origin"),
		q.Get("destination"),
		deliveryType,
		weight,
	)
	if err != nil {
		util.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	util.SuccessJSON(w, http.StatusOK, map[string]interface{}{
		"message": "status updated",
		"trackID": r.PathValue("trackId"),
		"status":  util.NormalizeStatus(req.Status),
	})
}

	util.WriteJSON(w, http.StatusOK, quote)
}
