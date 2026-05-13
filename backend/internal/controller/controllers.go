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

// ด้า 1) POST /api/auth/login
func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := util.DecodeJSON(r, &req); err != nil {
		util.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := a.auth.Login(r.Context(), req)
	if err != nil {
		util.ErrorJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     util.CookieName,
		Value:    resp.Token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(a.tokenTTL.Seconds()),
	})

	util.WriteJSON(w, http.StatusOK, resp)
}

// ด้า 2) GET /api/users/{userId}
func (a *API) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("userId"), 10, 64)
	if err != nil || userID <= 0 {
		util.ErrorJSON(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := a.users.GetCustomerByID(r.Context(), userID)
	if err != nil {
		util.ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, user)
}


// กัส 5) GET /api/parcels/{parcelId}
func (a *API) GetParcel(w http.ResponseWriter, r *http.Request) {
	detail, err := a.parcels.GetParcelDetail(r.Context(), r.PathValue("parcelId"))
	if err != nil {
		util.ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, tracking)
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

// หมวย 7) GET /api/trackings/{trackId}
func (a *API) GetTracking(w http.ResponseWriter, r *http.Request) {
	tracking, err := a.tracking.GetTracking(r.Context(), r.PathValue("trackId"))

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

// โอม 9) GET /api/messenger/tasks
func (a *API) GetMessengerTasks(w http.ResponseWriter, r *http.Request) {
	var employeeID int64
	if claims, ok := middleware.ClaimsFromContext(r.Context()); ok {
		employeeID = claims.EmployeeID
	}

	items, err := a.messenger.ListTasks(r.Context(), employeeID)
	if err != nil {
		util.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, items)
}

// โอม 10) POST /api/vehicle/assign
func (a *API) AssignVehicle(w http.ResponseWriter, r *http.Request) {
	var req dto.VehicleAssignRequest
	if err := util.DecodeJSON(r, &req); err != nil {
		util.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.MessengerID == 0 {
		if claims, ok := middleware.ClaimsFromContext(r.Context()); ok {
			req.MessengerID = claims.EmployeeID
		}
	}

	if err := a.vehicles.AssignVehicle(r.Context(), req.TrackID, req.VehicleID, req.MessengerID); err != nil {
		util.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	util.SuccessJSON(w, http.StatusOK, map[string]interface{}{
		"message": "vehicle assigned",
		"trackID": req.TrackID,
	})
}
