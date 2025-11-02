package services

import (
	"testing"
	"starter_pack/domain"
)

func TestEquipment(t *testing.T) {
	char := &domain.Character{}

	res := EquipWeapon(char, "Longsword", MainHand)
	if char.Equipment.MainHandWeapon == nil || res != "Equipped weapon Longsword to main hand" {
		t.Fatalf("failed to equip main hand weapon")
	}

	res = EquipWeapon(char, "Dagger", MainHand)
	if res != "main hand already occupied" {
		t.Fatalf("expected occupied error for main hand")
	}

	res = EquipWeapon(char, "Dagger", OffHand)
	if char.Equipment.OffHandWeapon == nil || res != "Equipped weapon Dagger to off hand" {
		t.Fatalf("failed to equip off hand weapon")
	}

	res = EquipWeapon(char, "Axe", "invalid")
	if res != "Invalid slot: invalid" {
		t.Fatalf("expected invalid slot error")
	}

	res = EquipArmor(char, "Leather", 11, true)
	if char.Equipment.Armor == nil || res != "Equipped armor Leather" {
		t.Fatalf("failed to equip armor")
	}

	res = EquipArmor(char, "Chainmail", 16, false)
	if res != "armor slot already occupied" {
		t.Fatalf("expected occupied error for armor")
	}

	res = EquipShield(char, "Wooden Shield", 2)
	if char.Equipment.Shield == nil || res != "Equipped shield Wooden Shield" {
		t.Fatalf("failed to equip shield")
	}

	res = EquipShield(char, "Iron Shield", 2)
	if res != "shield slot already occupied" {
		t.Fatalf("expected occupied error for shield")
	}
}
