package services

import (
	"context"
	"errors"
	"starter_pack/domain"
	"testing"
)

const SpellMagicMissile = "Magic Missile"

func TestPrepareSpellServiceSuccess(t *testing.T) {
	char := &domain.Character{
		Name:  "Merlin",
		Class: "Wizard",
		Level: 3,
		Spells: []domain.Spell{},
	}
	repo := &MockCharacterRepo{
		Characters: map[string]*domain.Character{"Merlin": char},
	}
	spellRepo := &MockSpellRepo{
		Spells: map[string]domain.Spell{
			SpellMagicMissile: {Name: SpellMagicMissile, Level: 1, Class: []string{"wizard"}},
		},
	}

	service := &PrepareSpellService{Repo: repo, SpellRepo: spellRepo}
	msg, err := service.Execute(context.Background(), "Merlin", SpellMagicMissile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != "Prepared spell "+SpellMagicMissile {
		t.Errorf("unexpected message: %s", msg)
	}
	if len(char.Spells) != 1 || char.Spells[0].Name != SpellMagicMissile {
		t.Errorf("spell not added to character")
	}
}

func TestPrepareSpellServiceCharacterNotFound(t *testing.T) {
	repo := &MockCharacterRepo{Characters: map[string]*domain.Character{}}
	spellRepo := &MockSpellRepo{}
	service := &PrepareSpellService{Repo: repo, SpellRepo: spellRepo}

	_, err := service.Execute(context.Background(), "Unknown", SpellMagicMissile)
	if err == nil || err.Error() != "character not found: character not found" {
		t.Errorf("expected character not found error, got %v", err)
	}
}

func TestPrepareSpellServiceSpellNotFound(t *testing.T) {
	char := &domain.Character{Name: "Merlin", Class: "Wizard", Level: 3}
	repo := &MockCharacterRepo{
		Characters: map[string]*domain.Character{"Merlin": char},
	}
	spellRepo := &MockSpellRepo{Spells: map[string]domain.Spell{}}

	service := &PrepareSpellService{Repo: repo, SpellRepo: spellRepo}
	_, err := service.Execute(context.Background(), "Merlin", "Unknown Spell")
	if err == nil || err.Error() != "spell not found: Unknown Spell" {
		t.Errorf("expected spell not found error, got %v", err)
	}
}

func TestPrepareSpellServiceSaveError(t *testing.T) {
	char := &domain.Character{Name: "Merlin", Class: "Wizard", Level: 3}
	repo := &MockCharacterRepo{
		Characters: map[string]*domain.Character{"Merlin": char},
		SaveErr:    errors.New("save failed"),
	}
	spellRepo := &MockSpellRepo{
		Spells: map[string]domain.Spell{SpellMagicMissile: {Name: SpellMagicMissile, Level: 1, Class: []string{"wizard"}}},
	}

	service := &PrepareSpellService{Repo: repo, SpellRepo: spellRepo}
	_, err := service.Execute(context.Background(), "Merlin", SpellMagicMissile)
	if err == nil || err.Error() != "failed to save character: save failed" {
		t.Errorf("expected save error, got %v", err)
	}
}
