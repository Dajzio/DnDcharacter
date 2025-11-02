package services

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"starter_pack/domain"
)

type CharacterSheetService struct {
	Repo domain.CharacterRepository
}

func (s *CharacterSheetService) Execute(ctx context.Context, name string, format string) (string, error) {
	char, err := s.Repo.GetByName(ctx, name)
	if err != nil {
		return "", fmt.Errorf("character not found: %w", err)
	}

	if strings.ToLower(format) != "markdown" {
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", char.Name))

	sb.WriteString("## Character\n")
	sb.WriteString(fmt.Sprintf("Class: %s\n", char.Class))
	sb.WriteString(fmt.Sprintf("Race: %s\n", char.Race))
	sb.WriteString(fmt.Sprintf("Background: %s\n", char.Background))
	sb.WriteString(fmt.Sprintf("Level: %d\n", char.Level))
	sb.WriteString(fmt.Sprintf("Proficiency bonus: +%d\n", char.ProficiencyBonus))
	sb.WriteString(fmt.Sprintf("Passive perception: %d\n\n", char.PassivePerception))

	sb.WriteString("## Ability scores\n")
	sb.WriteString(fmt.Sprintf("STR: %d (%+d)\n", char.AbilityScores.Str, domain.Modifier(char.AbilityScores.Str)))
	sb.WriteString(fmt.Sprintf("DEX: %d (%+d)\n", char.AbilityScores.Dex, domain.Modifier(char.AbilityScores.Dex)))
	sb.WriteString(fmt.Sprintf("CON: %d (%+d)\n", char.AbilityScores.Con, domain.Modifier(char.AbilityScores.Con)))
	sb.WriteString(fmt.Sprintf("INT: %d (%+d)\n", char.AbilityScores.Int, domain.Modifier(char.AbilityScores.Int)))
	sb.WriteString(fmt.Sprintf("WIS: %d (%+d)\n", char.AbilityScores.Wis, domain.Modifier(char.AbilityScores.Wis)))
	sb.WriteString(fmt.Sprintf("CHA: %d (%+d)\n\n", char.AbilityScores.Cha, domain.Modifier(char.AbilityScores.Cha)))

	sb.WriteString("## Skills\n")
	repo := domain.NewSkillRepository()
	for _, skill := range repo.AllSkills() {
		marked := "[]"
		if repo.HasSkill(char.SkillProficiencies, skill.Name) {
			marked = "[x]"
		}
		sb.WriteString(fmt.Sprintf("%s %s (%s)\n", marked, skill.Name, skill.Ability))
	}
	sb.WriteString("\n")

	sb.WriteString("## Equipment\n")
	if char.Equipment.MainHandWeapon != nil {
		sb.WriteString(fmt.Sprintf("Main hand: %s\n", char.Equipment.MainHandWeapon.Name))
	}
	if char.Equipment.Armor != nil {
		sb.WriteString(fmt.Sprintf("Armor: %s\n", char.Equipment.Armor.Name))
	}
	if char.Equipment.Shield != nil {
		sb.WriteString(fmt.Sprintf("Shield: %s\n", char.Equipment.Shield.Name))
	}
	sb.WriteString("\n")

	sb.WriteString("## Combat stats\n")
	sb.WriteString(fmt.Sprintf("Armor class: %d\n", char.ArmorClass))
	sb.WriteString(fmt.Sprintf("Initiative bonus: %+d\n\n", char.Initiative))

	if len(char.SpellSlots) > 0 {
		sb.WriteString("## Spell slots [leave empty on non-casters]\n")
		keys := make([]int, 0, len(char.SpellSlots))
		for k := range char.SpellSlots {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, lvl := range keys {
			sb.WriteString(fmt.Sprintf("Level %d: %d\n", lvl, char.SpellSlots[lvl]))
		}
		sb.WriteString("\n")
	}

	if char.SpellcastingAbility != "" {
		sb.WriteString("## Spellcasting [leave empty on non-casters]\n")
		sb.WriteString(fmt.Sprintf("Spellcasting ability: %s\n", FullAbilityName(char.SpellcastingAbility)))
		sb.WriteString(fmt.Sprintf("Spell save DC: %d\n", char.SpellSaveDC))
		sb.WriteString(fmt.Sprintf("Spell attack bonus: %+d\n\n", char.SpellAttackBonus))
	}

	if len(char.Spells) > 0 {
		sb.WriteString("## Spells [leave empty on non-casters, only add levels that contain spells]\n\n")
		spellsByLevel := map[int][]string{}
		for _, sp := range char.Spells {
			spellsByLevel[sp.Level] = append(spellsByLevel[sp.Level], sp.Name)
		}

		levels := make([]int, 0, len(spellsByLevel))
		for lvl := range spellsByLevel {
			levels = append(levels, lvl)
		}
		sort.Ints(levels)

		for _, lvl := range levels {
			sb.WriteString(fmt.Sprintf("### Level %d\n", lvl))
			for _, spName := range spellsByLevel[lvl] {
				sb.WriteString(fmt.Sprintf("- %s\n", spName))
			}
			sb.WriteString("\n")
		}
	}

	return sb.String(), nil
}
