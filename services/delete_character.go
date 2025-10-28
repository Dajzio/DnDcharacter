package services

import (
	"context"
	"starter_pack/domain"
)

type DeleteCharacterService struct {
	Repo domain.CharacterRepository
}

func (s *DeleteCharacterService) Execute(ctx context.Context, name string) error {
	return s.Repo.Delete(ctx, name)
}
