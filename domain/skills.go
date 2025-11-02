package domain

import (
	"sort"
	"strings"
)

type SkillRepository struct{}

func NewSkillRepository() *SkillRepository {
	return &SkillRepository{}
}

const (
	SkillAnimalHandling = "Animal Handling"
	SkillAthletics      = "Athletics"
	SkillIntimidation   = "Intimidation"
	SkillNature         = "Nature"
	SkillPerception     = "Perception"
	SkillSurvival       = "Survival"
	SkillArcana         = "Arcana"
	SkillDeception      = "Deception"
	SkillInsight        = "Insight"
	SkillPerformance    = "Performance"
	SkillPersuasion     = "Persuasion"
	SkillReligion       = "Religion"
	SkillHistory        = "History"
	SkillMedicine       = "Medicine"
	SkillInvestigation  = "Investigation"
	SkillStealth        = "Stealth"
	SkillSleightOfHand  = "Sleight of Hand"
	SkillAcrobatics     = "Acrobatics"
)

func (r *SkillRepository) GetAllClassSkills(class string) []string {
	classSkills := map[string][]string{
		"barbarian": {SkillAnimalHandling, SkillAthletics, SkillIntimidation, SkillNature, SkillPerception, SkillSurvival},
		"bard":      {SkillArcana, SkillDeception, SkillInsight, SkillIntimidation, SkillPerformance, SkillPersuasion, SkillReligion},
		"cleric":    {SkillHistory, SkillInsight, SkillMedicine, SkillPersuasion, SkillReligion},
		"druid":     {SkillArcana, SkillAnimalHandling, SkillInsight, SkillMedicine, SkillNature, SkillPerception, SkillReligion, SkillSurvival},
		"fighter":   {SkillAcrobatics, SkillAnimalHandling, SkillAthletics, SkillHistory, SkillInsight, SkillIntimidation, SkillPerception, SkillSurvival},
		"monk":      {SkillAcrobatics, SkillAthletics, SkillHistory, SkillInsight, SkillReligion, SkillStealth},
		"paladin":   {SkillAthletics, SkillInsight, SkillIntimidation, SkillMedicine, SkillPersuasion, SkillReligion},
		"ranger":    {SkillAnimalHandling, SkillAthletics, SkillInsight, SkillInvestigation, SkillNature, SkillPerception, SkillStealth, SkillSurvival},
		"rogue":     {SkillAcrobatics, SkillAthletics, SkillDeception, SkillInsight, SkillIntimidation, SkillInvestigation, SkillPerception, SkillPerformance, SkillPersuasion, SkillSleightOfHand, SkillStealth},
		"sorcerer":  {SkillArcana, SkillDeception, SkillInsight, SkillIntimidation, SkillPersuasion, SkillReligion},
		"warlock":   {SkillArcana, SkillDeception, SkillHistory, SkillIntimidation, SkillInvestigation, SkillNature, SkillReligion},
		"wizard":    {SkillArcana, SkillHistory, SkillInsight, SkillInvestigation, SkillMedicine, SkillReligion},
	}

	skills := classSkills[strings.ToLower(class)]
	sort.Strings(skills)
	return skills
}

func (r *SkillRepository) GetAllBackgroundSkills(background string) []string {
	backgroundSkills := map[string][]string{
		"acolyte":       {SkillInsight, SkillReligion},
		"charlatan":     {SkillDeception, SkillSleightOfHand},
		"criminal":      {SkillDeception, SkillStealth},
		"entertainer":   {SkillAcrobatics, SkillPerformance},
		"folk hero":     {SkillAnimalHandling, SkillSurvival},
		"guild artisan": {SkillInsight, SkillPersuasion},
		"hermit":        {SkillMedicine, SkillReligion},
		"noble":         {SkillHistory, SkillPersuasion},
		"outlander":     {SkillAthletics, SkillSurvival},
		"sage":          {SkillArcana, SkillHistory},
		"sailor":        {SkillAthletics, SkillPerception},
		"soldier":       {SkillAthletics, SkillIntimidation},
		"urchin":        {SkillSleightOfHand, SkillStealth},
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
