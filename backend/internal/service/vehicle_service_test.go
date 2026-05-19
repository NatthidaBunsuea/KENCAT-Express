// โอม
package service

import (
	"context"
	"database/sql"
	"testing"

	"kencatexpress/backend/internal/domain"
)

// โอม
func TestVehicleService_AssignVehicle_SuccessAvailable(t *testing.T) {
	parcelStore := &testParcelStore{}
	vehicleStore := &testVehicleStore{vehicle: &domain.Vehicle{ID: 2, Status: "AVAILABLE"}}
	svc := NewVehicleService(parcelStore, vehicleStore)

	if err := svc.AssignVehicle(context.Background(), " TRK-001 ", 2, 8); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if vehicleStore.lastAssignVehID != 2 || vehicleStore.lastAssignEmpID != 8 {
		t.Fatalf("expected vehicle assignment to employee, got vehicle=%d employee=%d", vehicleStore.lastAssignVehID, vehicleStore.lastAssignEmpID)
	}
	if parcelStore.assignInput == nil || parcelStore.assignInput.TrackID != "TRK-001" {
		t.Fatalf("expected parcel assignment with trimmed track id, got %+v", parcelStore.assignInput)
	}
}

// โอม
func TestVehicleService_AssignVehicle_SuccessAlreadyAssignedToSameMessenger(t *testing.T) {
	parcelStore := &testParcelStore{}
	vehicleStore := &testVehicleStore{vehicle: &domain.Vehicle{ID: 3, Status: "IN_USE", AssignedEmployeeID: ptrInt64(9)}}
	svc := NewVehicleService(parcelStore, vehicleStore)

	if err := svc.AssignVehicle(context.Background(), "TRK-002", 3, 9); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// โอม
func TestVehicleService_AssignVehicle_Failures(t *testing.T) {
	svc := NewVehicleService(&testParcelStore{}, &testVehicleStore{})

	if err := svc.AssignVehicle(context.Background(), "", 1, 1); err == nil || err.Error() != "track id is required" {
		t.Fatalf("expected track id error, got %v", err)
	}
	if err := svc.AssignVehicle(context.Background(), "TRK-001", 0, 1); err == nil || err.Error() != "vehicle id is required" {
		t.Fatalf("expected vehicle id error, got %v", err)
	}
	if err := svc.AssignVehicle(context.Background(), "TRK-001", 1, 0); err == nil || err.Error() != "messenger id is required" {
		t.Fatalf("expected messenger id error, got %v", err)
	}

	svc = NewVehicleService(&testParcelStore{}, &testVehicleStore{findErr: sql.ErrNoRows})
	if err := svc.AssignVehicle(context.Background(), "TRK-001", 1, 1); err == nil || err.Error() != "vehicle not found" {
		t.Fatalf("expected vehicle not found, got %v", err)
	}

	svc = NewVehicleService(&testParcelStore{}, &testVehicleStore{findErr: errTestBoom})
	if err := svc.AssignVehicle(context.Background(), "TRK-001", 1, 1); err == nil || err.Error() != errTestBoom.Error() {
		t.Fatalf("expected repository error, got %v", err)
	}

	svc = NewVehicleService(&testParcelStore{}, &testVehicleStore{vehicle: &domain.Vehicle{ID: 1, Status: "IN_USE", AssignedEmployeeID: ptrInt64(77)}})
	if err := svc.AssignVehicle(context.Background(), "TRK-001", 1, 1); err == nil || err.Error() != "vehicle is not available" {
		t.Fatalf("expected vehicle unavailable error, got %v", err)
	}

	svc = NewVehicleService(&testParcelStore{}, &testVehicleStore{vehicle: &domain.Vehicle{ID: 1, Status: "AVAILABLE"}, assignErr: errTestBoom})
	if err := svc.AssignVehicle(context.Background(), "TRK-001", 1, 1); err == nil || err.Error() != errTestBoom.Error() {
		t.Fatalf("expected assign to employee error, got %v", err)
	}

	svc = NewVehicleService(&testParcelStore{assignErr: errTestBoom}, &testVehicleStore{vehicle: &domain.Vehicle{ID: 1, Status: "AVAILABLE"}})
	if err := svc.AssignVehicle(context.Background(), "TRK-001", 1, 1); err == nil || err.Error() != errTestBoom.Error() {
		t.Fatalf("expected parcel assign error, got %v", err)
	}
}
