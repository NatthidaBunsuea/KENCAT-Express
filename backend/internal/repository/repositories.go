// หมวย
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"kencatexpress/backend/internal/domain"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

// หมวย
func (s *Store) GetParcelTrackingView(ctx context.Context, trackID string) (*domain.ParcelTrackingView, error) {
	const q = `
SELECT p.tracking_code, p.parcel_code, p.status,
       su.first_name, su.last_name,
       ru.first_name, ru.last_name,
       ra.home_number, COALESCE(ra.soi, ''), COALESCE(ra.road, ''), ra.subdistrict, ra.district, ra.province, ra.zipcode,
       COALESCE(me.first_name, ''), COALESCE(me.last_name, ''),
       COALESCE(v.type, ''), COALESCE(v.license_plate, ''),
       p.delivered_at, p.updated_at
FROM parcels p
JOIN users su ON su.id = p.sender_user_id
JOIN users ru ON ru.id = p.receiver_user_id
JOIN user_addresses ra ON ra.id = p.receiver_address_id
LEFT JOIN employees me ON me.id = p.messenger_employee_id
LEFT JOIN vehicles v ON v.id = p.vehicle_id
WHERE p.tracking_code = ?
LIMIT 1`
	var view domain.ParcelTrackingView
	var senderFirst, senderLast, receiverFirst, receiverLast string
	var homeNo, soi, road, subdistrict, district, province, zipcode string
	var delivererFirst, delivererLast string
	var deliveredAt sql.NullTime
	if err := s.db.QueryRowContext(ctx, q, trackID).Scan(
		&view.TrackID,
		&view.ParcelID,
		&view.Status,
		&senderFirst,
		&senderLast,
		&receiverFirst,
		&receiverLast,
		&homeNo,
		&soi,
		&road,
		&subdistrict,
		&district,
		&province,
		&zipcode,
		&delivererFirst,
		&delivererLast,
		&view.TypeCar,
		&view.License,
		&deliveredAt,
		&view.UpdatedAt,
	); err != nil {
		return nil, err
	}
	view.Sender = strings.TrimSpace(senderFirst + " " + senderLast)
	view.Receiver = strings.TrimSpace(receiverFirst + " " + receiverLast)
	view.Deliverer = strings.TrimSpace(delivererFirst + " " + delivererLast)
	parts := []string{homeNo}
	if soi != "" {
		parts = append(parts, soi)
	}
	if road != "" {
		parts = append(parts, road)
	}
	parts = append(parts, subdistrict, district, province, zipcode)
	view.Address = strings.Join(parts, ", ")
	if deliveredAt.Valid {
		t := deliveredAt.Time
		view.DeliveredAt = &t
	}
	view.Events, _ = s.listTrackingEvents(ctx, view.TrackID)
	return &view, nil
}

// หมวย
func (s *Store) FindParcelByTrackingCode(ctx context.Context, trackID string) (*domain.Parcel, error) {
	const q = `
SELECT id, parcel_code, tracking_code, sender_user_id, sender_address_id, receiver_user_id, receiver_address_id,
       clerk_employee_id, messenger_employee_id, vehicle_id, delivery_type, weight, shipping_cost,
       origin_zone, destination_zone, status, COALESCE(notes, ''), deposited_at, delivered_at, created_at, updated_at
FROM parcels
WHERE tracking_code = ?
LIMIT 1`
	row := s.db.QueryRowContext(ctx, q, trackID)
	return scanParcel(row)
}

// หมวย
func (s *Store) UpdateParcelStatus(ctx context.Context, input domain.TrackingUpdateInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var parcelID int64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM parcels WHERE tracking_code = ? LIMIT 1`, input.TrackID).Scan(&parcelID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
UPDATE parcels
SET status = ?, messenger_employee_id = COALESCE(?, messenger_employee_id), delivered_at = CASE WHEN ? IS NULL THEN delivered_at ELSE ? END, updated_at = NOW()
WHERE tracking_code = ?`, input.Status, input.UpdatedByEmployeeID, input.DeliveredAt, input.DeliveredAt, input.TrackID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
INSERT INTO tracking_events (parcel_id, tracking_code, status, location, description, updated_by_employee_id, created_at)
VALUES (?, ?, ?, ?, ?, ?, NOW())`, parcelID, input.TrackID, input.Status, input.Location, input.Description, input.UpdatedByEmployeeID); err != nil {
		return err
	}
	return tx.Commit()
}
