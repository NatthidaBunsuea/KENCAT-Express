// ด้า
package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"kencatexpress/backend/internal/domain"
	"kencatexpress/backend/internal/util"
)

// ด้า
func TestUserService_GetCustomerByID_Success(t *testing.T) {
	user := &domain.User{ID: 7, FirstName: "Mali", LastName: "Customer"}
	svc := NewUserService(&testEmployeeStore{}, &testUserStore{user: user})

	got, err := svc.GetCustomerByID(context.Background(), 7)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got == nil || got.ID != 7 {
		t.Fatalf("expected user id 7, got %+v", got)
	}
}

// ด้า
func TestUserService_GetCustomerByID_Failures(t *testing.T) {
	svc := NewUserService(&testEmployeeStore{}, &testUserStore{})

	if _, err := svc.GetCustomerByID(context.Background(), 0); err == nil || err.Error() != "invalid user id" {
		t.Fatalf("expected invalid user id error, got %v", err)
	}

	svc = NewUserService(&testEmployeeStore{}, &testUserStore{userErr: sql.ErrNoRows})
	if _, err := svc.GetCustomerByID(context.Background(), 99); err == nil || err.Error() != "user not found" {
		t.Fatalf("expected user not found, got %v", err)
	}

	svc = NewUserService(&testEmployeeStore{}, &testUserStore{userErr: errTestBoom})
	if _, err := svc.GetCustomerByID(context.Background(), 99); err == nil || err.Error() != errTestBoom.Error() {
		t.Fatalf("expected repository error, got %v", err)
	}
}

// ด้า
func TestUserService_GetEmployeeProfile_Success(t *testing.T) {
	dob := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
	emp := makeActiveEmployee(util.RoleMessenger)
	emp.BirthDate = &dob
	emp.EmployeeCode = "EMP002"
	svc := NewUserService(&testEmployeeStore{byID: emp}, &testUserStore{})

	got, err := svc.GetEmployeeProfile(context.Background(), 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.EmployeeCode != "EMP002" {
		t.Fatalf("expected employee code EMP002, got %s", got.EmployeeCode)
	}
	if got.DisplayRole != "Messenger" {
		t.Fatalf("expected display role Messenger, got %s", got.DisplayRole)
	}
	if got.Birthdate == "" {
		t.Fatal("expected birthdate to be populated")
	}
}

// ด้า
func TestUserService_GetEmployeeProfile_Failures(t *testing.T) {
	svc := NewUserService(&testEmployeeStore{}, &testUserStore{})

	if _, err := svc.GetEmployeeProfile(context.Background(), 0); err == nil || err.Error() != "invalid employee id" {
		t.Fatalf("expected invalid employee id, got %v", err)
	}

	svc = NewUserService(&testEmployeeStore{byIDErr: sql.ErrNoRows}, &testUserStore{})
	if _, err := svc.GetEmployeeProfile(context.Background(), 10); err == nil || err.Error() != "employee not found" {
		t.Fatalf("expected employee not found, got %v", err)
	}

	svc = NewUserService(&testEmployeeStore{byIDErr: errTestBoom}, &testUserStore{})
	if _, err := svc.GetEmployeeProfile(context.Background(), 10); err == nil || err.Error() != errTestBoom.Error() {
		t.Fatalf("expected repository error, got %v", err)
	}
}
