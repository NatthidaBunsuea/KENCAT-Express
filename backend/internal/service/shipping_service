// กัส
package service

import (
	"context"
	"database/sql"
	"errors"

	"kencatexpress/backend/internal/dto"
	"kencatexpress/backend/internal/util"
)

// กัส
// หน้าที่: GET /api/shipping/calculate
// Service นี้รับผิดชอบคำนวณค่าส่งจากโซน ประเภทการส่ง และน้ำหนัก
// Test ที่เกี่ยวข้อง: shipping_service_test.go
type ShippingService struct{ rates ShippingRateStore }

// กัส - สร้าง ShippingService โดยรับ rate repository เพื่อ mock ใน test ได้
func NewShippingService(rates ShippingRateStore) *ShippingService {
	return &ShippingService{rates: rates}
}

// กัส - คำนวณค่าส่งจากจังหวัดต้นทาง/ปลายทาง
func (s *ShippingService) Calculate(ctx context.Context, originProvince, destinationProvince, deliveryType string, weight float64) (*dto.ShippingCalculationResponse, error) {
	if weight <= 0 {
		return nil, errors.New("weight must be greater than 0")
	}
	dType := util.NormalizeDeliveryType(deliveryType)
	if dType != util.DeliveryStandard && dType != util.DeliveryExpress {
		return nil, errors.New("delivery type must be STANDARD or EXPRESS")
	}
	originZone := util.ZoneFromProvince(originProvince)
	destinationZone := util.ZoneFromProvince(destinationProvince)
	cost, err := s.CalculateByZone(ctx, destinationZone, dType, weight)
	if err != nil {
		return nil, err
	}
	return &dto.ShippingCalculationResponse{OriginZone: originZone, DestinationZone: destinationZone, DeliveryType: dType, Weight: weight, ShippingCost: cost}, nil
}

// กัส - คำนวณค่าส่งจาก zone โดยตรง ใช้ร่วมกับ ParcelService ตอนสร้างพัสดุ
func (s *ShippingService) CalculateByZone(ctx context.Context, destinationZone, deliveryType string, weight float64) (float64, error) {
	rate, err := s.rates.FindMatchedRate(ctx, destinationZone, util.NormalizeDeliveryType(deliveryType), weight)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("shipping rate not found")
		}
		return 0, err
	}
	cost := rate.BasePrice
	if weight > rate.WeightMax {
		cost += (weight - rate.WeightMax) * rate.ExtraPerKg
	}
	return util.RoundCurrency(cost), nil
}
