// โอม
package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"kencatexpress/backend/internal/domain"
)

// โอม
// หน้าที่: POST /api/vehicle/assign
// Service นี้รับผิดชอบตรวจสอบรถและ assign รถให้ messenger กับพัสดุ
// Test ที่เกี่ยวข้อง: vehicle_service_test.go
type VehicleService struct {
	parcels  ParcelStore
	vehicles VehicleStore
}

// โอม - สร้าง VehicleService โดยรับ ParcelStore/VehicleStore เพื่อ mock ใน test ได้
func NewVehicleService(parcels ParcelStore, vehicles VehicleStore) *VehicleService {
	return &VehicleService{parcels: parcels, vehicles: vehicles}
}

// โอม - assign รถให้ messenger และผูกกับ tracking
func (s *VehicleService) AssignVehicle(ctx context.Context, trackID string, vehicleID, messengerID int64) error {
	if strings.TrimSpace(trackID) == "" {
		return errors.New("track id is required")
	}
	if vehicleID <= 0 {
		return errors.New("vehicle id is required")
	}
	if messengerID <= 0 {
		return errors.New("messenger id is required")
	}
	vehicle, err := s.vehicles.FindVehicleByID(ctx, vehicleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("vehicle not found")
		}
		return err
	}
	if strings.ToUpper(strings.TrimSpace(vehicle.Status)) != "AVAILABLE" {
		if vehicle.AssignedEmployeeID == nil || *vehicle.AssignedEmployeeID != messengerID {
			return errors.New("vehicle is not available")
		}
	}
	if err := s.vehicles.AssignVehicleToEmployee(ctx, vehicleID, messengerID); err != nil {
		return err
	}
	return s.parcels.AssignVehicle(ctx, domain.VehicleAssignmentInput{TrackID: strings.TrimSpace(trackID), VehicleID: vehicleID, MessengerEmployeeID: messengerID})
}
