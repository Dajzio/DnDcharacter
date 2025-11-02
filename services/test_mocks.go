package services

import (
	"context"
	"errors"
	"starter_pack/domain"
)

type MockCharacterRepo struct {
	Characters map[string]*domain.Character
	SaveErr    error
}

func (m *MockCharacterRepo) Save(ctx context.Context, c *domain.Character) error {
	if m.SaveErr != nil {
		return m.SaveErr
	}
	m.Characters[c.Name] = c
	return nil
}
func (m *MockCharacterRepo) GetByID(ctx context.Context, id string) (*domain.Character, error) {
	for _, c := range m.Characters {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, errors.New("character not found")
}
func (m *MockCharacterRepo) GetByName(ctx context.Context, name string) (*domain.Character, error) {
	c, ok := m.Characters[name]
	if !ok {
		return nil, errors.New("character not found")
	}
	return c, nil
}
func (m *MockCharacterRepo) List(ctx context.Context) ([]*domain.Character, error) {
	var list []*domain.Character
	for _, c := range m.Characters {
		list = append(list, c)
	}
	return list, nil
}
func (m *MockCharacterRepo) Delete(ctx context.Context, name string) error {
	if _, ok := m.Characters[name]; !ok {
		return errors.New("character not found")
	}
	delete(m.Characters, name)
	return nil
}

type MockSpellRepo struct {
	Spells map[string]domain.Spell
}

func (m *MockSpellRepo) LoadFromCSV(path string) error { return nil }
func (m *MockSpellRepo) GetSpellsForClass(class domain.Class) []domain.Spell {
	var list []domain.Spell
	for _, s := range m.Spells {
		for _, c := range s.Class {
			if c == string(class) {
				list = append(list, s)
			}
		}
	}
	return list
}
func (m *MockSpellRepo) FindSpellByName(name string) *domain.Spell {
	if s, ok := m.Spells[name]; ok {
		return &s
	}
	return nil
}
func (m *MockSpellRepo) ClassHasSpell(class domain.Class, spellName string) bool {
	s := m.FindSpellByName(spellName)
	if s == nil {
		return false
	}
	for _, c := range s.Class {
		if c == string(class) {
			return true
		}
	}
	return false
}
