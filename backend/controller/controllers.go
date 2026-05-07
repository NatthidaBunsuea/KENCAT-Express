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
