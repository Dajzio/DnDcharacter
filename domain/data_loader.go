package domain

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CachedSpell struct {
	Name   string `json:"name"`
	School string `json:"school"`
	Range  string `json:"range"`
	Level  int    `json:"level"`
}

type CachedWeapon struct {
	Name      string `json:"name"`
	Category  string `json:"category"`
	Range     string `json:"range"`
	TwoHanded bool   `json:"two_handed"`
}

type CachedArmor struct {
	Name           string `json:"name"`
	ArmorClass     int    `json:"armor_class"`
	DexterityBonus bool   `json:"dex_bonus"`
}

type EquipmentCache struct {
	Weapons []CachedWeapon `json:"weapons"`
	Armors  []CachedArmor  `json:"armors"`
}


func LoadCachedSpells() ([]CachedSpell, error) {
	data, err := os.ReadFile("spells_data.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []CachedSpell{}, nil
		}
		return nil, fmt.Errorf("cannot read spells_data.json: %w", err)
	}
	if len(data) == 0 {
		return []CachedSpell{}, nil
	}
	var spells []CachedSpell
	if err := json.Unmarshal(data, &spells); err != nil {
		return nil, fmt.Errorf("invalid spell data: %w", err)
	}
	return spells, nil
}

func LoadCachedEquipment() (EquipmentCache, error) {
	data, err := os.ReadFile("eq_data.json")
	if err != nil {
		if os.IsNotExist(err) {
			return EquipmentCache{}, nil
		}
		return EquipmentCache{}, fmt.Errorf("cannot read eq_data.json: %w", err)
	}
	if len(data) == 0 {
		return EquipmentCache{}, nil
	}
	var eq EquipmentCache
	if err := json.Unmarshal(data, &eq); err != nil {
		return EquipmentCache{}, fmt.Errorf("invalid equipment data: %w", err)
	}
	return eq, nil
}


func FindCachedWeapon(name string) *CachedWeapon {
	eq, _ := LoadCachedEquipment()
	for _, w := range eq.Weapons {
		if strings.EqualFold(w.Name, name) {
			return &w
		}
	}
	return &CachedWeapon{
		Name:      name,
		Category:  "unknown",
		Range:     "0",
		TwoHanded: false,
	}
}

func FindCachedArmor(name string) *CachedArmor {
	eq, _ := LoadCachedEquipment()
	for _, a := range eq.Armors {
		if strings.EqualFold(a.Name, name) {
			return &a
		}
	}
	return &CachedArmor{
		Name:           name,
		ArmorClass:     0,
		DexterityBonus: false,
	}
}


func FindCachedSpell(name string) *CachedSpell {
	name = strings.ToLower(strings.TrimSpace(name))

	spells, _ := LoadCachedSpells()
	for _, s := range spells {
		if strings.EqualFold(s.Name, name) && s.Level > 0 {
			return &s
		}
	}

	level, school, rng := loadSpellFromCSV(name)
	if level == 0 {
		return nil
	}

	return &CachedSpell{
		Name:   name,
		School: school,
		Range:  rng,
		Level:  level,
	}
}


func loadSpellFromCSV(name string) (int, string, string) {
	file, err := os.Open("5e-SRD-Spells.csv")
	if err != nil {
		return 0, "", ""
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, "", ""
	}

	name = strings.ToLower(strings.TrimSpace(name))
	for _, fields := range records {
		if len(fields) < 4 {
			continue
		}
		spellName := strings.ToLower(strings.TrimSpace(fields[0]))
		if spellName == name {
			lvl, _ := strconv.Atoi(strings.TrimSpace(fields[1]))
			school := strings.TrimSpace(fields[2])
			rng := strings.TrimSpace(fields[3])
			return lvl, school, rng
		}
	}

	return 0, "", ""
}
