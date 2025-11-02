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
	printBasicInfo(c)
	printAbilities(c)
	printProficiencies(c)
	printEquipment(c)
	printSpells(c)
	printCombatStats(c)
}

func printBasicInfo(c *domain.Character) {
	fmt.Printf("Name: %s\nClass: %s\nRace: %s\nBackground: %s\nLevel: %d\n\n",
		c.Name, c.Class, c.Race, c.Background, c.Level)
}

func printAbilities(c *domain.Character) {
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
}

func printProficiencies(c *domain.Character) {
	if len(c.SkillProficiencies) > 0 {
		fmt.Printf("Skill proficiencies: %s\n", strings.Join(c.SkillProficiencies, ", "))
	}
}

func printEquipment(c *domain.Character) {
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
}

func printSpells(c *domain.Character) {
	if !domain.IsSpellcastingClass(string(c.Class)) || len(c.SpellSlots) == 0 {
		return
	}

	fmt.Println("Spell slots:")
	for lvl := 0; lvl <= 9; lvl++ {
		if count, ok := c.SpellSlots[lvl]; ok && count > 0 {
			label := fmt.Sprintf("Level %d", lvl)
			if lvl == 0 {
				label = "Level 0"
			}
			fmt.Printf("  %s: %d\n", label, count)
		}
	}

	fullName := FullAbilityName(c.SpellcastingAbility)
	fmt.Printf("Spellcasting ability: %s\nSpell save DC: %d\nSpell attack bonus: +%d\n",
		fullName, c.SpellSaveDC, c.SpellAttackBonus)
}

func printCombatStats(c *domain.Character) {
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
