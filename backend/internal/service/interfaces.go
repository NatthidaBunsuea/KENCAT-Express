//หมวย

package service

import (
	"context"

	"kencatexpress/backend/internal/domain"
)

// EmployeeStore คือ repository interface ที่ AuthService และ UserService ใช้
// เพื่อให้สามารถ mock แล้วส่งไปทดสอบใน auth_service_test/user_service_test ได้
type EmployeeStore interface {
	FindEmployeeByIdentity(ctx context.Context, identity string) (*domain.Employee, error)
	FindEmployeeByID(ctx context.Context, id int64) (*domain.Employee, error)
}

// UserStore คือ repository interface สำหรับดึงข้อมูลลูกค้า
type UserStore interface {
	FindCustomerByID(ctx context.Context, id int64) (*domain.User, error)
	FindActiveAddressByUserID(ctx context.Context, userID int64) (*domain.UserAddress, error)
}

// ShippingRateStore คือ repository interface สำหรับคำนวณค่าส่ง
type ShippingRateStore interface {
	FindMatchedRate(ctx context.Context, zoneCode, deliveryType string, weight float64) (*domain.ShippingRate, error)
}

// ParcelStore คือ repository interface หลักของพัสดุ ใช้ร่วมกันหลาย service
type ParcelStore interface {
	CreateParcelGraph(ctx context.Context, input domain.ParcelCreateInput) (*domain.ParcelCreated, error)
	ListParcels(ctx context.Context) ([]domain.ParcelListItem, error)
	GetParcelDetail(ctx context.Context, identifier string) (*domain.ParcelDetail, error)
	GetParcelTrackingView(ctx context.Context, trackID string) (*domain.ParcelTrackingView, error)
	FindParcelByTrackingCode(ctx context.Context, trackID string) (*domain.Parcel, error)
	UpdateParcelStatus(ctx context.Context, input domain.TrackingUpdateInput) error
	ListMessengerTasks(ctx context.Context, employeeID int64) ([]domain.ParcelListItem, error)
	AssignVehicle(ctx context.Context, input domain.VehicleAssignmentInput) error
}

// VehicleStore คือ repository interface สำหรับจัดการรถ
type VehicleStore interface {
	FindVehicleByID(ctx context.Context, id int64) (*domain.Vehicle, error)
	AssignVehicleToEmployee(ctx context.Context, vehicleID, employeeID int64) error
}

// ReportStore คือ repository interface สำหรับรายงานปัญหา
type ReportStore interface {
	CreateReport(ctx context.Context, report *domain.Report) (int64, error)
	ListReports(ctx context.Context) ([]domain.Report, error)
}
