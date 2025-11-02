package infrastructure

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"starter_pack/domain"
)

type LearnSpellService struct {
	Repo      domain.CharacterRepository
	SpellRepo domain.SpellRepository
}

type PrepareSpellService struct {
	Repo      domain.CharacterRepository
	SpellRepo domain.SpellRepository
}

type SpellRepository struct {
	classSpells map[string][]domain.Spell
}

func NewSpellRepository() *SpellRepository {
	return &SpellRepository{
		classSpells: make(map[string][]domain.Spell),
	}
}

func (r *SpellRepository) LoadFromCSV(path string) error {
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

		name := strings.TrimSpace(record[0])
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

		spell := domain.Spell{
			Name:  name,
			Level: level,
			Class: classList,
		}

		for _, c := range classList {
			r.classSpells[c] = append(r.classSpells[c], spell)
		}
	}

	return nil
}

func (r *SpellRepository) GetSpellsForClass(class domain.Class) []domain.Spell {
	classKey := strings.ToLower(strings.TrimSpace(string(class)))
	return r.classSpells[classKey]
}

func (r *SpellRepository) FindSpellByName(name string) *domain.Spell {
	nameLower := strings.ToLower(strings.TrimSpace(name))
	for _, spells := range r.classSpells {
		for i := range spells {
			if spells[i].HasName(nameLower) {
				return &spells[i]
			}
		}
	}
	return nil
}

func (r *SpellRepository) ClassHasSpell(class domain.Class, spellName string) bool {
	classKey := strings.ToLower(strings.TrimSpace(string(class)))
	spells, ok := r.classSpells[classKey]
	if !ok {
		return false
	}

	for _, s := range spells {
		if s.HasName(spellName) {
			return true
		}
	}

	return false
}

