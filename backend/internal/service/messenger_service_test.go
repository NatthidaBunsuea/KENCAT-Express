// โอม
package service

import (
	"context"
	"testing"

	"kencatexpress/backend/internal/domain"
)

// โอม
func TestMessengerService_ListTasks(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		parcelStore := &testParcelStore{
			listResp: []domain.ParcelListItem{makeParcelListItem()},
		}
		svc := NewMessengerService(parcelStore)

		got, err := svc.ListTasks(context.Background(), 99)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("expected 1 task, got %d", len(got))
		}
		if parcelStore.lastTasksEmpID != 99 {
			t.Fatalf("expected employee id 99, got %d", parcelStore.lastTasksEmpID)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		parcelStore := &testParcelStore{
			listErr: errTestBoom,
		}
		svc := NewMessengerService(parcelStore)

		got, err := svc.ListTasks(context.Background(), 99)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != errTestBoom.Error() {
			t.Fatalf("expected error %q, got %q", errTestBoom.Error(), err.Error())
		}
		if got != nil {
			t.Fatalf("expected nil tasks, got %v", got)
		}
	})
}
