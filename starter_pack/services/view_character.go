package services

import (
	"context"
	"fmt"
	"starter_pack/domain"
	"starter_pack/infrastructure"
	"strings"
)

func abilityMod(score int) int {
	return domain.Modifier(score)
}

func ViewCharacter(filename, name string) error {
	repo := infrastructure.NewFileCharacterRepo(filename)

	list, err := repo.List(context.Background())
	if err != nil {
		return err
	}

	var character *domain.Character
	for _, c := range list {
		if strings.EqualFold(c.Name, name) {
			character = c
			break
		}
	}

	if character == nil {
		return fmt.Errorf("character \"%s\" not found", name)
	}

	character.UpdateStats()

	PrintCharacter(character)
	return nil
}

func PrintCharacter(c *domain.Character) {
	fmt.Printf("Name: %s\n", c.Name)
	fmt.Printf("Class: %s\n", c.Class)
	fmt.Printf("Race: %s\n", c.Race)
	fmt.Printf("Background: %s\n", c.Background)
	fmt.Printf("Level: %d\n", c.Level)

	fmt.Println("Ability scores:")
	fmt.Printf("  STR: %d (%+d)\n", c.AbilityScores.Str, abilityMod(c.AbilityScores.Str))
	fmt.Printf("  DEX: %d (%+d)\n", c.AbilityScores.Dex, abilityMod(c.AbilityScores.Dex))
	fmt.Printf("  CON: %d (%+d)\n", c.AbilityScores.Con, abilityMod(c.AbilityScores.Con))
	fmt.Printf("  INT: %d (%+d)\n", c.AbilityScores.Int, abilityMod(c.AbilityScores.Int))
	fmt.Printf("  WIS: %d (%+d)\n", c.AbilityScores.Wis, abilityMod(c.AbilityScores.Wis))
	fmt.Printf("  CHA: %d (%+d)\n", c.AbilityScores.Cha, abilityMod(c.AbilityScores.Cha))

	fmt.Printf("Proficiency bonus: +%d\n", c.ProficiencyBonus)
	if len(c.SkillProficiencies) > 0 {
		fmt.Printf("Skill proficiencies: %s\n", strings.Join(c.SkillProficiencies, ", "))
	}

	if c.Equipment.MainHandWeapon != nil {
		fmt.Printf("Main hand: %s\n", c.Equipment.MainHandWeapon.Name)
	}
	if c.Equipment.OffHandWeapon != nil {
		fmt.Printf("Off hand: %s\n", c.Equipment.OffHandWeapon.Name)
	}
	if c.Equipment.Armor != nil {
		fmt.Printf("Armor: %s\n", c.Equipment.Armor.Name)
	}
	if c.Equipment.Shield != nil {
		fmt.Printf("Shield: %s\n", c.Equipment.Shield.Name)
	}

	if domain.IsSpellcastingClass(string(c.Class)) && len(c.SpellSlots) > 0 {
		if len(c.SpellSlots) > 0 {
			fmt.Println("Spell slots:")
			for lvl := 0; lvl <= 9; lvl++ {
				if count, ok := c.SpellSlots[lvl]; ok {
					if lvl == 0 {
						fmt.Printf("  Level 0: %d\n", count)
					} else if count > 0 {
						fmt.Printf("  Level %d: %d\n", lvl, count)
					}
				}
			}
		}

		if domain.IsSpellcastingClass(string(c.Class)) {
			fullName := FullAbilityName(c.SpellcastingAbility)
			fmt.Printf("Spellcasting ability: %s\n", fullName)
			fmt.Printf("Spell save DC: %d\n", c.SpellSaveDC)
			fmt.Printf("Spell attack bonus: +%d\n", c.SpellAttackBonus)
		}
	}

	fmt.Printf("Armor class: %d\n", c.ArmorClass)
	fmt.Printf("Initiative bonus: %d\n", c.Initiative)
	fmt.Printf("Passive perception: %d\n", c.PassivePerception)

}

func FullAbilityName(short string) string {
	switch strings.ToUpper(short) {
	case "STR":
		return "strength"
	case "DEX":
		return "dexterity"
	case "CON":
		return "constitution"
	case "INT":
		return "intelligence"
	case "WIS":
		return "wisdom"
	case "CHA":
		return "charisma"
	default:
		return strings.ToLower(short)
	}
}
