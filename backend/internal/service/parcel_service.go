package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"kencatexpress/backend/internal/domain"
	"kencatexpress/backend/internal/dto"
	"kencatexpress/backend/internal/util"
)

// วุ่น
// หน้าที่: POST /api/parcels, GET /api/parcels
// Service นี้รับผิดชอบสร้างพัสดุและดึงรายการพัสดุ
// Test ที่เกี่ยวข้อง: parcel_service_test.go
type ParcelService struct {
	parcels  ParcelStore
	shipping *ShippingService
}

// วุ่น - สร้าง ParcelService โดยรับ ParcelStore และ ShippingService เพื่อ mock ใน test ได้
func NewParcelService(parcels ParcelStore, shipping *ShippingService) *ParcelService {
	return &ParcelService{parcels: parcels, shipping: shipping}
}

// วุ่น - สร้างพัสดุใหม่ พร้อมสร้าง tracking code และ event เริ่มต้น
func (s *ParcelService) CreateParcel(ctx context.Context, req dto.ParcelRequest, clerkID *int64) (*dto.ParcelCreateResponse, error) {
	if err := validatePerson(req.Deliver, "sender"); err != nil {
		return nil, err
	}
	if err := validatePerson(req.Receiver, "receiver"); err != nil {
		return nil, err
	}
	if req.Parcel.Weight <= 0 {
		return nil, errors.New("parcel weight must be greater than 0")
	}
	quote, err := s.shipping.Calculate(ctx, req.Deliver.Province, req.Receiver.Province, req.Parcel.Type, req.Parcel.Weight)
	if err != nil {
		return nil, err
	}
	input := domain.ParcelCreateInput{Sender: toPerson(req.Deliver), Receiver: toPerson(req.Receiver), ParcelCode: util.GenerateCode("PCL"), TrackingCode: util.GenerateCode("TRK"), ClerkEmployeeID: clerkID, DeliveryType: quote.DeliveryType, Weight: req.Parcel.Weight, ShippingCost: quote.ShippingCost, OriginZone: quote.OriginZone, DestinationZone: quote.DestinationZone, Status: util.StatusPending, Notes: util.NormalizeWhitespace(req.Parcel.Notes), DepositedAt: time.Now(), InitialEvent: domain.TrackingEvent{Status: util.StatusPending, Location: req.Deliver.Province, Description: util.DefaultTrackingDescription(util.StatusPending), CreatedAt: time.Now()}}
	if clerkID != nil {
		input.InitialEvent.UpdatedByEmployeeID = clerkID
	}
	created, err := s.parcels.CreateParcelGraph(ctx, input)
	if err != nil {
		return nil, err
	}
	return &dto.ParcelCreateResponse{Success: true, Message: "parcel created successfully", ParcelID: created.ParcelCode, TrackID: created.TrackingCode, ShippingCost: created.ShippingCost}, nil
}

// วุ่น - ดึงรายการพัสดุทั้งหมด
func (s *ParcelService) ListParcels(ctx context.Context) ([]domain.ParcelListItem, error) {
	return s.parcels.ListParcels(ctx)
}

// กัส
// หน้าที่: GET /api/parcels/{parcelId}
// อยู่ใน ParcelService เพราะข้อมูลมาจาก ParcelStore เหมือนกัน
func (s *ParcelService) GetParcelDetail(ctx context.Context, identifier string) (*domain.ParcelDetail, error) {
	if strings.TrimSpace(identifier) == "" {
		return nil, errors.New("parcel identifier is required")
	}
	detail, err := s.parcels.GetParcelDetail(ctx, strings.TrimSpace(identifier))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("parcel not found")
		}
		return nil, err
	}
	return detail, nil
}

// วุ่น - ตรวจข้อมูลผู้ส่ง/ผู้รับก่อนสร้างพัสดุ
func validatePerson(person dto.PersonAddressRequest, label string) error {
	fields := map[string]string{
		"name":        person.Name,
		"surname":     person.Surname,
		"phone":       person.Phone,
		"homeNumber":  person.HomeNumber,
		"district":    person.District,
		"subdistrict": person.Subdistrict,
		"province":    person.Province,
		"zipcode":     person.Zipcode,
	}
	for field, value := range fields {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s %s is required", label, field)
		}
	}
	if !util.IsValidPhone(person.Phone) {
		return fmt.Errorf("%s phone must be 9-10 digits", label)
	}
	return nil
}

// วุ่น - แปลง request ผู้ส่ง/ผู้รับให้เป็น domain model สำหรับบันทึกลงฐานข้อมูล
func toPerson(req dto.PersonAddressRequest) domain.PersonAddress {
	return domain.PersonAddress{
		FirstName:   util.NormalizeWhitespace(req.Name),
		LastName:    util.NormalizeWhitespace(req.Surname),
		Phone:       strings.TrimSpace(req.Phone),
		Email:       strings.TrimSpace(req.Email),
		HomeNumber:  util.NormalizeWhitespace(req.HomeNumber),
		Soi:         util.NormalizeWhitespace(req.Soi),
		Road:        util.NormalizeWhitespace(req.Road),
		District:    util.NormalizeWhitespace(req.District),
		Subdistrict: util.NormalizeWhitespace(req.Subdistrict),
		Province:    util.NormalizeWhitespace(req.Province),
		Zipcode:     strings.TrimSpace(req.Zipcode),
	}
}
