// หมวย
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"kencatexpress/backend/internal/config"
	"kencatexpress/backend/internal/controller"
	"kencatexpress/backend/internal/database"
	"kencatexpress/backend/internal/repository"
	"kencatexpress/backend/internal/router"
	"kencatexpress/backend/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := database.Bootstrap(ctx, db, cfg.AutoMigrate, cfg.AutoSeed); err != nil {
		log.Fatalf("database bootstrap failed: %v", err)
	}

	store := repository.New(db)
	authSvc := service.NewAuthService(store, cfg.JWTSecret, cfg.PasswordSalt, cfg.JWTTTL)
	userSvc := service.NewUserService(store, store)
	shippingSvc := service.NewShippingService(store)
	parcelSvc := service.NewParcelService(store, shippingSvc)
	trackingSvc := service.NewTrackingService(store)
	messengerSvc := service.NewMessengerService(store)
	vehicleSvc := service.NewVehicleService(store, store)
	reportSvc := service.NewReportService(store, store)

	api := controller.NewAPI(authSvc, userSvc, shippingSvc, parcelSvc, trackingSvc, messengerSvc, vehicleSvc, reportSvc, cfg.JWTTTL)

	frontendDir := cfg.FrontendDir
	if !filepath.IsAbs(frontendDir) {
		cwd, _ := os.Getwd()
		frontendDir = filepath.Join(cwd, frontendDir)
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router.New(api, cfg, frontendDir),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("KencatExpress API listening on http://localhost:%s", cfg.Port)
	log.Printf("Serving frontend from %s", frontendDir)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server stopped with error: %v", err)
	}
}
