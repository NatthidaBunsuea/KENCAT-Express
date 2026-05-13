// กัส
package router

import (
	"net/http"
	"strings"

	"kencatexpress/backend/internal/config"
	"kencatexpress/backend/internal/controller"
	"kencatexpress/backend/internal/middleware"
	"kencatexpress/backend/internal/util"
)

func New(api *controller.API, cfg config.Config, frontendDir string) http.Handler {
	mux := http.NewServeMux()

	wrap := func(h http.Handler, auth ...middleware.Middleware) http.Handler {
		mws := []middleware.Middleware{middleware.Recoverer, middleware.RequestLogger, middleware.CORS(cfg)}
		mws = append(mws, auth...)
		return middleware.Chain(h, mws...)
	}

//ด้า
	mux.Handle("POST /api/auth/login", wrap(http.HandlerFunc(api.Login)))
	mux.Handle("GET /api/users/{userId}", wrap(http.HandlerFunc(api.GetUserByID), middleware.RequireAuth(cfg.JWTSecret)))
	
	
//กัส
	mux.Handle("GET /api/parcels/{parcelId}", wrap(http.HandlerFunc(api.GetParcel), middleware.RequireAuth(cfg.JWTSecret)))
	mux.Handle("GET /api/shipping/calculate", wrap(http.HandlerFunc(api.CalculateShipping)))
//หมวย
	mux.Handle("GET /api/trackings/{trackId}", wrap(http.HandlerFunc(api.GetTracking)))
	mux.Handle("PUT /api/trackings/{trackId}/status", wrap(http.HandlerFunc(api.UpdateTrackingStatus), middleware.RequireAuth(cfg.JWTSecret)))
//โอม
	mux.Handle("GET /api/messenger/tasks", wrap(http.HandlerFunc(api.GetMessengerTasks), middleware.RequireAuth(cfg.JWTSecret)))
	mux.Handle("POST /api/vehicle/assign", wrap(http.HandlerFunc(api.AssignVehicle), middleware.RequireAuth(cfg.JWTSecret)))
//กัส
	mux.HandleFunc("GET /docs/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/openapi.yaml")
	})
//กัส
	fileServer := http.FileServer(http.Dir(frontendDir))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api" || strings.HasPrefix(r.URL.Path, "/api/") {
			util.ErrorJSON(w, http.StatusNotFound, "endpoint not found")
			return
		}
		fileServer.ServeHTTP(w, r)
	}))
	return mux
}
