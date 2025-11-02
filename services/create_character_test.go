package services

import (
	"context"
	"errors"
	"starter_pack/domain"
	"testing"
)

func TestCreateCharacterService_Success(t *testing.T) {
	repo := &MockCharacterRepo{Characters: make(map[string]*domain.Character)}
	factory := &domain.CharacterFactory{}
	service := &CreateCharacterService{Repo: repo, Factory: factory}

	input := CreateCharacterInput{
		Name:       "TestChar",
		Race:       "Human",
		Class:      "Wizard",
		Level:      1,
		Str:        10,
		Dex:        12,
		Con:        14,
		Int:        16,
		Wis:        8,
		Cha:        10,
		Background: "Sage",
		Skills:     []string{"Arcana", "History"},
	}

	result, err := service.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatalf("expected character, got nil")
	}

	if result.Name != input.Name {
		t.Errorf("expected name %s, got %s", input.Name, result.Name)
	}

	if _, ok := repo.Characters[result.Name]; !ok {
		t.Errorf("expected character to be saved in repo")
	}
}

func TestCreateCharacterService_RepoError(t *testing.T) {
	repo := &MockCharacterRepo{
		Characters: make(map[string]*domain.Character),
		SaveErr:    errors.New("save error"),
	}
	factory := &domain.CharacterFactory{}
	service := &CreateCharacterService{Repo: repo, Factory: factory}

	input := CreateCharacterInput{
		Name:       "TestChar",
		Race:       "Elf",
		Class:      "Rogue",
		Level:      1,
		Str:        8,
		Dex:        16,
		Con:        12,
		Int:        10,
		Wis:        10,
		Cha:        10,
		Background: "Criminal",
		Skills:     []string{"Stealth", "Acrobatics"},
	}

	_, err := service.Execute(context.Background(), input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestCreateCharacterService_InvalidInput(t *testing.T) {
	repo := &MockCharacterRepo{Characters: make(map[string]*domain.Character)}
	factory := &domain.CharacterFactory{}
	service := &CreateCharacterService{Repo: repo, Factory: factory}

	input := CreateCharacterInput{
		Name:       "",
		Race:       "Human",
		Class:      "Wizard",
		Level:      0,
		Str:        10,
		Dex:        10,
		Con:        10,
		Int:        10,
		Wis:        10,
		Cha:        10,
		Background: "Sage",
		Skills:     []string{},
	}

	_, err := service.Execute(context.Background(), input)
	if err == nil {
		t.Fatalf("expected error for invalid input")
	}
}
