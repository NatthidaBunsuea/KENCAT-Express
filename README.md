# KENCAT-Express
# KencatExpress API Test Guide

README นี้เป็นเวอร์ชันสรุปการติดตั้งและตัวอย่างการทดสอบ API ที่ตรงกับโปรเจกต์ KencatExpress สำหรับ 12 เส้นที่ต้องยื่น

## 1. Project Overview

KencatExpress เป็นระบบจัดการพัสดุ พัฒนาโดยใช้

- Backend: Go (Golang)
- Database: MySQL
- Architecture: Layered Architecture (Controller -> Service -> Repository)

## โครงสร้างหลัก

- `Frontend/` หน้าเว็บเดิมจากโปรเจกต์ต้นฉบับ
- `backend/cmd/api` จุดเริ่มต้นของเซิร์ฟเวอร์ Go
- `backend/internal/controller` controller layer
- `backend/internal/service` service layer
- `backend/internal/repository` repository layer
- `backend/internal/database` database bootstrap
- `backend/internal/middleware` auth / cors / logging / recover
- `backend/migrations` schema + seed สำหรับ MySQL
- `backend/docs/openapi.yaml` เอกสาร API แบบ OpenAPI

## 2. System Requirements

ต้องติดตั้งก่อนใช้งาน

- Go 1.21 หรือใหม่กว่า
- MySQL 8.0 หรือใหม่กว่า
- PowerShell / Terminal
- Postman (ถ้าต้องการใช้ทดสอบแบบ GUI)

## 3. Project Setup

### 3.1 แตกไฟล์โปรเจกต์

```bash
unzip KencatExpress.zip
cd KencatExpress
```

### 3.2 สร้างฐานข้อมูล

เข้า MySQL แล้วสร้างฐานข้อมูล

```bash
mysql -u root -p
```

```sql
CREATE DATABASE kencat;
```

### 3.3 Import Schema และ Seed

เปิด Terminal เข้าไปในโฟลเดอร์ backend แล้ว import ไฟล์ migration

```bash
cd backend
mysql -u root -p kencat < migrations/001_schema.sql
mysql -u root -p kencat < migrations/002_seed.sql
```

เช็คว่าตารางถูกสร้างแล้ว
```sql
USE kencat;
SHOW TABLES;
```

### 3.4 ตั้งค่า Database

กำหนดค่าการเชื่อมต่อฐานข้อมูลให้ตรงกับเครื่องที่ใช้งาน ใน `.env` 

```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=*yourpassword* เปลี่ยนตามของตัวเองในโค้ดใส่ไว้เป็น Muaymin_ly12548
DB_NAME=kencat
APP_PORT=8080
JWT_SECRET=secret
PASSWORD_SALT=kencat-express-salt
```

### 3.5 Run Backend

กลับไปใน Terminal ใหม่
```bash
cd KencatExpress
cd backend
go mod init kencatexpress/backend
go mod tidy
go run ./cmd/api/main.go
```

เมื่อรันสำเร็จ server จะทำงานที่

```text
http://localhost:8080
```

## 4. Authentication Flow

ลำดับการใช้งาน API ที่มีการป้องกันสิทธิ์

```text
Login -> Receive Token -> Send Token in Header -> Access Protected APIs
```

ตัวอย่าง Header

```http
Authorization: Bearer <token>
```

## 5. API Endpoints (12 Required Routes)

รายการ API ที่ต้องยื่น

- POST /api/auth/login (ด้า)
- GET /api/users/{userId} (ด้า)
- POST /api/parcels (วุ่น)
- GET /api/parcels (วุ่น)
- GET /api/parcels/{parcelId} (กัส)
- GET /api/trackings/{trackId} (หมวย)
- GET /api/shipping/calculate (กัส)
- PUT /api/trackings/{trackId}/status (หมวย)
- GET /api/messenger/tasks (โอม)
- POST /api/vehicle/assign (โอม)
- POST /api/reports (ยีน)
- GET /api/reports (ยีน)

## 6. API Testing Examples

Base URL ที่ใช้ในตัวอย่างทั้งหมด

```text
http://localhost:8080
```

### 6.1 POST /api/auth/login

ใช้สำหรับเข้าสู่ระบบและรับ token

**Request**

```http
POST /api/auth/login
Content-Type: application/json
```

```json
{
  "employeeID": "EMP001",
  "password": "1234"
}
```

**PowerShell**

```powershell
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/auth/login" -ContentType "application/json" -Body '{"employeeID":"EMP001","password":"1234"}'
```

หมายเหตุ:
- ถ้าจะทดสอบ messenger tasks แนะนำให้ login ด้วย `EMP002`
- แรกเริ่มรหัสของ EMP001 จะถูก salt ไว้ให้ไปเปลี่ยนใน DB ก่อน เช็ครหัสต้องตรงกัน กล่าวคือไม่ต้องเปลี่ยนไป login ของ EMP002 แล้วเพราะรหัสเดียวกันหมด
```sql
USE kencat;

UPDATE employees
SET password_hash = SHA2(CONCAT('kencat-express-salt', '1234'), 256)
WHERE employee_code IN ('EMP001', 'EMP002', 'EMP003');

SELECT employee_code, password_hash
FROM employees
WHERE employee_code IN ('EMP001', 'EMP002', 'EMP003');
```
---

### 6.2 GET /api/users/{userId}

ใช้สำหรับดูข้อมูลผู้ใช้ตาม user id เมื่อ login 6.1 ได้แล้วจะได้ token มา ให้ใส่ token นั้นใน Auth แบบ  Bearer <token>

**Request**

```http
GET /api/users/3
Authorization: Bearer <token> 
```

**PowerShell**

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/users/3" -Headers @{ Authorization = "Bearer <token>" }
```

---

### 6.3 POST /api/parcels

ใช้สำหรับสร้างพัสดุใหม่

**Request**

```http
POST /api/parcels
Content-Type: application/json
```

```json
{
  "deliver": {
    "name": "Narin",
    "surname": "Customer",
    "phone": "0893333333",
    "email": "narin@example.com",
    "homeNumber": "99/2",
    "soi": "Soi C",
    "road": "Road C",
    "district": "Mueang Chiang Mai",
    "subdistrict": "Suthep",
    "province": "Chiang Mai",
    "zipcode": "50200"
  },
  "receiver": {
    "name": "Mali",
    "surname": "Customer",
    "phone": "0894444444",
    "email": "mali@example.com",
    "homeNumber": "108",
    "soi": "Soi D",
    "road": "Road D",
    "district": "Mueang Phuket",
    "subdistrict": "Wichit",
    "province": "Phuket",
    "zipcode": "83000"
  },
  "parcel": {
    "type": "EXPRESS",
    "weight": 2.5,
    "notes": "Handle with care"
  }
}
```

**PowerShell**

```powershell
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/parcels" `
-ContentType "application/json" `
-Body '{
  "deliver": {
    "name": "Narin",
    "surname": "Customer",
    "phone": "0893333333",
    "email": "narin@example.com",
    "homeNumber": "99/2",
    "soi": "Soi C",
    "road": "Road C",
    "district": "Mueang Chiang Mai",
    "subdistrict": "Suthep",
    "province": "Chiang Mai",
    "zipcode": "50200"
  },
  "receiver": {
    "name": "Mali",
    "surname": "Customer",
    "phone": "0894444444",
    "email": "mali@example.com",
    "homeNumber": "108",
    "soi": "Soi D",
    "road": "Road D",
    "district": "Mueang Phuket",
    "subdistrict": "Wichit",
    "province": "Phuket",
    "zipcode": "83000"
  },
  "parcel": {
    "type": "EXPRESS",
    "weight": 2.5,
    "notes": "Handle with care"
  }
}'
```

---

### 6.4 GET /api/parcels

ใช้สำหรับดูรายการพัสดุทั้งหมด

**Request**

```http
GET /api/parcels
```

**PowerShell**

```powershell
curl.exe http://localhost:8080/api/parcels
```

---

### 6.5 GET /api/parcels/{parcelId}

ใช้สำหรับดูรายละเอียดพัสดุตาม parcel id

**Request**

```http
GET /api/parcels/PCL-000002
```

**PowerShell**

```powershell
curl.exe http://localhost:8080/api/parcels/PCL-000002
```

---

### 6.6 GET /api/trackings/{trackId}

ใช้สำหรับดูสถานะและประวัติการติดตามพัสดุ

**Request**

```http
GET /api/trackings/TRK-000002
```

**PowerShell**

```powershell
curl.exe http://localhost:8080/api/trackings/TRK-000002
```

---

### 6.7 GET /api/shipping/calculate

ใช้สำหรับคำนวณค่าจัดส่ง

**Request**

```http
GET /api/shipping/calculate?origin=Bangkok&destination=Phuket&deliveryType=EXPRESS&weight=2.5
```

**PowerShell**

```powershell
curl.exe "http://localhost:8080/api/shipping/calculate?origin=Bangkok&destination=Phuket&deliveryType=EXPRESS&weight=2.5"
```

หมายเหตุ:
- `deliveryType` ต้องเป็น `STANDARD` หรือ `EXPRESS`

---

### 6.8 PUT /api/trackings/{trackId}/status

ใช้สำหรับอัปเดตสถานะการขนส่งของพัสดุ

**Request**

```http
PUT /api/trackings/TRK-000002/status
Content-Type: application/json
```

```json
{
  "status": "Delivered",
  "description": "Parcel delivered successfully",
  "location": "Phuket Destination Hub"
}
```

**PowerShell**

```powershell
Invoke-RestMethod -Method Put -Uri "http://localhost:8080/api/trackings/TRK-000002/status" -ContentType "application/json" -Body '{"status":"Delivered","description":"Parcel delivered successfully","location":"Phuket Destination Hub"}'
```

สถานะที่รองรับ

- Pending
- In Transit
- Delivered
- Delivery Failed

---

### 6.9 GET /api/messenger/tasks

ใช้สำหรับดูงานของ messenger

**Request**

```http
GET /api/messenger/tasks
Authorization: Bearer <messenger-token>
```

**PowerShell**

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/messenger/tasks" -Headers @{ Authorization = "Bearer <messenger-token>" }
```

หมายเหตุ:
- ถ้าไม่มีงาน อาจได้ `null` หรือรายการว่าง
- แนะนำให้ login ด้วย `EMP002` เพื่อทดสอบเส้นนี้

---

### 6.10 POST /api/vehicle/assign

ใช้สำหรับ assign รถและ messenger ให้พัสดุผ่าน track id

**Request**

```http
POST /api/vehicle/assign
Content-Type: application/json
```

```json
{
  "trackID": "TRK-000002",
  "vehicleID": 1,
  "messengerID": 2
}
```

**PowerShell**

```powershell
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/vehicle/assign" -ContentType "application/json" -Body '{"trackID":"TRK-000002","vehicleID":1,"messengerID":2}'
```

หมายเหตุ:
- ต้องใช้ชื่อ field แบบ camelCase
- ใช้ `trackID` ไม่ใช่ `track_id`

---

### 6.11 POST /api/reports

ใช้สำหรับสร้างรายงานปัญหา

**Request**

```http
POST /api/reports
Content-Type: application/json
```

```json
{
  "trackID": "TRK-000002",
  "issue": "Delayed Delivery",
  "reason": "Receiver was not available at the address",
  "firstName": "Aida",
  "lastName": "Bunsuea",
  "phone": "0812345678",
  "email": "clerk@kencat.local"
}
```

**PowerShell**

```powershell
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/reports" -ContentType "application/json" -Body '{"trackID":"TRK-000002","issue":"Delayed Delivery","reason":"Receiver was not available at the address","firstName":"Aida","lastName":"Bunsuea","phone":"0812345678","email":"clerk@kencat.local"}'
```

---

### 6.12 GET /api/reports

ใช้สำหรับดูรายงานทั้งหมด

**Request**

```http
GET /api/reports
```

**PowerShell**

```powershell
curl.exe http://localhost:8080/api/reports
```


## 7. Notes

- Test

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html

```

- login หน้าเว็บ

EmployeeID: messenger@kencat.local
Password: 1234
Role: Messenger

EmployeeID: clerk@kencat.local
Password: 1234
Role: Parcel Clerk

EmployeeID: admin@kencat.local
Password: 1234

- Route `/api/users/profile` มีอยู่ในโปรเจกต์ แต่ไม่รวมใน 12 เส้นที่ต้องยื่น
- Route `/api/vehicles/select` มีอยู่ในโปรเจกต์ แต่ไม่รวมใน 12 เส้นที่ต้องยื่น
- ถ้าทดสอบใน PowerShell แนะนำใช้ `curl.exe` หรือ `Invoke-RestMethod`
- Route ที่มีการป้องกันสิทธิ์ต้องส่ง `Authorization: Bearer <token>`
