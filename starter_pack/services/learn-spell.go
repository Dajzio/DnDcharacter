package services

import (
	"context"
	"fmt"
	"strings"

	"starter_pack/domain"
)

type LearnSpellService struct {
	Repo interface {
		GetByName(ctx context.Context, name string) (*domain.Character, error)
		Save(ctx context.Context, c *domain.Character) error
	}
}

func (s *LearnSpellService) Execute(ctx context.Context, name string, spellName string) (string, error) {
	char, err := s.Repo.GetByName(ctx, name)
	if err != nil {
		return "", fmt.Errorf("character not found: %w", err)
	}

	class := strings.ToLower(string(char.Class))
	spellName = strings.ToLower(spellName)

	if !domain.IsSpellcastingClass(class) {
		return "", fmt.Errorf("this class can't cast spells")
	}

	if domain.PreparesSpells(class) {
		return "", fmt.Errorf("this class prepares spells and can't learn them")
	}

	if !domain.ClassHasSpell(char.Class, spellName) {
		return "", fmt.Errorf("spell \"%s\" not found for this class", spellName)
	}

	spell := domain.FindSpellByName(spellName)
	if spell == nil {
		return "", fmt.Errorf("spell not found: %s", spellName)
	}

	for _, known := range char.Spells {
		if strings.EqualFold(known.Name, spell.Name) {
			return "", fmt.Errorf("spell already known: %s", spell.Name)
		}
	}

	char.Spells = append(char.Spells, *spell)

	if err := s.Repo.Save(ctx, char); err != nil {
		return "", fmt.Errorf("failed to save character: %w", err)
	}

	return fmt.Sprintf("Learned spell %s", spell.Name), nil
}
