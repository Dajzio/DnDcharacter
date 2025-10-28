package services

import (
	"context"
	"fmt"
	"strings"

	"starter_pack/domain"
)

type EquipItemService struct {
	Repo interface {
		GetByName(ctx context.Context, name string) (*domain.Character, error)
		Save(ctx context.Context, c *domain.Character) error
	}
}

func (s *EquipItemService) Execute(ctx context.Context, name string, itemType string, itemName string, slot string) (string, error) {
	char, err := s.Repo.GetByName(ctx, name)
	if err != nil {
		return "", fmt.Errorf("character not found: %w", err)
	}

	itemType = strings.ToLower(itemType)
	itemName = strings.ToLower(itemName)
	slot = strings.ToLower(slot)

	var message string

	switch itemType {
	case "weapon":
		if slot == "" {
			return "", fmt.Errorf("please specify a slot for the weapon (main hand/off hand)")
		}
		message = char.EquipWeapon(itemName, domain.EquipmentSlot(slot))

	case "armor":
		message = char.EquipArmor(itemName, 0, true) 

	case "shield":
		message = char.EquipShield(itemName, 0) 

	default:
		return "", fmt.Errorf("unknown item type: %s", itemType)
	}

	if err := s.Repo.Save(ctx, char); err != nil {
		return "", fmt.Errorf("failed to save character: %w", err)
	}

	return message, nil
}
