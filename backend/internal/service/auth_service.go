// ด้า
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"kencatexpress/backend/internal/dto"
	"kencatexpress/backend/internal/util"
)

// ด้า
// หน้าที่: POST /api/auth/login
// Service นี้รับผิดชอบตรวจสอบ employeeID/password/role และสร้าง JWT token
// Test ที่เกี่ยวข้อง: auth_service_test.go
type AuthService struct {
	employees    EmployeeStore
	jwtSecret    string
	passwordSalt string
	tokenTTL     time.Duration
}

// ด้า - สร้าง AuthService โดยรับ repository interface เพื่อให้ mock ใน test ได้
func NewAuthService(employees EmployeeStore, jwtSecret, passwordSalt string, tokenTTL time.Duration) *AuthService {
	return &AuthService{employees: employees, jwtSecret: jwtSecret, passwordSalt: passwordSalt, tokenTTL: tokenTTL}
}

// ด้า - Login ตรวจ credentials แล้วคืน token และข้อมูลผู้ใช้
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	identity := util.NormalizeWhitespace(req.EmployeeID)
	if identity == "" || strings.TrimSpace(req.Password) == "" {
		return nil, errors.New("employeeID and password are required")
	}
	role := util.NormalizeRole(req.Role)
	emp, err := s.employees.FindEmployeeByIdentity(ctx, identity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}
	if !emp.IsActive {
		return nil, errors.New("employee account is inactive")
	}
	if role != "" && role != util.NormalizeRole(emp.Role) {
		return nil, errors.New("role does not match the account")
	}
	if !util.CheckPassword(req.Password, emp.PasswordHash, s.passwordSalt) {
		return nil, errors.New("invalid credentials")
	}
	now := time.Now()
	claims := util.Claims{Subject: fmt.Sprintf("%d", emp.ID), EmployeeID: emp.ID, EmployeeCode: emp.EmployeeCode, Email: emp.Email, Role: util.NormalizeRole(emp.Role), Name: strings.TrimSpace(emp.FirstName + " " + emp.LastName), IssuedAt: now.Unix(), ExpiresAt: now.Add(s.tokenTTL).Unix()}
	token, err := util.BuildJWT(claims, s.jwtSecret)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{Success: true, Token: token, User: map[string]interface{}{"id": emp.ID, "employeeCode": emp.EmployeeCode, "email": emp.Email, "firstName": emp.FirstName, "lastName": emp.LastName, "role": util.NormalizeRole(emp.Role), "roleCode": util.RoleToLegacyCode(emp.Role), "displayRole": util.RoleDisplayName(emp.Role)}}, nil
}
