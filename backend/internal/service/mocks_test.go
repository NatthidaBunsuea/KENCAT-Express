// ยีน
package service

// Mock กลางสำหรับ service_test ใช้ร่วมกันทั้ง 6 คน

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"kencatexpress/backend/internal/domain"
	"kencatexpress/backend/internal/dto"
	"kencatexpress/backend/internal/util"
)

var errTestBoom = errors.New("boom")

type testEmployeeStore struct {
	byIdentity      *domain.Employee
	byIdentityErr   error
	byID            *domain.Employee
	byIDErr         error
	lastIdentityArg string
	lastIDArg       int64
}

func (m *testEmployeeStore) FindEmployeeByIdentity(ctx context.Context, identity string) (*domain.Employee, error) {
	m.lastIdentityArg = identity
	if m.byIdentityErr != nil {
		return nil, m.byIdentityErr
	}
	if m.byIdentity == nil {
		return nil, sql.ErrNoRows
	}
	return m.byIdentity, nil
}

func (m *testEmployeeStore) FindEmployeeByID(ctx context.Context, id int64) (*domain.Employee, error) {
	m.lastIDArg = id
	if m.byIDErr != nil {
		return nil, m.byIDErr
	}
	if m.byID == nil {
		return nil, sql.ErrNoRows
	}
	return m.byID, nil
}

type testUserStore struct {
	user        *domain.User
	userErr     error
	address     *domain.UserAddress
	addressErr  error
	lastUserID  int64
	lastAddrUID int64
}

func (m *testUserStore) FindCustomerByID(ctx context.Context, id int64) (*domain.User, error) {
	m.lastUserID = id
	if m.userErr != nil {
		return nil, m.userErr
	}
	if m.user == nil {
		return nil, sql.ErrNoRows
	}
	return m.user, nil
}

func (m *testUserStore) FindActiveAddressByUserID(ctx context.Context, userID int64) (*domain.UserAddress, error) {
	m.lastAddrUID = userID
	if m.addressErr != nil {
		return nil, m.addressErr
	}
	if m.address == nil {
		return nil, sql.ErrNoRows
	}
	return m.address, nil
}

type testRateStore struct {
	rate             *domain.ShippingRate
	err              error
	lastZoneCode     string
	lastDeliveryType string
	lastWeight       float64
}

func (m *testRateStore) FindMatchedRate(ctx context.Context, zoneCode, deliveryType string, weight float64) (*domain.ShippingRate, error) {
	m.lastZoneCode = zoneCode
	m.lastDeliveryType = deliveryType
	m.lastWeight = weight
	if m.err != nil {
		return nil, m.err
	}
	if m.rate == nil {
		return nil, sql.ErrNoRows
	}
	return m.rate, nil
}

type testParcelStore struct {
	createResp   *domain.ParcelCreated
	createErr    error
	listResp     []domain.ParcelListItem
	listErr      error
	detailResp   *domain.ParcelDetail
	detailErr    error
	trackingResp *domain.ParcelTrackingView
	trackingErr  error
	findResp     *domain.Parcel
	findErr      error
	updateErr    error
	assignErr    error

	createdInput   *domain.ParcelCreateInput
	lastDetailID   string
	lastTrackingID string
	lastFindTrack  string
	updateInput    *domain.TrackingUpdateInput
	lastTasksEmpID int64
	assignInput    *domain.VehicleAssignmentInput
}

func (m *testParcelStore) CreateParcelGraph(ctx context.Context, input domain.ParcelCreateInput) (*domain.ParcelCreated, error) {
	cloned := input
	m.createdInput = &cloned
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.createResp != nil {
		return m.createResp, nil
	}
	return &domain.ParcelCreated{ParcelID: 1, ParcelCode: input.ParcelCode, TrackingCode: input.TrackingCode, ShippingCost: input.ShippingCost}, nil
}

func (m *testParcelStore) ListParcels(ctx context.Context) ([]domain.ParcelListItem, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.listResp, nil
}

func (m *testParcelStore) GetParcelDetail(ctx context.Context, identifier string) (*domain.ParcelDetail, error) {
	m.lastDetailID = identifier
	if m.detailErr != nil {
		return nil, m.detailErr
	}
	if m.detailResp == nil {
		return nil, sql.ErrNoRows
	}
	return m.detailResp, nil
}

func (m *testParcelStore) GetParcelTrackingView(ctx context.Context, trackID string) (*domain.ParcelTrackingView, error) {
	m.lastTrackingID = trackID
	if m.trackingErr != nil {
		return nil, m.trackingErr
	}
	if m.trackingResp == nil {
		return nil, sql.ErrNoRows
	}
	return m.trackingResp, nil
}

func (m *testParcelStore) FindParcelByTrackingCode(ctx context.Context, trackID string) (*domain.Parcel, error) {
	m.lastFindTrack = trackID
	if m.findErr != nil {
		return nil, m.findErr
	}
	if m.findResp == nil {
		return nil, sql.ErrNoRows
	}
	return m.findResp, nil
}

func (m *testParcelStore) UpdateParcelStatus(ctx context.Context, input domain.TrackingUpdateInput) error {
	cloned := input
	m.updateInput = &cloned
	if m.updateErr != nil {
		return m.updateErr
	}
	return nil
}

func (m *testParcelStore) ListMessengerTasks(ctx context.Context, employeeID int64) ([]domain.ParcelListItem, error) {
	m.lastTasksEmpID = employeeID
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.listResp, nil
}

func (m *testParcelStore) AssignVehicle(ctx context.Context, input domain.VehicleAssignmentInput) error {
	cloned := input
	m.assignInput = &cloned
	if m.assignErr != nil {
		return m.assignErr
	}
	return nil
}

type testVehicleStore struct {
	vehicle         *domain.Vehicle
	findErr         error
	assignErr       error
	lastFindID      int64
	lastAssignVehID int64
	lastAssignEmpID int64
}

func (m *testVehicleStore) FindVehicleByID(ctx context.Context, id int64) (*domain.Vehicle, error) {
	m.lastFindID = id
	if m.findErr != nil {
		return nil, m.findErr
	}
	if m.vehicle == nil {
		return nil, sql.ErrNoRows
	}
	return m.vehicle, nil
}

func (m *testVehicleStore) AssignVehicleToEmployee(ctx context.Context, vehicleID, employeeID int64) error {
	m.lastAssignVehID = vehicleID
	m.lastAssignEmpID = employeeID
	return m.assignErr
}

type testReportStore struct {
	createID  int64
	createErr error
	listResp  []domain.Report
	listErr   error
	created   *domain.Report
}

func (m *testReportStore) CreateReport(ctx context.Context, report *domain.Report) (int64, error) {
	copied := *report
	m.created = &copied
	if m.createErr != nil {
		return 0, m.createErr
	}
	if m.createID == 0 {
		m.createID = 1
	}
	return m.createID, nil
}

func (m *testReportStore) ListReports(ctx context.Context) ([]domain.Report, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.listResp, nil
}

func makeActiveEmployee(role string) *domain.Employee {
	return &domain.Employee{
		ID:           1,
		EmployeeCode: "EMP001",
		Email:        "employee@kencat.local",
		PasswordHash: util.HashPassword("secret123", "kencat-express-salt"),
		FirstName:    "Aida",
		LastName:     "Bunsuea",
		Role:         role,
		IsActive:     true,
	}
}

func makeValidParcelRequest() dto.ParcelRequest {
	return dto.ParcelRequest{
		Deliver: dto.PersonAddressRequest{
			Name:        " Somchai ",
			Surname:     " Sender ",
			Phone:       "0891111111",
			Email:       " sender@example.com ",
			HomeNumber:  "12/1",
			Soi:         " Soi 1 ",
			Road:        " Rama 9 ",
			District:    "Huai Khwang",
			Subdistrict: "Bang Kapi",
			Province:    "Bangkok",
			Zipcode:     "10310",
		},
		Receiver: dto.PersonAddressRequest{
			Name:        " Mali ",
			Surname:     " Customer ",
			Phone:       "0894444444",
			Email:       " receiver@example.com ",
			HomeNumber:  "108",
			Soi:         "",
			Road:        "",
			District:    "Mueang Phuket",
			Subdistrict: "Wichit",
			Province:    "Phuket",
			Zipcode:     "83000",
		},
		Parcel: dto.ParcelInfoRequest{
			Type:   util.DeliveryExpress,
			Weight: 3.2,
			Notes:  "  Express please  ",
		},
	}
}

func makeParcelListItem() domain.ParcelListItem {
	return domain.ParcelListItem{
		ParcelID:        "PCL-000001",
		TrackID:         "TRK-000001",
		ReceiverName:    "Mali",
		ReceiverSurname: "Customer",
		ReceiverTel:     "0894444444",
		Status:          util.StatusPending,
		DepositDate:     time.Now(),
	}
}

func ptrInt64(v int64) *int64 { return &v }
