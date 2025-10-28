package domain

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Spell struct {
    Name  string
    Level int
    Class []string
    School string `json:"school,omitempty"`
    Range  string `json:"range,omitempty"`
}


var ClassSpellList = map[string][]Spell{}

func LoadSpellsFromCSV(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open spells CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.ReuseRecord = true

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV: %w", err)
		}

		if len(record) < 3 {
			continue
		}

		name := strings.ToLower(strings.TrimSpace(record[0]))
		levelStr := strings.TrimSpace(record[1])
		classStr := strings.TrimSpace(record[2])

		level, err := strconv.Atoi(levelStr)
		if err != nil {
			continue
		}

		classList := strings.Split(classStr, ",")
		for i := range classList {
			classList[i] = strings.ToLower(strings.TrimSpace(classList[i]))
		}

		spell := Spell{
			Name:  name,
			Level: level,
			Class: classList,
		}

		for _, c := range classList {
			ClassSpellList[c] = append(ClassSpellList[c], spell)
		}
	}

	return nil
}

func FindSpellByName(spellName string) *Spell {
	nameLower := strings.ToLower(strings.TrimSpace(spellName))
	for _, spells := range ClassSpellList {
		for i := range spells {
			if spells[i].Name == nameLower {
				return &spells[i]
			}
		}
	}
	return nil
}

func ClassHasSpell(class Class, spellName string) bool {
	classKey := strings.ToLower(strings.TrimSpace(string(class)))
	classSpells, ok := ClassSpellList[classKey]
	if !ok {
		return false
	}

	nameLower := strings.ToLower(strings.TrimSpace(spellName))
	for i := range classSpells {
		if classSpells[i].Name == nameLower {
			return true
		}
	}
	return false
}