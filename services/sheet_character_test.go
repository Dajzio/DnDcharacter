package services

import (
	"context"
	"starter_pack/domain"
	"strings"
	"testing"
)

func TestCharacterSheetService_Success(t *testing.T) {
	char := &domain.Character{
		Name:      "Aragorn",
		Class:     "Ranger",
		Race:      "Human",
		Level:     5,
		Background: "Noble",
		ProficiencyBonus: 3,
		AbilityScores: domain.AbilityScores{Str: 16, Dex: 14, Con: 14, Int: 12, Wis: 13, Cha: 12},
		SkillProficiencies: []string{"Athletics", "Survival"},
		Equipment: domain.Equipment{
			MainHandWeapon: &domain.Weapon{Name: "Longsword"},
			Armor:          &domain.Armor{Name: "Chain Shirt"},
			Shield:         &domain.Shield{Name: "Shield"},
		},
		ArmorClass: 16,
		Initiative: 2,
		PassivePerception: 11,
		SpellSlots: map[int]int{1: 4, 2: 2},
		SpellcastingAbility: "WIS",
		SpellSaveDC: 13,
		SpellAttackBonus: 5,
		Spells: []domain.Spell{
			{Name: "Hunter's Mark", Level: 1, Class: []string{"ranger"}},
			{Name: "Cure Wounds", Level: 1, Class: []string{"ranger"}},
		},
	}

	repo := &MockCharacterRepo{
		Characters: map[string]*domain.Character{"Aragorn": char},
	}

	service := &CharacterSheetService{Repo: repo}
	output, err := service.Execute(context.Background(), "Aragorn", "markdown")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "# Aragorn") {
		t.Errorf("expected name in sheet")
	}
	if !strings.Contains(output, "## Ability scores") {
		t.Errorf("expected ability scores section")
	}
	if !strings.Contains(output, "## Equipment") {
		t.Errorf("expected equipment section")
	}
	if !strings.Contains(output, "### Level 1") {
		t.Errorf("expected spell section for level 1")
	}
	if !strings.Contains(output, "Hunter's Mark") {
		t.Errorf("expected spell Hunter's Mark")
	}
}

func TestCharacterSheetService_CharacterNotFound(t *testing.T) {
	repo := &MockCharacterRepo{
		Characters: map[string]*domain.Character{},
	}
	service := &CharacterSheetService{Repo: repo}
	_, err := service.Execute(context.Background(), "Unknown", "markdown")
	if err == nil || !strings.Contains(err.Error(), "character not found") {
		t.Errorf("expected character not found error, got %v", err)
	}
}

func TestCharacterSheetService_UnsupportedFormat(t *testing.T) {
	char := &domain.Character{Name: "Frodo"}
	repo := &MockCharacterRepo{Characters: map[string]*domain.Character{"Frodo": char}}
	service := &CharacterSheetService{Repo: repo}
	_, err := service.Execute(context.Background(), "Frodo", "pdf")
	if err == nil || !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("expected unsupported format error, got %v", err)
	}
}
