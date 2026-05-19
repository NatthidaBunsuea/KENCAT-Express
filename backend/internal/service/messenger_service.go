// โอม
package service

import (
	"context"

	"kencatexpress/backend/internal/domain"
)

// โอม
// หน้าที่: GET /api/messenger/tasks
// Service นี้รับผิดชอบดึงรายการงานของ messenger
// Test ที่เกี่ยวข้อง: messenger_service_test.go
type MessengerService struct{ parcels ParcelStore }

// โอม - สร้าง MessengerService โดยรับ ParcelStore เพื่อ mock ใน test ได้
func NewMessengerService(parcels ParcelStore) *MessengerService {
	return &MessengerService{parcels: parcels}
}

// โอม - ดึงงานของ messenger ตาม employee id
func (s *MessengerService) ListTasks(ctx context.Context, employeeID int64) ([]domain.ParcelListItem, error) {
	return s.parcels.ListMessengerTasks(ctx, employeeID)
}
