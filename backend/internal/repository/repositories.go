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

// ด้า Find employee by employee code or email
func (s *Store) FindEmployeeByIdentity(ctx context.Context, identity string) (*domain.Employee, error) {
	const q = `
SELECT id, employee_code, email, password_hash, first_name, last_name, phone, role, birth_date, is_active, created_at, updated_at
FROM employees
WHERE employee_code = ? OR email = ?
LIMIT 1`
	row := s.db.QueryRowContext(ctx, q, identity, identity)
	return scanEmployee(row)
}

// ด้า Find employee by employee id
func (s *Store) FindEmployeeByID(ctx context.Context, id int64) (*domain.Employee, error) {
	const q = `
SELECT id, employee_code, email, password_hash, first_name, last_name, phone, role, birth_date, is_active, created_at, updated_at
FROM employees
WHERE id = ?
LIMIT 1`
	row := s.db.QueryRowContext(ctx, q, id)
	return scanEmployee(row)
}

// ด้า Find customer by user ID and load active address
func (s *Store) FindCustomerByID(ctx context.Context, id int64) (*domain.User, error) {
	const q = `
SELECT id, first_name, last_name, phone, COALESCE(email, ''), created_at, updated_at
FROM users
WHERE id = ?
LIMIT 1`
	var user domain.User
	if err := s.db.QueryRowContext(ctx, q, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}
	addr, err := s.FindActiveAddressByUserID(ctx, user.ID)
	if err == nil {
		user.Address = addr
	}
	return &user, nil
}

// ด้า Find the latest active address of a user
func (s *Store) FindActiveAddressByUserID(ctx context.Context, userID int64) (*domain.UserAddress, error) {
	const q = `
SELECT id, user_id, home_number, COALESCE(soi, ''), COALESCE(road, ''), subdistrict, district, province, zipcode, is_active, created_at, updated_at
FROM user_addresses
WHERE user_id = ? AND is_active = TRUE
ORDER BY id DESC
LIMIT 1`
	var addr domain.UserAddress
	if err := s.db.QueryRowContext(ctx, q, userID).Scan(
		&addr.ID,
		&addr.UserID,
		&addr.HomeNumber,
		&addr.Soi,
		&addr.Road,
		&addr.Subdistrict,
		&addr.District,
		&addr.Province,
		&addr.Zipcode,
		&addr.IsActive,
		&addr.CreatedAt,
		&addr.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &addr, nil
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

// วุ่น
func (s *Store) CreateParcelGraph(ctx context.Context, input domain.ParcelCreateInput) (*domain.ParcelCreated, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	senderID, senderAddrID, err := insertUserAndAddressTx(ctx, tx, input.Sender)
	if err != nil {
		return nil, err
	}
	receiverID, receiverAddrID, err := insertUserAndAddressTx(ctx, tx, input.Receiver)
	if err != nil {
		return nil, err
	}

	if input.DepositedAt.IsZero() {
		input.DepositedAt = time.Now()
	}
	res, err := tx.ExecContext(ctx, `
INSERT INTO parcels (
	parcel_code, tracking_code, sender_user_id, sender_address_id, receiver_user_id, receiver_address_id,
	clerk_employee_id, messenger_employee_id, vehicle_id, delivery_type, weight, shipping_cost,
	origin_zone, destination_zone, status, notes, deposited_at
) VALUES (?, ?, ?, ?, ?, ?, ?, NULL, NULL, ?, ?, ?, ?, ?, ?, ?, ?)`,
		input.ParcelCode,
		input.TrackingCode,
		senderID,
		senderAddrID,
		receiverID,
		receiverAddrID,
		input.ClerkEmployeeID,
		input.DeliveryType,
		input.Weight,
		input.ShippingCost,
		input.OriginZone,
		input.DestinationZone,
		input.Status,
		input.Notes,
		input.DepositedAt,
	)
	if err != nil {
		return nil, err
	}
	parcelID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	event := input.InitialEvent
	if event.Status == "" {
		event.Status = input.Status
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = input.DepositedAt
	}
	if _, err := tx.ExecContext(ctx, `
INSERT INTO tracking_events (parcel_id, tracking_code, status, location, description, updated_by_employee_id, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		parcelID,
		input.TrackingCode,
		event.Status,
		event.Location,
		event.Description,
		event.UpdatedByEmployeeID,
		event.CreatedAt,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &domain.ParcelCreated{ParcelID: parcelID, ParcelCode: input.ParcelCode, TrackingCode: input.TrackingCode, ShippingCost: input.ShippingCost}, nil
}

// วุ่น
func (s *Store) ListParcels(ctx context.Context) ([]domain.ParcelListItem, error) {
	const q = `
SELECT p.parcel_code, p.tracking_code,
       ru.first_name, ru.last_name, ru.phone,
       ra.home_number, COALESCE(ra.soi, ''), COALESCE(ra.road, ''), ra.subdistrict, ra.district, ra.province, ra.zipcode,
       p.status, p.deposited_at, p.messenger_employee_id
FROM parcels p
JOIN users ru ON ru.id = p.receiver_user_id
JOIN user_addresses ra ON ra.id = p.receiver_address_id
ORDER BY p.id DESC`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.ParcelListItem, 0)
	for rows.Next() {
		var item domain.ParcelListItem
		var messengerID sql.NullInt64
		if err := rows.Scan(
			&item.ParcelID,
			&item.TrackID,
			&item.ReceiverName,
			&item.ReceiverSurname,
			&item.ReceiverTel,
			&item.HomeNumber,
			&item.Soi,
			&item.Road,
			&item.Subdistrict,
			&item.DistrictName,
			&item.ProvinceName,
			&item.Zipcode,
			&item.Status,
			&item.DepositDate,
			&messengerID,
		); err != nil {
			return nil, err
		}
		if messengerID.Valid {
			id := messengerID.Int64
			item.AssignedMessengerID = &id
		}
		items = append(items, item)
	}
	return items, rows.Err()
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

// โอม
func (s *Store) ListMessengerTasks(ctx context.Context, employeeID int64) ([]domain.ParcelListItem, error) {
	query := `
SELECT p.parcel_code, p.tracking_code,
       ru.first_name, ru.last_name, ru.phone,
       ra.home_number, COALESCE(ra.soi, ''), COALESCE(ra.road, ''), ra.subdistrict, ra.district, ra.province, ra.zipcode,
       p.status, p.deposited_at, p.messenger_employee_id
FROM parcels p
JOIN users ru ON ru.id = p.receiver_user_id
JOIN user_addresses ra ON ra.id = p.receiver_address_id`
	args := []interface{}{}
	if employeeID > 0 {
		query += ` WHERE p.messenger_employee_id = ? OR p.messenger_employee_id IS NULL`
		args = append(args, employeeID)
	}
	query += ` ORDER BY p.id DESC`
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.ParcelListItem
	for rows.Next() {
		var item domain.ParcelListItem
		var messengerID sql.NullInt64
		if err := rows.Scan(&item.ParcelID, &item.TrackID, &item.ReceiverName, &item.ReceiverSurname, &item.ReceiverTel, &item.HomeNumber, &item.Soi, &item.Road, &item.Subdistrict, &item.DistrictName, &item.ProvinceName, &item.Zipcode, &item.Status, &item.DepositDate, &messengerID); err != nil {
			return nil, err
		}
		if messengerID.Valid {
			id := messengerID.Int64
			item.AssignedMessengerID = &id
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// โอม
func (s *Store) AssignVehicle(ctx context.Context, input domain.VehicleAssignmentInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var parcelID int64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM parcels WHERE tracking_code = ? LIMIT 1`, input.TrackID).Scan(&parcelID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE parcels SET vehicle_id = ?, messenger_employee_id = ?, updated_at = NOW() WHERE tracking_code = ?`, input.VehicleID, input.MessengerEmployeeID, input.TrackID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO vehicle_assignments (vehicle_id, parcel_id, messenger_employee_id, status, assigned_at) VALUES (?, ?, ?, 'ASSIGNED', NOW())`, input.VehicleID, parcelID, input.MessengerEmployeeID); err != nil {
		return err
	}
	return tx.Commit()
}

// โอม
func (s *Store) FindVehicleByID(ctx context.Context, id int64) (*domain.Vehicle, error) {
	const q = `
SELECT id, vehicle_code, type, license_plate, status, assigned_employee_id, created_at, updated_at
FROM vehicles
WHERE id = ?
LIMIT 1`
	var vehicle domain.Vehicle
	var assignedID sql.NullInt64
	if err := s.db.QueryRowContext(ctx, q, id).Scan(&vehicle.ID, &vehicle.VehicleCode, &vehicle.Type, &vehicle.LicensePlate, &vehicle.Status, &assignedID, &vehicle.CreatedAt, &vehicle.UpdatedAt); err != nil {
		return nil, err
	}
	if assignedID.Valid {
		id := assignedID.Int64
		vehicle.AssignedEmployeeID = &id
	}
	return &vehicle, nil
}

// โอม
func (s *Store) AssignVehicleToEmployee(ctx context.Context, vehicleID, employeeID int64) error {
	res, err := s.db.ExecContext(ctx, `UPDATE vehicles SET assigned_employee_id = ?, status = 'IN_USE', updated_at = NOW() WHERE id = ?`, employeeID, vehicleID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
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
