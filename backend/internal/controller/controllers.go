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



// กัส 5) GET /api/parcels/{parcelId}
func (a *API) GetParcel(w http.ResponseWriter, r *http.Request) {
	detail, err := a.parcels.GetParcelDetail(r.Context(), r.PathValue("parcelId"))
	if err != nil {
		util.ErrorJSON(w, http.StatusNotFound, err.Error())
		return
	}

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

	util.WriteJSON(w, http.StatusOK, quote)
}
