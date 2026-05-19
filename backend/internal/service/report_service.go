// ยีน
package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"kencatexpress/backend/internal/domain"
	"kencatexpress/backend/internal/dto"
	"kencatexpress/backend/internal/util"
)

// ยีน
// หน้าที่: POST /api/reports, GET /api/reports
// Service นี้รับผิดชอบสร้างรายงานปัญหาและดึงรายการรายงาน
// Test ที่เกี่ยวข้อง: report_service_test.go
type ReportService struct {
	reports ReportStore
	parcels ParcelStore
}

// ยีน - สร้าง ReportService โดยรับ ReportStore/ParcelStore เพื่อ mock ใน test ได้
func NewReportService(reports ReportStore, parcels ParcelStore) *ReportService {
	return &ReportService{reports: reports, parcels: parcels}
}

// ยีน - สร้าง report โดยตรวจสอบ track id ก่อนว่ามีพัสดุจริง
func (s *ReportService) CreateReport(ctx context.Context, req dto.ReportCreateRequest) (*domain.Report, error) {
	if strings.TrimSpace(req.TrackID) == "" {
		return nil, errors.New("trackID is required")
	}
	if strings.TrimSpace(req.Issue) == "" {
		return nil, errors.New("issue is required")
	}
	if strings.TrimSpace(req.Reason) == "" {
		return nil, errors.New("reason is required")
	}
	parcel, err := s.parcels.FindParcelByTrackingCode(ctx, strings.TrimSpace(req.TrackID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("tracking not found")
		}
		return nil, err
	}
	report := &domain.Report{ReportCode: util.GenerateCode("RPT"), ParcelID: &parcel.ID, TrackingCode: strings.TrimSpace(req.TrackID), ReporterFirstName: util.NormalizeWhitespace(req.FirstName), ReporterLastName: util.NormalizeWhitespace(req.LastName), ReporterPhone: strings.TrimSpace(req.Phone), ReporterEmail: strings.TrimSpace(req.Email), IssueType: util.NormalizeWhitespace(req.Issue), Subject: util.NormalizeWhitespace(req.Issue), Description: util.NormalizeWhitespace(req.Reason), Status: "OPEN"}
	id, err := s.reports.CreateReport(ctx, report)
	if err != nil {
		return nil, err
	}
	report.ID = id
	return report, nil
}

// ยีน - ดึงรายการ report ทั้งหมด
func (s *ReportService) ListReports(ctx context.Context) ([]domain.Report, error) {
	return s.reports.ListReports(ctx)
}
