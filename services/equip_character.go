package services

import (
	"context"
	"fmt"
	"strings"
	"starter_pack/domain"
)

type EquipItemService struct {
	Repo domain.CharacterRepository
}

func (s *EquipItemService) Execute(ctx context.Context, name, itemType, itemName, slot string) (string, error) {
	char, err := s.Repo.GetByName(ctx, name)
	if err != nil {
		return "", fmt.Errorf("character not found: %w", err)
	}

	itemType = strings.ToLower(itemType)
	itemName = strings.ToLower(itemName)
	slot = strings.ToLower(slot)

	switch itemType {
	case "weapon":
		if slot == "" {
			return "", fmt.Errorf("please specify a slot for the weapon (main hand/off hand)")
		}
		if err := char.EquipWeapon(itemName, slot); err != nil {
			return "", err
		}
	case "armor":
		if err := char.EquipArmor(itemName); err != nil {
			return "", err
		}
	case "shield":
		if err := char.EquipShield(itemName); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unknown item type: %s", itemType)
	}

	if err := s.Repo.Save(ctx, char); err != nil {
		return "", fmt.Errorf("failed to save character: %w", err)
	}

	return fmt.Sprintf("Equipped %s", itemName), nil
}
