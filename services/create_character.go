package services

import (
	"context"
	"fmt"
	"math/rand"
	"starter_pack/domain"
	"time"
)

type CreateCharacterService struct {
	Repo domain.CharacterRepository
}

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

func (s *CreateCharacterService) Execute(ctx context.Context, in CreateCharacterInput) (*domain.Character, error) {
	ab := domain.AbilityScores{
		Str: in.Str,
		Dex: in.Dex,
		Con: in.Con,
		Int: in.Int,
		Wis: in.Wis,
		Cha: in.Cha,
	}

	id := generateID()
	c, err := domain.NewCharacter(id, in.Name, in.Race, in.Class, in.Level, ab, in.Background, in.Skills)
	if err != nil {
		return nil, err
	}

	if err := s.Repo.Save(ctx, c); err != nil {
		return nil, fmt.Errorf("save failed: %w", err)
	}

	return c, nil
}

func generateID() string {
	rand.Seed(time.Now().UnixNano())
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}


