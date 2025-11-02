package domain

import (
	"sort"
	"strings"
)

type SkillRepository struct{}

func NewSkillRepository() *SkillRepository {
	return &SkillRepository{}
}


func (r *SkillRepository) GetAllClassSkills(class string) []string {
	classSkills := map[string][]string{
		"barbarian": {"Animal Handling", "Athletics", "Intimidation", "Nature", "Perception", "Survival"},
		"bard":      {"Arcana", "Deception", "Insight", "Intimidation", "Performance", "Persuasion", "Religion"},
		"cleric":    {"History", "Insight", "Medicine", "Persuasion", "Religion"},
		"druid":     {"Arcana", "Animal Handling", "Insight", "Medicine", "Nature", "Perception", "Religion", "Survival"},
		"fighter":   {"Acrobatics", "Animal Handling", "Athletics", "History", "Insight", "Intimidation", "Perception", "Survival"},
		"monk":      {"Acrobatics", "Athletics", "History", "Insight", "Religion", "Stealth"},
		"paladin":   {"Athletics", "Insight", "Intimidation", "Medicine", "Persuasion", "Religion"},
		"ranger":    {"Animal Handling", "Athletics", "Insight", "Investigation", "Nature", "Perception", "Stealth", "Survival"},
		"rogue":     {"Acrobatics", "Athletics", "Deception", "Insight", "Intimidation", "Investigation", "Perception", "Performance", "Persuasion", "Sleight of Hand", "Stealth"},
		"sorcerer":  {"Arcana", "Deception", "Insight", "Intimidation", "Persuasion", "Religion"},
		"warlock":   {"Arcana", "Deception", "History", "Intimidation", "Investigation", "Nature", "Religion"},
		"wizard":    {"Arcana", "History", "Insight", "Investigation", "Medicine", "Religion"},
	}

	skills := classSkills[strings.ToLower(class)]
	sort.Strings(skills)
	return skills
}

func (r *SkillRepository) GetAllBackgroundSkills(background string) []string {
	backgroundSkills := map[string][]string{
		"acolyte":       {"Insight", "Religion"},
		"charlatan":     {"Deception", "Sleight of Hand"},
		"criminal":      {"Deception", "Stealth"},
		"entertainer":   {"Acrobatics", "Performance"},
		"folk hero":     {"Animal Handling", "Survival"},
		"guild artisan": {"Insight", "Persuasion"},
		"hermit":        {"Medicine", "Religion"},
		"noble":         {"History", "Persuasion"},
		"outlander":     {"Athletics", "Survival"},
		"sage":          {"Arcana", "History"},
		"sailor":        {"Athletics", "Perception"},
		"soldier":       {"Athletics", "Intimidation"},
		"urchin":        {"Sleight of Hand", "Stealth"},
	}

	skills := backgroundSkills[strings.ToLower(background)]
	sort.Strings(skills)
	return skills
}

func (r *SkillRepository) GetDefaultSkills(class string, background string) []string {
	classSkills := r.GetAllClassSkills(class)
	backgroundSkills := r.GetAllBackgroundSkills(background)

	classSkillCount := map[string]int{
		"barbarian": 2,
		"bard":      3,
		"cleric":    2,
		"druid":     2,
		"fighter":   2,
		"monk":      2,
		"paladin":   2,
		"ranger":    3,
		"rogue":     4,
		"sorcerer":  2,
		"warlock":   2,
		"wizard":    2,
	}

	count := classSkillCount[strings.ToLower(class)]
	if count == 0 {
		count = 2 
	}

	selected := []string{}

	for i := 0; i < count && i < len(classSkills); i++ {
		selected = append(selected, strings.ToLower(classSkills[i]))
	}

	for i := 0; i < 2 && i < len(backgroundSkills); i++ {
		selected = append(selected, strings.ToLower(backgroundSkills[i]))
	}

	return selected
}
