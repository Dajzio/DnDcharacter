package services

import (
	"context"
	"starter_pack/domain"
	"testing"
)

func TestLearnSpellServiceSuccess(t *testing.T) {
	char := &domain.Character{
		Name:   "Gandalf",
		Class:  "Wizard",
		Level:  3,
		Spells: []domain.Spell{},
	}
	repo := &MockCharacterRepo{
		Characters: map[string]*domain.Character{"Gandalf": char},
	}
	spellRepo := &MockSpellRepo{
		Spells: map[string]domain.Spell{"Fireball": {Name: "Fireball", Level: 3, Class: []string{"wizard"}}},
	}

	service := &LearnSpellService{Repo: repo, SpellRepo: spellRepo}
	msg, err := service.Execute(context.Background(), "Gandalf", "Fireball")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != "Learned spell Fireball" {
		t.Errorf("unexpected message: %s", msg)
	}
	if len(char.Spells) != 1 || char.Spells[0].Name != "Fireball" {
		t.Errorf("spell not added to character")
	}
}

func TestLearnSpellServiceCharacterNotFound(t *testing.T) {
	repo := &MockCharacterRepo{Characters: map[string]*domain.Character{}}
	spellRepo := &MockSpellRepo{Spells: map[string]domain.Spell{"Fireball": {Name: "Fireball"}}}

	service := &LearnSpellService{Repo: repo, SpellRepo: spellRepo}
	_, err := service.Execute(context.Background(), "Unknown", "Fireball")
	if err == nil || err.Error() != "character not found: character not found" {
		t.Errorf("expected character not found error, got %v", err)
	}
}

func TestLearnSpellServiceSpellNotFound(t *testing.T) {
	char := &domain.Character{Name: "Gandalf"}
	repo := &MockCharacterRepo{Characters: map[string]*domain.Character{"Gandalf": char}}
	spellRepo := &MockSpellRepo{Spells: map[string]domain.Spell{}}

	service := &LearnSpellService{Repo: repo, SpellRepo: spellRepo}
	_, err := service.Execute(context.Background(), "Gandalf", "Unknown Spell")
	if err == nil || err.Error() != "spell not found: Unknown Spell" {
		t.Errorf("expected spell not found error, got %v", err)
	}
}
