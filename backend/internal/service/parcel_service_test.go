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

// วุ่น
func TestParcelService_CreateParcel_Success(t *testing.T) {
	parcelStore := &testParcelStore{createResp: &domain.ParcelCreated{ParcelID: 10, ParcelCode: "PCL-000010", TrackingCode: "TRK-000010", ShippingCost: 65}}
	shippingSvc := NewShippingService(&testRateStore{rate: &domain.ShippingRate{ZoneCode: util.ZoneUpcountry, DeliveryType: util.DeliveryExpress, WeightMin: 0, WeightMax: 5, BasePrice: 65, ExtraPerKg: 5}})
	svc := NewParcelService(parcelStore, shippingSvc)
	clerkID := int64(12)

	resp, err := svc.CreateParcel(context.Background(), makeValidParcelRequest(), &clerkID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !resp.Success || resp.TrackID == "" || resp.ParcelID == "" {
		t.Fatalf("expected successful response with ids, got %+v", resp)
	}
	if parcelStore.createdInput == nil {
		t.Fatal("expected create input to be captured")
	}
	if parcelStore.createdInput.Status != util.StatusPending {
		t.Fatalf("expected status %s, got %s", util.StatusPending, parcelStore.createdInput.Status)
	}
	if parcelStore.createdInput.InitialEvent.Description != util.DefaultTrackingDescription(util.StatusPending) {
		t.Fatalf("expected default tracking description, got %s", parcelStore.createdInput.InitialEvent.Description)
	}
	if parcelStore.createdInput.ClerkEmployeeID == nil || *parcelStore.createdInput.ClerkEmployeeID != clerkID {
		t.Fatal("expected clerk id to be propagated")
	}
	if parcelStore.createdInput.Notes != "Express please" {
		t.Fatalf("expected trimmed notes, got %q", parcelStore.createdInput.Notes)
	}
}

// วุ่น
func TestParcelService_CreateParcel_Failures(t *testing.T) {
	shippingSvc := NewShippingService(&testRateStore{rate: &domain.ShippingRate{ZoneCode: util.ZoneUpcountry, DeliveryType: util.DeliveryExpress, WeightMin: 0, WeightMax: 5, BasePrice: 65}})
	svc := NewParcelService(&testParcelStore{}, shippingSvc)

	req := makeValidParcelRequest()
	req.Deliver.Name = ""
	if _, err := svc.CreateParcel(context.Background(), req, nil); err == nil || !strings.Contains(err.Error(), "sender name is required") {
		t.Fatalf("expected sender validation error, got %v", err)
	}

	req = makeValidParcelRequest()
	req.Receiver.Phone = "123"
	if _, err := svc.CreateParcel(context.Background(), req, nil); err == nil || !strings.Contains(err.Error(), "receiver phone must be 9-10 digits") {
		t.Fatalf("expected receiver phone validation error, got %v", err)
	}

	req = makeValidParcelRequest()
	req.Parcel.Weight = 0
	if _, err := svc.CreateParcel(context.Background(), req, nil); err == nil || err.Error() != "parcel weight must be greater than 0" {
		t.Fatalf("expected parcel weight error, got %v", err)
	}

	req = makeValidParcelRequest()
	req.Parcel.Type = "bad-type"
	if _, err := svc.CreateParcel(context.Background(), req, nil); err == nil || err.Error() != "delivery type must be STANDARD or EXPRESS" {
		t.Fatalf("expected delivery type error, got %v", err)
	}

	req = makeValidParcelRequest()
	svc = NewParcelService(&testParcelStore{createErr: errTestBoom}, shippingSvc)
	if _, err := svc.CreateParcel(context.Background(), req, nil); err == nil || err.Error() != errTestBoom.Error() {
		t.Fatalf("expected create error, got %v", err)
	}
}

// วุ่น
func TestParcelService_ListParcels(t *testing.T) {
	items := []domain.ParcelListItem{makeParcelListItem()}
	svc := NewParcelService(&testParcelStore{listResp: items}, nil)

	got, err := svc.ListParcels(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}

	svc = NewParcelService(&testParcelStore{listErr: errTestBoom}, nil)
	if _, err := svc.ListParcels(context.Background()); err == nil || err.Error() != errTestBoom.Error() {
		t.Fatalf("expected list error, got %v", err)
	}
}

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

// วุ่น
func TestValidatePerson_AndToPerson(t *testing.T) {
	person := dto.PersonAddressRequest{
		Name: " A ", Surname: " B ", Phone: "0812345678", Email: " a@example.com ",
		HomeNumber: "1", Soi: " 2 ", Road: " 3 ", District: "D", Subdistrict: "S", Province: "P", Zipcode: "10000",
	}
	if err := validatePerson(person, "sender"); err != nil {
		t.Fatalf("expected valid person, got %v", err)
	}
	converted := toPerson(person)
	if converted.FirstName != "A" || converted.LastName != "B" || converted.Soi != "2" || converted.Road != "3" {
		t.Fatalf("expected normalized person fields, got %+v", converted)
	}
}
