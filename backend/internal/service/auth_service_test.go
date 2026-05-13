// ด้า
package service

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"kencatexpress/backend/internal/dto"
	"kencatexpress/backend/internal/util"
)

// ด้า
func TestAuthService_Login_Success(t *testing.T) {
	store := &testEmployeeStore{byIdentity: makeActiveEmployee(util.RoleParcelClerk)}
	svc := NewAuthService(store, "secret", "kencat-express-salt", time.Hour)

	resp, err := svc.Login(context.Background(), dto.LoginRequest{
		EmployeeID: " employee@kencat.local ",
		Password:   "secret123",
		Role:       "Parcel Clerk",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || !resp.Success {
		t.Fatal("expected success response")
	}
	if resp.Token == "" {
		t.Fatal("expected token to be generated")
	}
	if store.lastIdentityArg != "employee@kencat.local" {
		t.Fatalf("expected normalized identity, got %q", store.lastIdentityArg)
	}
	if got := resp.User["role"]; got != util.RoleParcelClerk {
		t.Fatalf("expected role %q, got %v", util.RoleParcelClerk, got)
	}
}

// ด้า
func TestAuthService_Login_ValidationErrors(t *testing.T) {
	svc := NewAuthService(&testEmployeeStore{}, "secret", "kencat-express-salt", time.Hour)

	cases := []struct {
		name string
		req  dto.LoginRequest
		want string
	}{
		{name: "missing employee id", req: dto.LoginRequest{Password: "secret123"}, want: "employeeID and password are required"},
		{name: "missing password", req: dto.LoginRequest{EmployeeID: "EMP001"}, want: "employeeID and password are required"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.Login(context.Background(), tc.req)
			if err == nil || err.Error() != tc.want {
				t.Fatalf("expected error %q, got %v", tc.want, err)
			}
		})
	}
}

// ด้า
func TestAuthService_Login_FailureCases(t *testing.T) {
	active := makeActiveEmployee(util.RoleParcelClerk)
	inactive := makeActiveEmployee(util.RoleParcelClerk)
	inactive.IsActive = false

	cases := []struct {
		name  string
		store *testEmployeeStore
		req   dto.LoginRequest
		want  string
	}{
		{
			name:  "employee not found",
			store: &testEmployeeStore{byIdentityErr: sql.ErrNoRows},
			req:   dto.LoginRequest{EmployeeID: "EMP404", Password: "secret123"},
			want:  "invalid credentials",
		},
		{
			name:  "employee inactive",
			store: &testEmployeeStore{byIdentity: inactive},
			req:   dto.LoginRequest{EmployeeID: "EMP001", Password: "secret123"},
			want:  "employee account is inactive",
		},
		{
			name:  "role mismatch",
			store: &testEmployeeStore{byIdentity: active},
			req:   dto.LoginRequest{EmployeeID: "EMP001", Password: "secret123", Role: "Messenger"},
			want:  "role does not match the account",
		},
		{
			name:  "wrong password",
			store: &testEmployeeStore{byIdentity: active},
			req:   dto.LoginRequest{EmployeeID: "EMP001", Password: "wrongpass", Role: "Parcel Clerk"},
			want:  "invalid credentials",
		},
		{
			name:  "repository error",
			store: &testEmployeeStore{byIdentityErr: errTestBoom},
			req:   dto.LoginRequest{EmployeeID: "EMP001", Password: "secret123"},
			want:  errTestBoom.Error(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewAuthService(tc.store, "secret", "kencat-express-salt", time.Hour)
			_, err := svc.Login(context.Background(), tc.req)
			if err == nil {
				t.Fatal("expected an error")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected error containing %q, got %v", tc.want, err)
			}
		})
	}
}
