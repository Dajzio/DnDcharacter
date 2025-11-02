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
	switch strings.ToLower(class) {
	case "wizard", "sorcerer", "cleric", "paladin", "druid", "bard", "warlock", "ranger":
		return true
	default:
		return false
	}
}

func PreparesSpells(class string) bool {
	switch strings.ToLower(class) {
	case "cleric", "paladin", "druid", "wizard":
		return true
	default:
		return false
	}
}

func GetSpellSlots(class Class, level int) map[int]int {
	className := strings.ToLower(string(class))
	slots := map[int]int{}

	switch className {
	case "cleric":
		switch {
		case level >= 17:
			slots = map[int]int{0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1, 7: 1, 8: 1, 9: 1}
		case level >= 15:
			slots = map[int]int{0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1, 7: 1, 8: 0}
		case level >= 13:
			slots = map[int]int{0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 1, 6: 1}
		case level >= 10:
			slots = map[int]int{0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 2}
		case level >= 9:
			slots = map[int]int{0: 4, 1: 4, 2: 3, 3: 3, 4: 3, 5: 1}
		case level >= 7:
			slots = map[int]int{0: 4, 1: 4, 2: 3, 3: 3, 4: 1}
		case level >= 5:
			slots = map[int]int{0: 4, 1: 4, 2: 3, 3: 2}
		case level >= 3:
			slots = map[int]int{0: 3, 1: 4, 2: 2}
		case level >= 2:
			slots = map[int]int{0: 3, 1: 3}
		default:
			slots = map[int]int{0: 3, 1: 2}
		}

	case "druid", "sorcerer", "bard":
		switch {
		case level >= 17:
			slots = map[int]int{0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1, 7: 1, 8: 1, 9: 1}
		case level >= 15:
			slots = map[int]int{0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1, 7: 1, 8: 0}
		case level >= 13:
			slots = map[int]int{0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 1, 6: 1}
		case level >= 11:
			slots = map[int]int{0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1}
		case level >= 9:
			slots = map[int]int{0: 4, 1: 4, 2: 3, 3: 3, 4: 3, 5: 1}
		case level >= 7:
			slots = map[int]int{0: 4, 1: 4, 2: 3, 3: 3, 4: 1}
		case level >= 5:
			slots = map[int]int{0: 4, 1: 4, 2: 3, 3: 2}
		case level >= 3:
			slots = map[int]int{0: 3, 1: 4, 2: 2}
		case level >= 2:
			slots = map[int]int{0: 3, 1: 3}
		default:
			slots = map[int]int{0: 3, 1: 2}
		}

	case "wizard":
		switch {
		case level >= 17:
			slots = map[int]int{0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 3, 6: 2, 7: 2, 8: 1, 9: 1}
		case level >= 15:
			slots = map[int]int{0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1, 7: 1, 8: 0}
		case level >= 13:
			slots = map[int]int{0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 1, 6: 1}
		case level >= 11:
			slots = map[int]int{0: 5, 1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1}
		case level >= 9:
			slots = map[int]int{0: 4, 1: 4, 2: 3, 3: 3, 4: 3, 5: 1}
		case level >= 7:
			slots = map[int]int{0: 4, 1: 4, 2: 3, 3: 3, 4: 1}
		case level >= 5:
			slots = map[int]int{0: 4, 1: 4, 2: 3, 3: 2}
		case level >= 3:
			slots = map[int]int{0: 3, 1: 4, 2: 2}
		case level >= 2:
			slots = map[int]int{0: 3, 1: 3}
		default:
			slots = map[int]int{0: 3, 1: 2}
		}

	case "paladin", "ranger":
		switch {
		case level >= 17:
			slots[1] = 4
			slots[2] = 3
			slots[3] = 3
			slots[4] = 3
			slots[5] = 2
		case level >= 13:
			slots[1] = 4
			slots[2] = 3
			slots[3] = 3
			slots[4] = 2
		case level >= 9:
			slots[1] = 4
			slots[2] = 3
			slots[3] = 2
		case level >= 5:
			slots[1] = 4
			slots[2] = 2
		}

	case "warlock":
		switch {
		case level >= 17:
			slots = map[int]int{0: 4, 5: 4}
		case level >= 11:
			slots = map[int]int{0: 4, 5: 3}
		case level >= 9:
			slots = map[int]int{0: 3, 5: 2}
		case level >= 7:
			slots = map[int]int{0: 3, 4: 2}
		case level >= 5:
			slots = map[int]int{0: 3, 3: 2}
		case level >= 3:
			slots = map[int]int{0: 2, 2: 2}
		default:
			slots = map[int]int{0: 2, 1: 1}
		}

	default:
		return map[int]int{}
	}

	return slots
}
