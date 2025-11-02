package domain

import (
	"context"
)

type CharacterRepository interface {
	Save(context context.Context, c *Character) error
	GetByID(context context.Context, id string) (*Character, error)
	GetByName(context context.Context, name string) (*Character, error)
	List(context context.Context) ([]*Character, error)
	Delete(ctx context.Context, name string) error
}
