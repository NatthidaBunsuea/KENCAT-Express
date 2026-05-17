// หมวย
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"kencatexpress/backend/internal/domain"
	"kencatexpress/backend/internal/util"
)

// หมวย
// หน้าที่: GET /api/trackings/{trackId}, PUT /api/trackings/{trackId}/status
// Service นี้รับผิดชอบดูสถานะพัสดุและอัปเดตสถานะตามลำดับที่อนุญาต
// Test ที่เกี่ยวข้อง: tracking_service_test.go
type TrackingService struct{ parcels ParcelStore }

// หมวย - สร้าง TrackingService โดยรับ ParcelStore เพื่อ mock ใน test ได้
func NewTrackingService(parcels ParcelStore) *TrackingService {
	return &TrackingService{parcels: parcels}
}

// หมวย - ดึงข้อมูล tracking ตาม track id
func (s *TrackingService) GetTracking(ctx context.Context, trackID string) (*domain.ParcelTrackingView, error) {
	if strings.TrimSpace(trackID) == "" {
		return nil, errors.New("track id is required")
	}
	tracking, err := s.parcels.GetParcelTrackingView(ctx, strings.TrimSpace(trackID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("tracking not found")
		}
		return nil, err
	}
	return tracking, nil
}

// หมวย - อัปเดตสถานะ tracking พร้อมตรวจ transition ก่อนบันทึก
func (s *TrackingService) UpdateStatus(ctx context.Context, trackID, rawStatus, description, location string, employeeID *int64) error {
	if strings.TrimSpace(trackID) == "" {
		return errors.New("track id is required")
	}
	newStatus := util.NormalizeStatus(rawStatus)
	if newStatus == "" {
		return errors.New("status is required")
	}
	parcel, err := s.parcels.FindParcelByTrackingCode(ctx, strings.TrimSpace(trackID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("tracking not found")
		}
		return err
	}
	currentStatus := util.NormalizeStatus(parcel.Status)
	if currentStatus == newStatus {
		return nil
	}
	if !isAllowedTransition(currentStatus, newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", currentStatus, newStatus)
	}
	if strings.TrimSpace(description) == "" {
		description = util.DefaultTrackingDescription(newStatus)
	}
	if strings.TrimSpace(location) == "" {
		location = defaultLocation(newStatus)
	}
	var deliveredAt *time.Time
	if newStatus == util.StatusDelivered {
		now := time.Now()
		deliveredAt = &now
	}
	return s.parcels.UpdateParcelStatus(ctx, domain.TrackingUpdateInput{TrackID: strings.TrimSpace(trackID), Status: newStatus, Location: location, Description: description, UpdatedByEmployeeID: employeeID, DeliveredAt: deliveredAt})
}

// หมวย - ตรวจลำดับการเปลี่ยนสถานะ tracking ที่ระบบอนุญาต
func isAllowedTransition(currentStatus, newStatus string) bool {
	switch currentStatus {
	case util.StatusPending:
		return newStatus == util.StatusInTransit || newStatus == util.StatusDeliveryFailed
	case util.StatusInTransit:
		return newStatus == util.StatusDelivered || newStatus == util.StatusDeliveryFailed
	case util.StatusDelivered, util.StatusDeliveryFailed:
		return false
	default:
		return newStatus == util.StatusPending || newStatus == util.StatusInTransit || newStatus == util.StatusDelivered || newStatus == util.StatusDeliveryFailed
	}
}

// หมวย - กำหนด location เริ่มต้นเมื่อผู้ใช้ไม่ได้ส่ง location มากับการอัปเดตสถานะ
func defaultLocation(status string) string {
	switch util.NormalizeStatus(status) {
	case util.StatusPending:
		return "Parcel Counter"
	case util.StatusInTransit:
		return "On Route"
	case util.StatusDelivered, util.StatusDeliveryFailed:
		return "Destination"
	default:
		return "Kencat Express"
	}
}
