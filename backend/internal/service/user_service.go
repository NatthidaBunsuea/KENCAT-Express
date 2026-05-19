// ด้า
package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"kencatexpress/backend/internal/domain"
	"kencatexpress/backend/internal/dto"
	"kencatexpress/backend/internal/util"
)

// ด้า
// หน้าที่: GET /api/users/{userId}
// Service นี้รับผิดชอบดึงข้อมูลลูกค้า และมีฟังก์ชันดึง profile พนักงานสำหรับหน้า profile
// Test ที่เกี่ยวข้อง: user_service_test.go
type UserService struct {
	employees EmployeeStore
	users     UserStore
}

// ด้า - สร้าง UserService โดยรับ repository interface เพื่อให้ mock ใน test ได้
func NewUserService(employees EmployeeStore, users UserStore) *UserService {
	return &UserService{employees: employees, users: users}
}

// ด้า - ดึงข้อมูลลูกค้าตาม id
func (s *UserService) GetCustomerByID(ctx context.Context, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user id")
	}
	user, err := s.users.FindCustomerByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

// ด้า - ดึง profile พนักงานที่ login อยู่
func (s *UserService) GetEmployeeProfile(ctx context.Context, employeeID int64) (*dto.EmployeeProfileResponse, error) {
	if employeeID <= 0 {
		return nil, errors.New("invalid employee id")
	}
	emp, err := s.employees.FindEmployeeByID(ctx, employeeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("employee not found")
		}
		return nil, err
	}
	birthDate := ""
	if emp.BirthDate != nil {
		birthDate = emp.BirthDate.Format(time.RFC3339)
	}
	return &dto.EmployeeProfileResponse{EmployeeID: emp.EmployeeCode, Role: util.RoleToLegacyCode(emp.Role), Firstname: emp.FirstName, Lastname: emp.LastName, Email: emp.Email, Birthdate: birthDate, EmployeeCode: emp.EmployeeCode, DisplayRole: util.RoleDisplayName(emp.Role), FirstName: emp.FirstName, LastName: emp.LastName, RoleName: util.NormalizeRole(emp.Role)}, nil
}
