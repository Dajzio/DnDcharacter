package domain

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

type Race string
type Class string

type AbilityScores struct {
	Str int
	Dex int
	Con int
	Int int
	Wis int
	Cha int
}

type Character struct {
	ID                  string
	Name                string
	Race                Race
	Class               Class
	Background          string
	Level               int
	AbilityScores       AbilityScores
	SkillProficiencies  []string
	ProficiencyBonus    int
	Equipment           Equipment
	Spells              []Spell
	SpellSlots          map[int]int `json:"spell_slots,omitempty"`
	SpellcastingAbility string
	SpellSaveDC         int
	SpellAttackBonus    int
	ArmorClass          int `json:"armor_class,omitempty"`
	Initiative          int `json:"initiative,omitempty"`
	PassivePerception   int `json:"passive_perception,omitempty"`
}

type Equipment struct {
	MainHandWeapon *Weapon
	OffHandWeapon  *Weapon
	Armor          *Armor
	Shield         *Shield
}
type Weapon struct {
	Name      string
	Category  string
	Range     int
	TwoHanded bool
}

type Armor struct {
	Name       string
	ArmorClass int
	DexBonus   bool
}

type Shield struct {
	Name       string
	ArmorClass int
}

type CharacterFactory struct{}

func NewCharacterFactory() *CharacterFactory {
	return &CharacterFactory{}
}

func (f *CharacterFactory) Create(id, name string, race Race, class Class, level int, ab AbilityScores, background string, skills []string) (*Character, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if level < 1 {
		return nil, fmt.Errorf("level cannot be lower than 1")
	}

	racialBonuses := GetRacialBonuses(race)
	ab.Str += racialBonuses["Str"]
	ab.Dex += racialBonuses["Dex"]
	ab.Con += racialBonuses["Con"]
	ab.Int += racialBonuses["Int"]
	ab.Wis += racialBonuses["Wis"]
	ab.Cha += racialBonuses["Cha"]

	char := &Character{
		ID:                 id,
		Name:               name,
		Race:               race,
		Class:              class,
		Level:              level,
		AbilityScores:      ab,
		Background:         background,
		SkillProficiencies: skills,
		ProficiencyBonus:   CalculateProficiencyBonus(level),
	}

	char.UpdateStats()

	return char, nil
}

func (c *Character) UpdateStats() {
	c.Initiative = Modifier(c.AbilityScores.Dex)
	c.ArmorClass = c.CalculateArmorClass()
	c.PassivePerception = 10 + Modifier(c.AbilityScores.Wis)
	c.UpdateSpellcasting()
}

func (c *Character) UpdateSpellcasting() {
	if !IsSpellcastingClass(string(c.Class)) {
		c.SpellSlots = nil
		c.SpellcastingAbility = ""
		c.SpellSaveDC = 0
		c.SpellAttackBonus = 0
		return
	}

	var ability string
	switch strings.ToLower(string(c.Class)) {
	case "wizard":
		ability = "INT"
	case "cleric", "druid", "ranger":
		ability = "WIS"
	case "paladin", "sorcerer", "bard", "warlock":
		ability = "CHA"
	default:
		ability = "INT"
	}

	c.SpellcastingAbility = ability

	var mod int
	switch ability {
	case "STR":
		mod = Modifier(c.AbilityScores.Str)
	case "DEX":
		mod = Modifier(c.AbilityScores.Dex)
	case "CON":
		mod = Modifier(c.AbilityScores.Con)
	case "INT":
		mod = Modifier(c.AbilityScores.Int)
	case "WIS":
		mod = Modifier(c.AbilityScores.Wis)
	case "CHA":
		mod = Modifier(c.AbilityScores.Cha)
	}

	c.SpellSaveDC = 8 + c.ProficiencyBonus + mod
	c.SpellAttackBonus = c.ProficiencyBonus + mod

	c.SpellSlots = GetSpellSlots(c.Class, c.Level)
}

func (c *Character) CalculateArmorClass() int {
	dexMod := Modifier(c.AbilityScores.Dex)
	baseAC := 10
	maxDex := -1

	if c.Equipment.Armor != nil {
		switch strings.ToLower(c.Equipment.Armor.Name) {
		case "leather", "leather armor":
			baseAC = 11
			maxDex = -1
		case "studded leather":
			baseAC = 12
			maxDex = -1
		case "chain shirt":
			baseAC = 13
			maxDex = 2
		case "scale mail":
			baseAC = 14
			maxDex = 2
		case "half plate":
			baseAC = 15
			maxDex = 2
		case "plate armor":
			baseAC = 18
			maxDex = 0
		default:
			baseAC = 10
			maxDex = -1
		}

		if maxDex == -1 {
			baseAC += dexMod
		} else if maxDex > 0 {
			if dexMod > maxDex {
				baseAC += maxDex
			} else {
				baseAC += dexMod
			}
		}
	} else {
		switch strings.ToLower(string(c.Class)) {
		case "barbarian":
			baseAC = 10 + dexMod + Modifier(c.AbilityScores.Con)
		case "monk":
			baseAC = 10 + dexMod + Modifier(c.AbilityScores.Wis)
		default:
			baseAC = 10 + dexMod
		}
	}

	if c.Equipment.Shield != nil {
		baseAC += 2
	}

	return baseAC
}

func GetRacialBonuses(race Race) map[string]int {
	switch strings.ToLower(string(race)) {
	case "human":
		return map[string]int{"Str": 1, "Dex": 1, "Con": 1, "Int": 1, "Wis": 1, "Cha": 1}
	case "variant human":
		return map[string]int{"Str": 1, "Dex": 1}
	case "hill dwarf":
		return map[string]int{"Con": 2, "Wis": 1}
	case "dwarf":
		return map[string]int{"Con": 2}
	case "dwarf mountain":
		return map[string]int{"Con": 2, "Str": 2}
	case "elf high":
		return map[string]int{"Dex": 2, "Int": 1}
	case "elf wood":
		return map[string]int{"Dex": 2, "Wis": 1}
	case "elf drow":
		return map[string]int{"Dex": 2, "Cha": 1}
	case "lightfoot halfling":
		return map[string]int{"Dex": 2, "Cha": 1}
	case "stout halfling":
		return map[string]int{"Dex": 2, "Con": 1}
	case "dragonborn":
		return map[string]int{"Str": 2, "Cha": 1}
	case "gnome forest":
		return map[string]int{"Int": 2, "Dex": 1}
	case "gnome rock":
		return map[string]int{"Int": 2, "Con": 1}
	case "gnome":
		return map[string]int{"Int": 2}
	case "half elf":
		return map[string]int{"Cha": 2, "Con": 1, "Dex": 1}
	case "half orc":
		return map[string]int{"Str": 2, "Con": 1}
	case "tiefling":
		return map[string]int{"Cha": 2, "Int": 1}
	case "aasimar":
		return map[string]int{"Cha": 2, "Con": 1}
	case "firbolg":
		return map[string]int{"Wis": 2, "Str": 1}
	case "goliath":
		return map[string]int{"Str": 2, "Con": 1}
	case "tabaxi":
		return map[string]int{"Dex": 2, "Cha": 1}
	case "triton":
		return map[string]int{"Str": 1, "Con": 1, "Cha": 1}
	default:
		return map[string]int{}
	}
}

func CalculateProficiencyBonus(level int) int {
	switch {
	case level >= 17:
		return 6
	case level >= 13:
		return 5
	case level >= 9:
		return 4
	case level >= 5:
		return 3
	default:
		return 2
	}
}

func Modifier(score int) int {
	return int(math.Floor(float64(score-10) / 2.0))
}

func GenerateID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Int63())
}

func (c *Character) LearnSpell(spell Spell) error {
	for _, s := range c.Spells {
		if s.Name == spell.Name {
			return fmt.Errorf("spell already known: %s", spell.Name)
		}
	}
	c.Spells = append(c.Spells, spell)
	return nil
}

func (c *Character) PrepareSpell(spell Spell) error {
	if !PreparesSpells(string(c.Class)) {
		return fmt.Errorf("this class cannot prepare spells")
	}

	slots := GetSpellSlots(c.Class, c.Level)
	maxSlotLevel := 0
	for lvl := range slots {
		if lvl > maxSlotLevel {
			maxSlotLevel = lvl
		}
	}

	if spell.Level > maxSlotLevel {
		return fmt.Errorf("the spell has higher level than the available spell slots")
	}

	if count, ok := slots[spell.Level]; !ok || count == 0 {
		return fmt.Errorf("no slots available for this spell level")
	}

	for _, s := range c.Spells {
		if s.Name == spell.Name {
			return fmt.Errorf("spell already prepared: %s", spell.Name)
		}
	}

	c.Spells = append(c.Spells, spell)
	c.UpdateStats()
	return nil
}

func (c *Character) EquipWeapon(name string, slot string) error {
	slot = strings.ToLower(slot)
	if slot != "main hand" && slot != "off hand" {
		return fmt.Errorf("invalid weapon slot: %s", slot)
	}
	weapon := &Weapon{Name: name}
	if slot == "main hand" {
		c.Equipment.MainHandWeapon = weapon
	} else {
		c.Equipment.OffHandWeapon = weapon
	}
	c.UpdateStats()
	return nil
}

func (c *Character) EquipArmor(name string) error {
	c.Equipment.Armor = &Armor{Name: name}
	c.UpdateStats()
	return nil
}

func (c *Character) EquipShield(name string) error {
	c.Equipment.Shield = &Shield{Name: name}
	c.UpdateStats()
	return nil
}
