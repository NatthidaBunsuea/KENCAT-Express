// ยีน
package service

import (
	"context"
	"database/sql"
	"testing"

	"kencatexpress/backend/internal/domain"
	"kencatexpress/backend/internal/dto"
)

// ยีน
func TestReportService_CreateReport_Success(t *testing.T) {
	parcelID := int64(42)
	reportStore := &testReportStore{createID: 101}
	parcelStore := &testParcelStore{findResp: &domain.Parcel{ID: parcelID, TrackingCode: "TRK-101"}}
	svc := NewReportService(reportStore, parcelStore)

	report, err := svc.CreateReport(context.Background(), dto.ReportCreateRequest{
		TrackID:   " TRK-101 ",
		Issue:     " Late Delivery ",
		Reason:    " Parcel arrived two days late ",
		FirstName: " Mali ",
		LastName:  " Customer ",
		Phone:     "0894444444",
		Email:     " mali@example.com ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if report.ID != 101 || report.Status != "OPEN" {
		t.Fatalf("expected report id 101 and status OPEN, got %+v", report)
	}
	if reportStore.created == nil {
		t.Fatal("expected created report to be captured")
	}
	if reportStore.created.TrackingCode != "TRK-101" || reportStore.created.IssueType != "Late Delivery" {
		t.Fatalf("expected normalized report fields, got %+v", reportStore.created)
	}
}

// ยีน
func TestReportService_CreateReport_Failures(t *testing.T) {
	svc := NewReportService(&testReportStore{}, &testParcelStore{})

	if _, err := svc.CreateReport(context.Background(), dto.ReportCreateRequest{}); err == nil || err.Error() != "trackID is required" {
		t.Fatalf("expected trackID required, got %v", err)
	}
	if _, err := svc.CreateReport(context.Background(), dto.ReportCreateRequest{TrackID: "TRK-1"}); err == nil || err.Error() != "issue is required" {
		t.Fatalf("expected issue required, got %v", err)
	}
	if _, err := svc.CreateReport(context.Background(), dto.ReportCreateRequest{TrackID: "TRK-1", Issue: "Late Delivery"}); err == nil || err.Error() != "reason is required" {
		t.Fatalf("expected reason required, got %v", err)
	}

	svc = NewReportService(&testReportStore{}, &testParcelStore{findErr: sql.ErrNoRows})
	if _, err := svc.CreateReport(context.Background(), dto.ReportCreateRequest{TrackID: "TRK-404", Issue: "Issue", Reason: "Reason"}); err == nil || err.Error() != "tracking not found" {
		t.Fatalf("expected tracking not found, got %v", err)
	}

	svc = NewReportService(&testReportStore{}, &testParcelStore{findErr: errTestBoom})
	if _, err := svc.CreateReport(context.Background(), dto.ReportCreateRequest{TrackID: "TRK-500", Issue: "Issue", Reason: "Reason"}); err == nil || err.Error() != errTestBoom.Error() {
		t.Fatalf("expected parcel lookup error, got %v", err)
	}

	svc = NewReportService(&testReportStore{createErr: errTestBoom}, &testParcelStore{findResp: &domain.Parcel{ID: 1, TrackingCode: "TRK-1"}})
	if _, err := svc.CreateReport(context.Background(), dto.ReportCreateRequest{TrackID: "TRK-1", Issue: "Issue", Reason: "Reason"}); err == nil || err.Error() != errTestBoom.Error() {
		t.Fatalf("expected create report error, got %v", err)
	}
}

// ยีน
func TestReportService_ListReports(t *testing.T) {
	expected := []domain.Report{{ID: 1, ReportCode: "RPT-001"}}
	svc := NewReportService(&testReportStore{listResp: expected}, &testParcelStore{})

	got, err := svc.ListReports(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != 1 || got[0].ReportCode != "RPT-001" {
		t.Fatalf("expected reports to be returned, got %+v", got)
	}

	svc = NewReportService(&testReportStore{listErr: errTestBoom}, &testParcelStore{})
	if _, err := svc.ListReports(context.Background()); err == nil || err.Error() != errTestBoom.Error() {
		t.Fatalf("expected list error, got %v", err)
	}
}
