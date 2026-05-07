INSERT INTO employees (id, employee_code, email, password_hash, first_name, last_name, phone, role, birth_date, is_active)
VALUES
    (1, 'EMP001', 'clerk@kencat.local', 'eaba00ddd64d3dedb21016d7714dadc884f9a17a10f1bb3b22f30c8184bf9fae ', 'Aida', 'Bunsuea', '0812345678', 'PARCEL_CLERK', '2004-01-15', TRUE),
    (2, 'EMP002', 'messenger@kencat.local', 'eaba00ddd64d3dedb21016d7714dadc884f9a17a10f1bb3b22f30c8184bf9fae ', 'Pumipat', 'Santhongkaew', '0823456789', 'MESSENGER', '2004-04-10', TRUE),
    (3, 'EMP003', 'admin@kencat.local', 'eaba00ddd64d3dedb21016d7714dadc884f9a17a10f1bb3b22f30c8184bf9fae ', 'Admin', 'Kencat', '0834567890', 'ADMIN', '2003-11-01', TRUE)
ON DUPLICATE KEY UPDATE email=VALUES(email), password_hash=VALUES(password_hash), first_name=VALUES(first_name), last_name=VALUES(last_name), phone=VALUES(phone), role=VALUES(role), birth_date=VALUES(birth_date), is_active=VALUES(is_active);

INSERT INTO vehicles (id, vehicle_code, type, license_plate, status, assigned_employee_id)
VALUES
    (1, 'VAN-001', 'Van', '1กข-1234', 'IN_USE', 2),
    (2, 'BIKE-001', 'Motorbike', '2กข-5678', 'AVAILABLE', NULL),
    (3, 'TRUCK-001', 'Truck', '3กข-9012', 'AVAILABLE', NULL)
ON DUPLICATE KEY UPDATE type=VALUES(type), license_plate=VALUES(license_plate), status=VALUES(status), assigned_employee_id=VALUES(assigned_employee_id);

INSERT INTO shipping_rates (id, zone_code, delivery_type, weight_min, weight_max, base_price, extra_per_kg)
VALUES
    (1, 'BANGKOK_METRO', 'STANDARD', 0, 1, 23, 5), (2, 'BANGKOK_METRO', 'STANDARD', 1, 2, 28, 5), (3, 'BANGKOK_METRO', 'STANDARD', 2, 3, 39, 5), (4, 'BANGKOK_METRO', 'STANDARD', 3, 4, 49, 5), (5, 'BANGKOK_METRO', 'STANDARD', 4, 5, 60, 5), (6, 'BANGKOK_METRO', 'STANDARD', 5, 999, 65, 5),
    (7, 'BANGKOK_METRO', 'EXPRESS', 0, 1, 30, 5), (8, 'BANGKOK_METRO', 'EXPRESS', 1, 2, 34, 5), (9, 'BANGKOK_METRO', 'EXPRESS', 2, 3, 44, 5), (10, 'BANGKOK_METRO', 'EXPRESS', 3, 4, 54, 5), (11, 'BANGKOK_METRO', 'EXPRESS', 4, 5, 65, 5), (12, 'BANGKOK_METRO', 'EXPRESS', 5, 999, 70, 5),
    (13, 'UPCOUNTRY', 'STANDARD', 0, 1, 30, 5), (14, 'UPCOUNTRY', 'STANDARD', 1, 2, 39, 5), (15, 'UPCOUNTRY', 'STANDARD', 2, 3, 49, 5), (16, 'UPCOUNTRY', 'STANDARD', 3, 4, 60, 5), (17, 'UPCOUNTRY', 'STANDARD', 4, 5, 65, 5), (18, 'UPCOUNTRY', 'STANDARD', 5, 999, 70, 5),
    (19, 'UPCOUNTRY', 'EXPRESS', 0, 1, 35, 5), (20, 'UPCOUNTRY', 'EXPRESS', 1, 2, 44, 5), (21, 'UPCOUNTRY', 'EXPRESS', 2, 3, 54, 5), (22, 'UPCOUNTRY', 'EXPRESS', 3, 4, 65, 5), (23, 'UPCOUNTRY', 'EXPRESS', 4, 5, 70, 5), (24, 'UPCOUNTRY', 'EXPRESS', 5, 999, 75, 5)
ON DUPLICATE KEY UPDATE base_price=VALUES(base_price), extra_per_kg=VALUES(extra_per_kg);

INSERT INTO users (id, first_name, last_name, phone, email)
VALUES
    (1, 'Somchai', 'Sender', '0891111111', 'somchai@example.com'),
    (2, 'Suda', 'Receiver', '0892222222', 'suda@example.com'),
    (3, 'Narin', 'Customer', '0893333333', 'narin@example.com'),
    (4, 'Mali', 'Customer', '0894444444', 'mali@example.com')
ON DUPLICATE KEY UPDATE first_name=VALUES(first_name), last_name=VALUES(last_name), phone=VALUES(phone), email=VALUES(email);

INSERT INTO user_addresses (id, user_id, home_number, soi, road, subdistrict, district, province, zipcode, is_active)
VALUES
    (1, 1, '12/1', 'Soi A', 'Road A', 'Bang Kapi', 'Huai Khwang', 'Bangkok', '10310', TRUE),
    (2, 2, '55/7', 'Soi B', 'Road B', 'Talat Khwan', 'Mueang Nonthaburi', 'Nonthaburi', '11000', TRUE),
    (3, 3, '99/2', 'Soi C', 'Road C', 'Suthep', 'Mueang Chiang Mai', 'Chiang Mai', '50200', TRUE),
    (4, 4, '108', 'Soi D', 'Road D', 'Wichit', 'Mueang Phuket', 'Phuket', '83000', TRUE)
ON DUPLICATE KEY UPDATE user_id=VALUES(user_id), home_number=VALUES(home_number), soi=VALUES(soi), road=VALUES(road), subdistrict=VALUES(subdistrict), district=VALUES(district), province=VALUES(province), zipcode=VALUES(zipcode), is_active=VALUES(is_active);

INSERT INTO parcels (id, parcel_code, tracking_code, sender_user_id, sender_address_id, receiver_user_id, receiver_address_id, clerk_employee_id, messenger_employee_id, vehicle_id, delivery_type, weight, shipping_cost, origin_zone, destination_zone, status, notes, deposited_at, delivered_at)
VALUES
    (1, 'PCL-000001', 'TRK-000001', 1, 1, 2, 2, 1, 2, 1, 'STANDARD', 2.50, 39.00, 'BANGKOK_METRO', 'BANGKOK_METRO', 'In Transit', 'Handle with care', '2026-03-20 09:30:00', NULL),
    (2, 'PCL-000002', 'TRK-000002', 3, 3, 4, 4, 1, 2, 1, 'EXPRESS', 3.20, 65.00, 'UPCOUNTRY', 'UPCOUNTRY', 'Delivered', 'Customer requested express', '2026-03-18 13:45:00', '2026-03-21 16:05:00')
ON DUPLICATE KEY UPDATE sender_user_id=VALUES(sender_user_id), sender_address_id=VALUES(sender_address_id), receiver_user_id=VALUES(receiver_user_id), receiver_address_id=VALUES(receiver_address_id), clerk_employee_id=VALUES(clerk_employee_id), messenger_employee_id=VALUES(messenger_employee_id), vehicle_id=VALUES(vehicle_id), delivery_type=VALUES(delivery_type), weight=VALUES(weight), shipping_cost=VALUES(shipping_cost), origin_zone=VALUES(origin_zone), destination_zone=VALUES(destination_zone), status=VALUES(status), notes=VALUES(notes), deposited_at=VALUES(deposited_at), delivered_at=VALUES(delivered_at);

INSERT INTO tracking_events (id, parcel_id, tracking_code, status, location, description, updated_by_employee_id, created_at)
VALUES
    (1, 1, 'TRK-000001', 'Pending', 'Bangkok Hub', 'Parcel registered at parcel counter', 1, '2026-03-20 09:35:00'),
    (2, 1, 'TRK-000001', 'In Transit', 'Bangkok Hub', 'Messenger picked up parcel for delivery', 2, '2026-03-20 11:15:00'),
    (3, 2, 'TRK-000002', 'Pending', 'Chiang Mai Hub', 'Parcel registered at parcel counter', 1, '2026-03-18 13:50:00'),
    (4, 2, 'TRK-000002', 'In Transit', 'Chiang Mai Hub', 'Parcel left sorting center', 2, '2026-03-19 07:30:00'),
    (5, 2, 'TRK-000002', 'Delivered', 'Phuket Destination Hub', 'Parcel delivered to receiver successfully', 2, '2026-03-21 16:05:00')
ON DUPLICATE KEY UPDATE status=VALUES(status), location=VALUES(location), description=VALUES(description), updated_by_employee_id=VALUES(updated_by_employee_id), created_at=VALUES(created_at);

INSERT INTO vehicle_assignments (id, vehicle_id, parcel_id, messenger_employee_id, status, assigned_at, released_at)
VALUES
    (1, 1, 1, 2, 'ASSIGNED', '2026-03-20 10:50:00', NULL),
    (2, 1, 2, 2, 'COMPLETED', '2026-03-19 07:10:00', '2026-03-21 16:05:00')
ON DUPLICATE KEY UPDATE status=VALUES(status), assigned_at=VALUES(assigned_at), released_at=VALUES(released_at);

INSERT INTO reports (id, report_code, parcel_id, tracking_code, reporter_first_name, reporter_last_name, reporter_phone, reporter_email, issue_type, subject, description, status)
VALUES
    (1, 'RPT-000001', 2, 'TRK-000002', 'Mali', 'Customer', '0894444444', 'mali@example.com', 'Delivery Feedback', 'Late arrival concern', 'Customer wants to confirm why the delivery arrived later than expected.', 'OPEN')
ON DUPLICATE KEY UPDATE parcel_id=VALUES(parcel_id), tracking_code=VALUES(tracking_code), reporter_first_name=VALUES(reporter_first_name), reporter_last_name=VALUES(reporter_last_name), reporter_phone=VALUES(reporter_phone), reporter_email=VALUES(reporter_email), issue_type=VALUES(issue_type), subject=VALUES(subject), description=VALUES(description), status=VALUES(status);
