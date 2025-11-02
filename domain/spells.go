package domain

import "strings"

type SpellRepository interface {
	LoadFromCSV(path string) error
	GetSpellsForClass(class Class) []Spell
	FindSpellByName(name string) *Spell
	ClassHasSpell(class Class, spellName string) bool
}

type Spell struct {
	Name   string
	Level  int
	Class  []string
	School string `json:"school,omitempty"`
	Range  string `json:"range,omitempty"`
}

func (s *Spell) HasName(name string) bool {
	return strings.EqualFold(s.Name, name)
}

func IsSpellcastingClass(class string) bool {
	_, ok := spellcastingClasses[strings.ToLower(class)]
	return ok
}

func PreparesSpells(class string) bool {
	_, ok := spellPreparationClasses[strings.ToLower(class)]
	return ok
}

var spellcastingClasses = map[string]struct{}{
	"wizard": {}, "sorcerer": {}, "cleric": {}, "paladin": {},
	"druid": {}, "bard": {}, "warlock": {}, "ranger": {},
}

var spellPreparationClasses = map[string]struct{}{
	"cleric": {}, "paladin": {}, "druid": {}, "wizard": {},
}

var spellSlotData = map[string]map[int]map[int]int{
	"cleric": {
		17: {0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1, 7: 1, 8: 1, 9: 1},
		15: {0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1, 7: 1, 8: 0},
		13: {0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 1, 6: 1},
		10: {0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 2},
		9:  {0: 4, 1: 4, 2: 3, 3: 3, 4: 3, 5: 1},
		7:  {0: 4, 1: 4, 2: 3, 3: 3, 4: 1},
		5:  {0: 4, 1: 4, 2: 3, 3: 2},
		3:  {0: 3, 1: 4, 2: 2},
		2:  {0: 3, 1: 3},
		1:  {0: 3, 1: 2},
	},
	"druid":  {}, "sorcerer": {}, "bard": {}, "wizard": {}, 
	"paladin": {}, "ranger": {}, "warlock": {},
}

func init() {
	for _, cls := range []string{"druid", "sorcerer", "bard"} {
		spellSlotData[cls] = copyMap(spellSlotData["cleric"])
	}
	spellSlotData["wizard"] = map[int]map[int]int{
		17: {0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 3, 6: 2, 7: 2, 8: 1, 9: 1},
		15: {0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1, 7: 1, 8: 0},
		13: {0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 1, 6: 1},
		11: {0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1},
		9:  {0: 4, 1: 4, 2: 3, 3: 3, 4: 3, 5: 1},
		7:  {0: 4, 1: 4, 2: 3, 3: 3, 4: 1},
		5:  {0: 4, 1: 4, 2: 3, 3: 2},
		3:  {0: 3, 1: 4, 2: 2},
		2:  {0: 3, 1: 3},
		1:  {0: 3, 1: 2},
	}

	spellSlotData["paladin"] = map[int]map[int]int{
		17: {1: 4, 2: 3, 3: 3, 4: 3, 5: 2},
		13: {1: 4, 2: 3, 3: 3, 4: 2},
		9:  {1: 4, 2: 3, 3: 2},
		5:  {1: 4, 2: 2},
	}
	spellSlotData["ranger"] = spellSlotData["paladin"]

	spellSlotData["warlock"] = map[int]map[int]int{
		17: {0: 4, 5: 4},
		11: {0: 4, 5: 3},
		9:  {0: 3, 5: 2},
		7:  {0: 3, 4: 2},
		5:  {0: 3, 3: 2},
		3:  {0: 2, 2: 2},
		1:  {0: 2, 1: 1},
	}
}

func copyMap(src map[int]map[int]int) map[int]map[int]int {
	dst := map[int]map[int]int{}
	for k, v := range src {
		sub := map[int]int{}
		for sk, sv := range v {
			sub[sk] = sv
		}
		dst[k] = sub
	}
	return dst
}

func GetSpellSlots(class Class, level int) map[int]int {
	className := strings.ToLower(string(class))
	data, ok := spellSlotData[className]
	if !ok {
		return map[int]int{}
	}
	
	for l := level; l >= 1; l-- {
		if slots, exists := data[l]; exists {
			return slots
		}
	}
	return map[int]int{}
}
