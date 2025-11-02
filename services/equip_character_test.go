package services

import (
	"context"
	"errors"
	"testing"
	"starter_pack/domain"
)

type mockRepo struct {
	characters map[string]*domain.Character
}

func newMockRepo() *mockRepo {
	return &mockRepo{characters: make(map[string]*domain.Character)}
}

func (r *mockRepo) Save(ctx context.Context, c *domain.Character) error {
	r.characters[c.Name] = c
	return nil
}

func (r *mockRepo) GetByName(ctx context.Context, name string) (*domain.Character, error) {
	c, ok := r.characters[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return c, nil
}

func (r *mockRepo) Delete(ctx context.Context, name string) error {
	delete(r.characters, name)
	return nil
}

func (r *mockRepo) GetByID(ctx context.Context, id string) (*domain.Character, error) {
	for _, c := range r.characters {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, errors.New("not found")
}

func (r *mockRepo) List(ctx context.Context) ([]*domain.Character, error) {
	list := []*domain.Character{}
	for _, c := range r.characters {
		list = append(list, c)
	}
	return list, nil
}

func TestEquipItemService(t *testing.T) {
	repo := newMockRepo()
	char := &domain.Character{Name: "Hero"}
	repo.Save(context.Background(), char)

	service := &EquipItemService{Repo: repo}

	_, err := service.Execute(context.Background(), "Hero", "weapon", "Longsword", "main hand")
	if err != nil {
		t.Fatalf("failed to equip weapon: %v", err)
	}

	_, err = service.Execute(context.Background(), "Hero", "armor", "Leather", "")
	if err != nil {
		t.Fatalf("failed to equip armor: %v", err)
	}

	_, err = service.Execute(context.Background(), "Hero", "shield", "Wooden Shield", "")
	if err != nil {
		t.Fatalf("failed to equip shield: %v", err)
	}

	_, err = service.Execute(context.Background(), "Hero", "weapon", "Dagger", "")
	if err == nil {
		t.Fatalf("expected error when weapon slot is missing")
	}

	_, err = service.Execute(context.Background(), "Hero", "unknown", "Item", "")
	if err == nil {
		t.Fatalf("expected error when item type is unknown")
	}
}
