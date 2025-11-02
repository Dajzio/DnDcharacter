package services

import (
	"context"
	"fmt"
	"starter_pack/domain"
)

type CreateCharacterInput struct {
	Name       string
	Race       domain.Race
	Class      domain.Class
	Background string
	Level      int
	Str        int
	Dex        int
	Con        int
	Int        int
	Wis        int
	Cha        int
	Skills     []string
}

type CreateCharacterService struct {
	Repo domain.CharacterRepository
	Factory *domain.CharacterFactory
}

func (s *CreateCharacterService) Execute(ctx context.Context, input CreateCharacterInput) (*domain.Character, error) {
	ab := domain.AbilityScores{
		Str: input.Str, Dex: input.Dex, Con: input.Con,
		Int: input.Int, Wis: input.Wis, Cha: input.Cha,
	}

	char, err := s.Factory.Create(
		domain.GenerateID(),
		input.Name,
		input.Race,
		input.Class,
		input.Level,
		ab,
		input.Background,
		input.Skills,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create character: %w", err)
	}

	if err := s.Repo.Save(ctx, char); err != nil {
		return nil, fmt.Errorf("cannot save character: %w", err)
	}

	return char, nil
}
