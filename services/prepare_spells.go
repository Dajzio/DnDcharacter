package services

import (
	"context"
	"fmt"

	"starter_pack/domain"
)

type PrepareSpellService struct {
	Repo    domain.CharacterRepository
	SpellRepo domain.SpellRepository
}

func (s *PrepareSpellService) Execute(ctx context.Context, name string, spellName string) (string, error) {
	char, err := s.Repo.GetByName(ctx, name)
	if err != nil {
		return "", fmt.Errorf("character not found: %w", err)
	}

	spell := s.SpellRepo.FindSpellByName(spellName)
	if spell == nil {
		return "", fmt.Errorf("spell not found: %s", spellName)
	}

	if err := char.PrepareSpell(*spell); err != nil {
		return "", err
	}

	if err := s.Repo.Save(ctx, char); err != nil {
		return "", fmt.Errorf("failed to save character: %w", err)
	}

	return fmt.Sprintf("Prepared spell %s", spell.Name), nil
}
