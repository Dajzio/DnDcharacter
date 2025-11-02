package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"starter_pack/domain"
	"strings"
	"sync"
)

const ErrCharacterNotFoundMsg = "character not found"

var ErrCharacterNotFound = errors.New(ErrCharacterNotFoundMsg)

type FileCharacterRepo struct {
	mu       sync.Mutex
	filename string
}

func NewFileCharacterRepo(filename string) *FileCharacterRepo {
	return &FileCharacterRepo{
		filename: filename,
	}
}

func (r *FileCharacterRepo) Save(ctx context.Context, c *domain.Character) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var characters []domain.Character
	if data, err := os.ReadFile(r.filename); err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &characters); err != nil {
			return err
		}
	}

	found := false
	for i := range characters {
		if characters[i].ID == c.ID {
			characters[i] = *c
			found = true
			break
		}
	}

	if !found {
		characters = append(characters, *c)
	}

	updated, err := json.MarshalIndent(characters, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filename, updated, 0644)
}

func (r *FileCharacterRepo) List(ctx context.Context) ([]*domain.Character, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var characters []domain.Character
	data, err := os.ReadFile(r.filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []*domain.Character{}, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, &characters); err != nil {
		return nil, err
	}

	result := make([]*domain.Character, 0, len(characters))
	for i := range characters {
		result = append(result, &characters[i])
	}
	return result, nil
}

func (r *FileCharacterRepo) GetByName(ctx context.Context, name string) (*domain.Character, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var characters []domain.Character
	data, err := os.ReadFile(r.filename)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &characters); err != nil {
		return nil, err
	}

	for i := range characters {
		if strings.EqualFold(characters[i].Name, name) {
			return &characters[i], nil
		}
	}
	return nil, ErrCharacterNotFound
}

func (r *FileCharacterRepo) Delete(ctx context.Context, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var characters []domain.Character
	data, err := os.ReadFile(r.filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &characters); err != nil {
		return err
	}

	newList := make([]domain.Character, 0, len(characters))
	found := false
	for _, c := range characters {
		if strings.EqualFold(c.Name, name) {
			found = true
			continue
		}
		newList = append(newList, c)
	}

	if !found {
		return ErrCharacterNotFound
	}

	updated, err := json.MarshalIndent(newList, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filename, updated, 0644)
}

func (r *FileCharacterRepo) GetByID(ctx context.Context, id string) (*domain.Character, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var characters []domain.Character
	data, err := os.ReadFile(r.filename)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &characters); err != nil {
		return nil, err
	}

	for i := range characters {
		if characters[i].ID == id {
			return &characters[i], nil
		}
	}
	return nil, ErrCharacterNotFound
}
