package services

import (
	"bytes"
	"context"
	"os"
	"starter_pack/domain"
	"strings"
	"testing"
)

func TestViewCharacterServiceSuccess(t *testing.T) {
	char := &domain.Character{
		Name:      "Gandalf",
		Class:     "Wizard",
		Race:      "Maia",
		Level:     10,
		Background: "Sage",
		ProficiencyBonus: 4,
		AbilityScores: domain.AbilityScores{Str: 10, Dex: 14, Con: 12, Int: 18, Wis: 16, Cha: 12},
		SkillProficiencies: []string{"Arcana", "History"},
		Equipment: domain.Equipment{
			MainHandWeapon: &domain.Weapon{Name: "Staff"},
			Armor:          &domain.Armor{Name: "Robe"},
		},
		ArmorClass: 12,
		Initiative: 2,
		PassivePerception: 13,
		SpellSlots: map[int]int{1: 4, 2: 3, 3: 3},
		SpellcastingAbility: "INT",
		SpellSaveDC: 15,
		SpellAttackBonus: 7,
	}

	repo := &MockCharacterRepo{
		Characters: map[string]*domain.Character{"Gandalf": char},
	}

	service := &ViewCharacterService{Repo: repo}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := service.Execute(context.Background(), "Gandalf")
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "Gandalf") {
		t.Errorf("expected character name in output")
	}
	if !strings.Contains(output, "Wizard") {
		t.Errorf("expected class in output")
	}
	if !strings.Contains(output, "Staff") {
		t.Errorf("expected weapon in output")
	}
	if !strings.Contains(output, "Spellcasting ability: intelligence") {
		t.Errorf("expected spellcasting ability")
	}
}

func TestViewCharacterServiceNotFound(t *testing.T) {
	repo := &MockCharacterRepo{Characters: map[string]*domain.Character{}}
	service := &ViewCharacterService{Repo: repo}

	err := service.Execute(context.Background(), "Unknown")
	if err == nil || !strings.Contains(err.Error(), "character") {
		t.Errorf("expected character not found error, got %v", err)
	}
}

func TestFullAbilityName(t *testing.T) {
	tests := map[string]string{
		"STR": "strength",
		"Dex": "dexterity",
		"con": "constitution",
		"Int": "intelligence",
		"WIS": "wisdom",
		"cha": "charisma",
		"UNK": "unk",
	}

	for input, expected := range tests {
		got := FullAbilityName(input)
		if got != expected {
			t.Errorf("FullAbilityName(%q) = %q, want %q", input, got, expected)
		}
	}
}
