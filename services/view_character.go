package services

import (
	"context"
	"fmt"
	"starter_pack/domain"
	"strings"
)

func abilityMod(score int) int {
	return domain.Modifier(score)
}

type ViewCharacterService struct {
	Repo domain.CharacterRepository
}

func (s *ViewCharacterService) Execute(ctx context.Context, name string) error {
	list, err := s.Repo.List(ctx)
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
		return fmt.Errorf("character %q not found", name)
	}

	character.UpdateStats()
	PrintCharacter(character)
	return nil
}

func PrintCharacter(c *domain.Character) {
	fmt.Printf("Name: %s\nClass: %s\nRace: %s\nBackground: %s\nLevel: %d\n\n",
		c.Name, c.Class, c.Race, c.Background, c.Level)

	abilities := map[string]int{
		"STR": c.AbilityScores.Str,
		"DEX": c.AbilityScores.Dex,
		"CON": c.AbilityScores.Con,
		"INT": c.AbilityScores.Int,
		"WIS": c.AbilityScores.Wis,
		"CHA": c.AbilityScores.Cha,
	}

	fmt.Println("Ability scores:")
	for k, v := range abilities {
		fmt.Printf("  %s: %d (%+d)\n", k, v, abilityMod(v))
	}

	fmt.Printf("\nProficiency bonus: +%d\n", c.ProficiencyBonus)

	if len(c.SkillProficiencies) > 0 {
		fmt.Printf("Skill proficiencies: %s\n", strings.Join(c.SkillProficiencies, ", "))
	}

	equipment := map[string]*domain.Weapon{
		"Main hand": c.Equipment.MainHandWeapon,
		"Off hand":  c.Equipment.OffHandWeapon,
	}
	for slot, w := range equipment {
		if w != nil {
			fmt.Printf("%s: %s\n", slot, w.Name)
		}
	}
	if c.Equipment.Armor != nil {
		fmt.Printf("Armor: %s\n", c.Equipment.Armor.Name)
	}
	if c.Equipment.Shield != nil {
		fmt.Printf("Shield: %s\n", c.Equipment.Shield.Name)
	}

	if domain.IsSpellcastingClass(string(c.Class)) && len(c.SpellSlots) > 0 {
		fmt.Println("Spell slots:")
		for lvl := 0; lvl <= 9; lvl++ {
			if count, ok := c.SpellSlots[lvl]; ok && count > 0 {
				if lvl == 0 {
					fmt.Printf("  Level 0: %d\n", count)
				} else {
					fmt.Printf("  Level %d: %d\n", lvl, count)
				}
			}
		}

		fullName := FullAbilityName(c.SpellcastingAbility)
		fmt.Printf("Spellcasting ability: %s\nSpell save DC: %d\nSpell attack bonus: +%d\n",
			fullName, c.SpellSaveDC, c.SpellAttackBonus)
	}

	fmt.Printf("\nArmor class: %d\nInitiative bonus: %d\nPassive perception: %d\n",
		c.ArmorClass, c.Initiative, c.PassivePerception)
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
