package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type Spell struct {
	Name   string `json:"name"`
	School string `json:"school"`
	Range  string `json:"range"`
}

type Weapon struct {
	Name      string `json:"name"`
	Category  string `json:"category"`
	Range     string `json:"range"`
	TwoHanded bool   `json:"two_handed"`
}

type Armor struct {
	Name          string `json:"name"`
	ArmorClass    int    `json:"armor_class"`
	DexterityBonus bool   `json:"dex_bonus"`
}

var requestLimiter = time.Tick(100 * time.Millisecond) 

func EnrichData() error {
	fmt.Println("Fetching D&D 5e data from API...")

	var wg sync.WaitGroup
	var spells []Spell
	var weapons []Weapon
	var armors []Armor
	var spellErr, weaponErr, armorErr error

	wg.Add(3)
	go func() {
		defer wg.Done()
		spells, spellErr = fetchSpells()
	}()
	go func() {
		defer wg.Done()
		weapons, weaponErr = fetchWeapons()
	}()
	go func() {
		defer wg.Done()
		armors, armorErr = fetchArmors()
	}()
	wg.Wait()

	if spellErr != nil {
		return fmt.Errorf("failed to fetch spells: %w", spellErr)
	}
	if weaponErr != nil {
		return fmt.Errorf("failed to fetch weapons: %w", weaponErr)
	}
	if armorErr != nil {
		return fmt.Errorf("failed to fetch armors: %w", armorErr)
	}

	if err := saveJSON("spells_data.json", spells); err != nil {
		return fmt.Errorf("saving spells failed: %w", err)
	}

	equipmentData := map[string]interface{}{
		"weapons": weapons,
		"armors":  armors,
	}
	if err := saveJSON("eq_data.json", equipmentData); err != nil {
		return fmt.Errorf("saving equipment failed: %w", err)
	}

	fmt.Printf("Data successfully enriched and saved: spells (%d), weapons (%d), armors (%d)\n",
		len(spells), len(weapons), len(armors))
	return nil
}

func fetchSpells() ([]Spell, error) {
	baseURL := "https://www.dnd5eapi.co/api/2014/spells"
	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var base struct {
		Results []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&base); err != nil {
		return nil, err
	}

	ch := make(chan Spell)
	var wg sync.WaitGroup
	for _, s := range base.Results {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			sp, err := fetchSpellDetail(url)
			if err == nil {
				ch <- sp
			}
		}(s.URL)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	spells := []Spell{}
	for sp := range ch {
		spells = append(spells, sp)
	}

	return spells, nil
}

func fetchSpellDetail(path string) (Spell, error) {
	<-requestLimiter
	data, err := fetchDetail(path)
	if err != nil {
		return Spell{}, err
	}

	var d struct {
		Name   string `json:"name"`
		School struct {
			Name string `json:"name"`
		} `json:"school"`
		Range string `json:"range"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return Spell{}, err
	}
	return Spell{Name: d.Name, School: d.School.Name, Range: d.Range}, nil
}

func fetchWeapons() ([]Weapon, error) {
	baseURL := "https://www.dnd5eapi.co/api/2014/equipment-categories/weapon"
	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var base struct {
		Equipment []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"equipment"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&base); err != nil {
		return nil, err
	}

	ch := make(chan Weapon)
	var wg sync.WaitGroup
	for _, w := range base.Equipment {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			wp, err := fetchWeaponDetail(url)
			if err == nil {
				ch <- wp
			}
		}(w.URL)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	weapons := []Weapon{}
	for wp := range ch {
		weapons = append(weapons, wp)
	}

	return weapons, nil
}

func fetchWeaponDetail(path string) (Weapon, error) {
	<-requestLimiter
	data, err := fetchDetail(path)
	if err != nil {
		return Weapon{}, err
	}

	var d struct {
		Name          string      `json:"name"`
		CategoryRange string      `json:"weapon_category"`
		TwoHanded     bool        `json:"two_handed"`
		Range         interface{} `json:"range"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return Weapon{}, err
	}

	rangeText := ""
	switch v := d.Range.(type) {
	case float64:
		rangeText = fmt.Sprintf("%d ft", int(v))
	case map[string]interface{}:
		if normal, ok := v["normal"].(float64); ok {
			rangeText = fmt.Sprintf("%d ft", int(normal))
		}
	}

	return Weapon{
		Name:      d.Name,
		Category:  d.CategoryRange,
		Range:     rangeText,
		TwoHanded: d.TwoHanded,
	}, nil
}

func fetchArmors() ([]Armor, error) {
	baseURL := "https://www.dnd5eapi.co/api/2014/equipment-categories/armor"
	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var base struct {
		Equipment []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"equipment"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&base); err != nil {
		return nil, err
	}

	ch := make(chan Armor)
	var wg sync.WaitGroup
	for _, a := range base.Equipment {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			ar, err := fetchArmorDetail(url)
			if err == nil {
				ch <- ar
			}
		}(a.URL)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	armors := []Armor{}
	for ar := range ch {
		armors = append(armors, ar)
	}

	return armors, nil
}

func fetchArmorDetail(path string) (Armor, error) {
	<-requestLimiter
	data, err := fetchDetail(path)
	if err != nil {
		return Armor{}, err
	}

	var d struct {
		Name       string `json:"name"`
		ArmorClass struct {
			Base     int  `json:"base"`
			DexBonus bool `json:"dex_bonus"`
		} `json:"armor_class"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return Armor{}, err
	}

	return Armor{
		Name:           d.Name,
		ArmorClass:     d.ArmorClass.Base,
		DexterityBonus: d.ArmorClass.DexBonus,
	}, nil
}

func fetchDetail(path string) ([]byte, error) {
	fullURL := "https://www.dnd5eapi.co" + path
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func saveJSON(filename string, data interface{}) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, bytes, 0644)
}
