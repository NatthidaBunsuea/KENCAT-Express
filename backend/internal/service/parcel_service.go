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
