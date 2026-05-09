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


// กัส
func (s *Store) FindMatchedRate(ctx context.Context, zoneCode, deliveryType string, weight float64) (*domain.ShippingRate, error) {
	const q = `
SELECT id, zone_code, delivery_type, weight_min, weight_max, base_price, extra_per_kg
FROM shipping_rates
WHERE zone_code = ?
  AND delivery_type = ?
  AND ? > weight_min
  AND ? <= weight_max
ORDER BY weight_max ASC
LIMIT 1`
	var rate domain.ShippingRate
	err := s.db.QueryRowContext(ctx, q, zoneCode, deliveryType, weight, weight).Scan(
		&rate.ID,
		&rate.ZoneCode,
		&rate.DeliveryType,
		&rate.WeightMin,
		&rate.WeightMax,
		&rate.BasePrice,
		&rate.ExtraPerKg,
	)
	if err == sql.ErrNoRows {
		fallback := `
SELECT id, zone_code, delivery_type, weight_min, weight_max, base_price, extra_per_kg
FROM shipping_rates
WHERE zone_code = ?
  AND delivery_type = ?
ORDER BY weight_max DESC
LIMIT 1`
		err = s.db.QueryRowContext(ctx, fallback, zoneCode, deliveryType).Scan(
			&rate.ID,
			&rate.ZoneCode,
			&rate.DeliveryType,
			&rate.WeightMin,
			&rate.WeightMax,
			&rate.BasePrice,
			&rate.ExtraPerKg,
		)
	}
	if err != nil {
		return nil, err
	}
	return &rate, nil
}


// กัส
func (s *Store) GetParcelDetail(ctx context.Context, identifier string) (*domain.ParcelDetail, error) {
	const q = `
SELECT p.id, p.parcel_code, p.tracking_code, p.sender_user_id, p.sender_address_id, p.receiver_user_id, p.receiver_address_id,
       p.clerk_employee_id, p.messenger_employee_id, p.vehicle_id, p.delivery_type, p.weight, p.shipping_cost,
       p.origin_zone, p.destination_zone, p.status, COALESCE(p.notes, ''), p.deposited_at, p.delivered_at, p.created_at, p.updated_at,
       su.id, su.first_name, su.last_name, su.phone, COALESCE(su.email, ''),
       sa.id, sa.home_number, COALESCE(sa.soi, ''), COALESCE(sa.road, ''), sa.subdistrict, sa.district, sa.province, sa.zipcode,
       ru.id, ru.first_name, ru.last_name, ru.phone, COALESCE(ru.email, ''),
       ra.id, ra.home_number, COALESCE(ra.soi, ''), COALESCE(ra.road, ''), ra.subdistrict, ra.district, ra.province, ra.zipcode
FROM parcels p
JOIN users su ON su.id = p.sender_user_id
JOIN user_addresses sa ON sa.id = p.sender_address_id
JOIN users ru ON ru.id = p.receiver_user_id
JOIN user_addresses ra ON ra.id = p.receiver_address_id
WHERE p.parcel_code = ? OR CAST(p.id AS CHAR) = ?
LIMIT 1`
	var detail domain.ParcelDetail
	var clerkID, messengerID, vehicleID sql.NullInt64
	var deliveredAt sql.NullTime
	var sender, receiver domain.User
	var senderAddr, receiverAddr domain.UserAddress

	err := s.db.QueryRowContext(ctx, q, identifier, identifier).Scan(
		&detail.Parcel.ID,
		&detail.Parcel.ParcelCode,
		&detail.Parcel.TrackingCode,
		&detail.Parcel.SenderUserID,
		&detail.Parcel.SenderAddressID,
		&detail.Parcel.ReceiverUserID,
		&detail.Parcel.ReceiverAddressID,
		&clerkID,
		&messengerID,
		&vehicleID,
		&detail.Parcel.DeliveryType,
		&detail.Parcel.Weight,
		&detail.Parcel.ShippingCost,
		&detail.Parcel.OriginZone,
		&detail.Parcel.DestinationZone,
		&detail.Parcel.Status,
		&detail.Parcel.Notes,
		&detail.Parcel.DepositedAt,
		&deliveredAt,
		&detail.Parcel.CreatedAt,
		&detail.Parcel.UpdatedAt,
		&sender.ID,
		&sender.FirstName,
		&sender.LastName,
		&sender.Phone,
		&sender.Email,
		&senderAddr.ID,
		&senderAddr.HomeNumber,
		&senderAddr.Soi,
		&senderAddr.Road,
		&senderAddr.Subdistrict,
		&senderAddr.District,
		&senderAddr.Province,
		&senderAddr.Zipcode,
		&receiver.ID,
		&receiver.FirstName,
		&receiver.LastName,
		&receiver.Phone,
		&receiver.Email,
		&receiverAddr.ID,
		&receiverAddr.HomeNumber,
		&receiverAddr.Soi,
		&receiverAddr.Road,
		&receiverAddr.Subdistrict,
		&receiverAddr.District,
		&receiverAddr.Province,
		&receiverAddr.Zipcode,
	)
	if err != nil {
		return nil, err
	}
	if clerkID.Valid {
		id := clerkID.Int64
		detail.Parcel.ClerkEmployeeID = &id
	}
	if messengerID.Valid {
		id := messengerID.Int64
		detail.Parcel.MessengerEmployeeID = &id
	}
	if vehicleID.Valid {
		id := vehicleID.Int64
		detail.Parcel.VehicleID = &id
	}
	if deliveredAt.Valid {
		t := deliveredAt.Time
		detail.Parcel.DeliveredAt = &t
	}
	sender.Address = &senderAddr
	receiver.Address = &receiverAddr
	detail.Sender = sender
	detail.SenderAddress = senderAddr
	detail.Receiver = receiver
	detail.ReceiverAddress = receiverAddr
	detail.Events, _ = s.listTrackingEvents(ctx, detail.Parcel.TrackingCode)
	return &detail, nil
}


// กัส
func insertUserAndAddressTx(ctx context.Context, tx *sql.Tx, person domain.PersonAddress) (int64, int64, error) {
	userRes, err := tx.ExecContext(ctx, `INSERT INTO users (first_name, last_name, phone, email) VALUES (?, ?, ?, ?)`, person.FirstName, person.LastName, person.Phone, nullString(person.Email))
	if err != nil {
		return 0, 0, err
	}
	userID, err := userRes.LastInsertId()
	if err != nil {
		return 0, 0, err
	}
	addrRes, err := tx.ExecContext(ctx, `INSERT INTO user_addresses (user_id, home_number, soi, road, subdistrict, district, province, zipcode, is_active) VALUES (?, ?, ?, ?, ?, ?, ?, ?, TRUE)`, userID, person.HomeNumber, nullString(person.Soi), nullString(person.Road), person.Subdistrict, person.District, person.Province, person.Zipcode)
	if err != nil {
		return 0, 0, err
	}
	addrID, err := addrRes.LastInsertId()
	if err != nil {
		return 0, 0, err
	}
	return userID, addrID, nil
}

// กัส
func (s *Store) listTrackingEvents(ctx context.Context, trackID string) ([]domain.TrackingEvent, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, parcel_id, tracking_code, status, COALESCE(location, ''), COALESCE(description, ''), updated_by_employee_id, created_at FROM tracking_events WHERE tracking_code = ? ORDER BY created_at ASC`, trackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []domain.TrackingEvent
	for rows.Next() {
		var event domain.TrackingEvent
		var updatedBy sql.NullInt64
		if err := rows.Scan(&event.ID, &event.ParcelID, &event.TrackingCode, &event.Status, &event.Location, &event.Description, &updatedBy, &event.CreatedAt); err != nil {
			return nil, err
		}
		if updatedBy.Valid {
			id := updatedBy.Int64
			event.UpdatedByEmployeeID = &id
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

// กัส
func scanEmployee(row *sql.Row) (*domain.Employee, error) {
	var employee domain.Employee
	var phone sql.NullString
	var birthDate sql.NullTime
	if err := row.Scan(&employee.ID, &employee.EmployeeCode, &employee.Email, &employee.PasswordHash, &employee.FirstName, &employee.LastName, &phone, &employee.Role, &birthDate, &employee.IsActive, &employee.CreatedAt, &employee.UpdatedAt); err != nil {
		return nil, err
	}
	if phone.Valid {
		employee.Phone = phone.String
	}
	if birthDate.Valid {
		t := birthDate.Time
		employee.BirthDate = &t
	}
	return &employee, nil
}

// กัส
func scanParcel(row *sql.Row) (*domain.Parcel, error) {
	var parcel domain.Parcel
	var clerkID, messengerID, vehicleID sql.NullInt64
	var deliveredAt sql.NullTime
	if err := row.Scan(&parcel.ID, &parcel.ParcelCode, &parcel.TrackingCode, &parcel.SenderUserID, &parcel.SenderAddressID, &parcel.ReceiverUserID, &parcel.ReceiverAddressID, &clerkID, &messengerID, &vehicleID, &parcel.DeliveryType, &parcel.Weight, &parcel.ShippingCost, &parcel.OriginZone, &parcel.DestinationZone, &parcel.Status, &parcel.Notes, &parcel.DepositedAt, &deliveredAt, &parcel.CreatedAt, &parcel.UpdatedAt); err != nil {
		return nil, err
	}
	if clerkID.Valid {
		id := clerkID.Int64
		parcel.ClerkEmployeeID = &id
	}
	if messengerID.Valid {
		id := messengerID.Int64
		parcel.MessengerEmployeeID = &id
	}
	if vehicleID.Valid {
		id := vehicleID.Int64
		parcel.VehicleID = &id
	}
	if deliveredAt.Valid {
		t := deliveredAt.Time
		parcel.DeliveredAt = &t
	}
	return &parcel, nil
}

// กัส
func nullString(v string) interface{} {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return v
}

// กัส
func MustNotFound(err error, message string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf(message)
	}
	return err
}
