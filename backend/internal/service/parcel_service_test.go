// วุ่นกัส
package service

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"kencatexpress/backend/internal/domain"
	"kencatexpress/backend/internal/dto"
	"kencatexpress/backend/internal/util"
)

// กัส
func TestParcelService_GetParcelDetail(t *testing.T) {
	svc := NewParcelService(&testParcelStore{detailResp: &domain.ParcelDetail{Parcel: domain.Parcel{ParcelCode: "PCL-000001"}}}, nil)

	got, err := svc.GetParcelDetail(context.Background(), "  PCL-000001  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Parcel.ParcelCode != "PCL-000001" {
		t.Fatalf("expected parcel code PCL-000001, got %s", got.Parcel.ParcelCode)
	}

	if _, err := svc.GetParcelDetail(context.Background(), " "); err == nil || err.Error() != "parcel identifier is required" {
		t.Fatalf("expected parcel identifier error, got %v", err)
	}

	svc = NewParcelService(&testParcelStore{detailErr: sql.ErrNoRows}, nil)
	if _, err := svc.GetParcelDetail(context.Background(), "PCL-404"); err == nil || err.Error() != "parcel not found" {
		t.Fatalf("expected parcel not found, got %v", err)
	}

	svc = NewParcelService(&testParcelStore{detailErr: errTestBoom}, nil)
	if _, err := svc.GetParcelDetail(context.Background(), "PCL-500"); err == nil || err.Error() != errTestBoom.Error() {
		t.Fatalf("expected repository error, got %v", err)
	}
}
