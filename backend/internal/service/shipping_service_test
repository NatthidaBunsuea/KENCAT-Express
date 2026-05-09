// กัส
package service

import (
	"context"
	"database/sql"
	"testing"

	"kencatexpress/backend/internal/domain"
	"kencatexpress/backend/internal/util"
)

// กัส
func TestShippingService_Calculate_Success(t *testing.T) {
	rateStore := &testRateStore{rate: &domain.ShippingRate{
		ZoneCode:     util.ZoneBangkokMetro,
		DeliveryType: util.DeliveryStandard,
		WeightMin:    0,
		WeightMax:    3,
		BasePrice:    39,
		ExtraPerKg:   5,
	}}
	svc := NewShippingService(rateStore)

	resp, err := svc.Calculate(context.Background(), "Bangkok", "Bangkok", util.DeliveryStandard, 2.5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ShippingCost != 39 {
		t.Fatalf("expected shipping cost 39, got %v", resp.ShippingCost)
	}
	if resp.OriginZone != util.ZoneBangkokMetro || resp.DestinationZone != util.ZoneBangkokMetro {
		t.Fatalf("expected bangkok metro zones, got %+v", resp)
	}
	if rateStore.lastDeliveryType != util.DeliveryStandard {
		t.Fatalf("expected normalized delivery type %s, got %s", util.DeliveryStandard, rateStore.lastDeliveryType)
	}
}

// กัส
func TestShippingService_Calculate_Failures(t *testing.T) {
	svc := NewShippingService(&testRateStore{})

	if _, err := svc.Calculate(context.Background(), "Bangkok", "Bangkok", util.DeliveryStandard, 0); err == nil || err.Error() != "weight must be greater than 0" {
		t.Fatalf("expected invalid weight error, got %v", err)
	}

	if _, err := svc.Calculate(context.Background(), "Bangkok", "Bangkok", "same-day", 1); err == nil || err.Error() != "delivery type must be STANDARD or EXPRESS" {
		t.Fatalf("expected invalid delivery type error, got %v", err)
	}
}

// กัส
func TestShippingService_CalculateByZone_SuccessWithExtraCharge(t *testing.T) {
	svc := NewShippingService(&testRateStore{rate: &domain.ShippingRate{
		ZoneCode:     util.ZoneUpcountry,
		DeliveryType: util.DeliveryExpress,
		WeightMin:    0,
		WeightMax:    5,
		BasePrice:    50,
		ExtraPerKg:   7.5,
	}})

	cost, err := svc.CalculateByZone(context.Background(), util.ZoneUpcountry, util.DeliveryExpress, 7)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cost != 65 {
		t.Fatalf("expected cost 65, got %v", cost)
	}
}

// กัส
func TestShippingService_CalculateByZone_Failures(t *testing.T) {
	svc := NewShippingService(&testRateStore{err: sql.ErrNoRows})
	if _, err := svc.CalculateByZone(context.Background(), util.ZoneUpcountry, util.DeliveryExpress, 1); err == nil || err.Error() != "shipping rate not found" {
		t.Fatalf("expected shipping rate not found, got %v", err)
	}

	svc = NewShippingService(&testRateStore{err: errTestBoom})
	if _, err := svc.CalculateByZone(context.Background(), util.ZoneUpcountry, util.DeliveryExpress, 1); err == nil || err.Error() != errTestBoom.Error() {
		t.Fatalf("expected repository error, got %v", err)
	}
}
