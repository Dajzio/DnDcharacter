package domain

import "fmt"

type EquipmentSlot string

const (
	MainHand  EquipmentSlot = "main hand"
	OffHand   EquipmentSlot = "off hand"
	ArmorSlot EquipmentSlot = "armor"
	ShieldSlot EquipmentSlot = "shield"
)

func (c *Character) EquipWeapon(weaponName string, slot EquipmentSlot) string {
	switch slot {
	case MainHand:
		if c.Equipment.MainHandWeapon != nil {
			return fmt.Sprintf("%s already occupied", slot)
		}
		c.Equipment.MainHandWeapon = &Weapon{Name: weaponName}
		return fmt.Sprintf("Equipped weapon %s to %s", weaponName, slot)

	case OffHand:
		if c.Equipment.OffHandWeapon != nil {
			return fmt.Sprintf("%s already occupied", slot)
		}
		c.Equipment.OffHandWeapon = &Weapon{Name: weaponName}
		return fmt.Sprintf("Equipped weapon %s to %s", weaponName, slot)

	default:
		return fmt.Sprintf("Invalid slot: %s", slot)
	}
}

func (c *Character) EquipArmor(armorName string, armorClass int, dexBonus bool) string {
	if c.Equipment.Armor != nil {
		return "armor slot already occupied"
	}
	c.Equipment.Armor = &Armor{
		Name:       armorName,
		ArmorClass: armorClass,
		DexBonus:   dexBonus,
	}

	c.UpdateStats() 
	return fmt.Sprintf("Equipped armor %s", armorName)
}

func (c *Character) EquipShield(shieldName string, armorClass int) string {
	if c.Equipment.Shield != nil {
		return "shield slot already occupied"
	}
	c.Equipment.Shield = &Shield{
		Name:       shieldName,
		ArmorClass: armorClass,
	}

	c.UpdateStats()
	return fmt.Sprintf("Equipped shield %s", shieldName)
}
