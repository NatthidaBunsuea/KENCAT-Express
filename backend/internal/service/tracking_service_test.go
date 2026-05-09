// หมวย
package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"kencatexpress/backend/internal/domain"
	"kencatexpress/backend/internal/util"
)

// หมวย
func TestTrackingService_GetTracking(t *testing.T) {
	now := time.Now()
	svc := NewTrackingService(&testParcelStore{trackingResp: &domain.ParcelTrackingView{TrackID: "TRK-001", Status: util.StatusPending, UpdatedAt: now}})

	got, err := svc.GetTracking(context.Background(), " TRK-001 ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.TrackID != "TRK-001" {
		t.Fatalf("expected track id TRK-001, got %s", got.TrackID)
	}

	if _, err := svc.GetTracking(context.Background(), " "); err == nil || err.Error() != "track id is required" {
		t.Fatalf("expected track id error, got %v", err)
	}

	svc = NewTrackingService(&testParcelStore{trackingErr: sql.ErrNoRows})
	if _, err := svc.GetTracking(context.Background(), "TRK-404"); err == nil || err.Error() != "tracking not found" {
		t.Fatalf("expected tracking not found, got %v", err)
	}

	svc = NewTrackingService(&testParcelStore{trackingErr: errTestBoom})
	if _, err := svc.GetTracking(context.Background(), "TRK-500"); err == nil || err.Error() != errTestBoom.Error() {
		t.Fatalf("expected repository error, got %v", err)
	}
}

// หมวย
func TestTrackingService_UpdateStatus_SuccessInTransit(t *testing.T) {
	parcelStore := &testParcelStore{findResp: &domain.Parcel{TrackingCode: "TRK-001", Status: util.StatusPending}}
	svc := NewTrackingService(parcelStore)
	empID := int64(2)

	if err := svc.UpdateStatus(context.Background(), "TRK-001", util.StatusInTransit, "", "", &empID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if parcelStore.updateInput == nil {
		t.Fatal("expected update input to be captured")
	}
	if parcelStore.updateInput.Status != util.StatusInTransit {
		t.Fatalf("expected status %s, got %s", util.StatusInTransit, parcelStore.updateInput.Status)
	}
	if parcelStore.updateInput.Description != util.DefaultTrackingDescription(util.StatusInTransit) {
		t.Fatalf("expected default description, got %s", parcelStore.updateInput.Description)
	}
	if parcelStore.updateInput.Location != "On Route" {
		t.Fatalf("expected default location On Route, got %s", parcelStore.updateInput.Location)
	}
}

// หมวย
func TestTrackingService_UpdateStatus_SuccessDelivered(t *testing.T) {
	parcelStore := &testParcelStore{findResp: &domain.Parcel{TrackingCode: "TRK-002", Status: util.StatusInTransit}}
	svc := NewTrackingService(parcelStore)

	if err := svc.UpdateStatus(context.Background(), "TRK-002", util.StatusDelivered, "left at front desk", "Lobby", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if parcelStore.updateInput == nil || parcelStore.updateInput.DeliveredAt == nil {
		t.Fatal("expected DeliveredAt to be set")
	}
	if parcelStore.updateInput.Description != "left at front desk" || parcelStore.updateInput.Location != "Lobby" {
		t.Fatalf("expected custom description and location to be preserved, got %+v", parcelStore.updateInput)
	}
}

// หมวย
func TestTrackingService_UpdateStatus_SameStatusNoOp(t *testing.T) {
	parcelStore := &testParcelStore{findResp: &domain.Parcel{TrackingCode: "TRK-003", Status: util.StatusPending}}
	svc := NewTrackingService(parcelStore)

	if err := svc.UpdateStatus(context.Background(), "TRK-003", util.StatusPending, "", "", nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if parcelStore.updateInput != nil {
		t.Fatal("expected no update when status is unchanged")
	}
}

// หมวย
func TestTrackingService_UpdateStatus_Failures(t *testing.T) {
	svc := NewTrackingService(&testParcelStore{})

	if err := svc.UpdateStatus(context.Background(), "", util.StatusPending, "", "", nil); err == nil || err.Error() != "track id is required" {
		t.Fatalf("expected track id error, got %v", err)
	}
	if err := svc.UpdateStatus(context.Background(), "TRK-001", "", "", "", nil); err == nil || err.Error() != "status is required" {
		t.Fatalf("expected status error, got %v", err)
	}

	svc = NewTrackingService(&testParcelStore{findErr: sql.ErrNoRows})
	if err := svc.UpdateStatus(context.Background(), "TRK-404", util.StatusPending, "", "", nil); err == nil || err.Error() != "tracking not found" {
		t.Fatalf("expected tracking not found, got %v", err)
	}

	svc = NewTrackingService(&testParcelStore{findResp: &domain.Parcel{TrackingCode: "TRK-005", Status: util.StatusDelivered}})
	if err := svc.UpdateStatus(context.Background(), "TRK-005", util.StatusInTransit, "", "", nil); err == nil {
		t.Fatal("expected invalid transition error")
	}

	svc = NewTrackingService(&testParcelStore{findResp: &domain.Parcel{TrackingCode: "TRK-006", Status: util.StatusPending}, updateErr: errTestBoom})
	if err := svc.UpdateStatus(context.Background(), "TRK-006", util.StatusInTransit, "", "", nil); err == nil || err.Error() != errTestBoom.Error() {
		t.Fatalf("expected update error, got %v", err)
	}
}

// หมวย
func TestTrackingHelpers(t *testing.T) {
	if !isAllowedTransition(util.StatusPending, util.StatusInTransit) {
		t.Fatal("expected pending -> in transit to be allowed")
	}
	if isAllowedTransition(util.StatusDelivered, util.StatusPending) {
		t.Fatal("expected delivered -> pending to be disallowed")
	}
	if got := defaultLocation(util.StatusDelivered); got != "Destination" {
		t.Fatalf("expected Destination, got %s", got)
	}
}
