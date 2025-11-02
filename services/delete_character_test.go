package services

import (
	"context"
	"errors"
	"testing"

	"starter_pack/domain"
)

type mockDeleteRepo struct {
	deleteFunc func(ctx context.Context, name string) error
}

func (m *mockDeleteRepo) Save(ctx context.Context, c *domain.Character) error           { return nil }
func (m *mockDeleteRepo) GetByID(ctx context.Context, id string) (*domain.Character, error) { return nil, nil }
func (m *mockDeleteRepo) GetByName(ctx context.Context, name string) (*domain.Character, error) {
	return nil, nil
}
func (m *mockDeleteRepo) List(ctx context.Context) ([]*domain.Character, error) { return nil, nil }
func (m *mockDeleteRepo) Delete(ctx context.Context, name string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, name)
	}
	return nil
}

func TestDeleteCharacterService_Success(t *testing.T) {
	mockRepo := &mockDeleteRepo{
		deleteFunc: func(ctx context.Context, name string) error {
			if name != "Gandalf" {
				t.Errorf("expected name 'Gandalf', got '%s'", name)
			}
			return nil
		},
	}

	service := &DeleteCharacterService{Repo: mockRepo}

	err := service.Execute(context.Background(), "Gandalf")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestDeleteCharacterService_Error(t *testing.T) {
	expectedErr := errors.New("character not found")
	mockRepo := &mockDeleteRepo{
		deleteFunc: func(ctx context.Context, name string) error {
			return expectedErr
		},
	}

	service := &DeleteCharacterService{Repo: mockRepo}

	err := service.Execute(context.Background(), "Frodo")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
