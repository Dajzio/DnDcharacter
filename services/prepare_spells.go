package services

import (
	"context"
	"fmt"
	"strings"

	"starter_pack/domain"
)

type PrepareSpellService struct {
	Repo interface {
		GetByName(ctx context.Context, name string) (*domain.Character, error)
		Save(ctx context.Context, c *domain.Character) error
	}
}

func (s *PrepareSpellService) Execute(ctx context.Context, name string, spellName string) (string, error) {
	char, err := s.Repo.GetByName(ctx, name)
	if err != nil {
		return "", fmt.Errorf("character not found: %w", err)
	}

	class := strings.ToLower(string(char.Class))
	spellName = strings.ToLower(strings.TrimSpace(spellName))

	if !domain.IsSpellcastingClass(class) {
		return "", fmt.Errorf("this class can't cast spells")
	}

	if !domain.PreparesSpells(class) {
		return "", fmt.Errorf("this class learns spells and can't prepare them")
	}

	spell := domain.FindSpellByName(spellName)
	if spell == nil || !domain.ClassHasSpell(char.Class, spellName) {
		return "", fmt.Errorf("spell \"%s\" not found for this class", spellName)
	}

	slots := domain.GetSpellSlots(char.Class, char.Level)
	maxSlotLevel := 0
	for lvl := range slots {
		if lvl > maxSlotLevel {
			maxSlotLevel = lvl
		}
	}

	if spell.Level > maxSlotLevel {
		return "", fmt.Errorf("the spell has higher level than the available spell slots")
	}

	if count, ok := slots[spell.Level]; !ok || count == 0 {
		return "", fmt.Errorf("the spell has higher level than the available spell slots")
	}

	for _, prepared := range char.Spells {
		if strings.EqualFold(prepared.Name, spell.Name) {
			return "", fmt.Errorf("spell already prepared: %s", spell.Name)
		}
	}

	char.Spells = append(char.Spells, *spell)

	if err := s.Repo.Save(ctx, char); err != nil {
		return "", fmt.Errorf("failed to save character: %w", err)
	}

	return fmt.Sprintf("Prepared spell %s", spell.Name), nil
}
