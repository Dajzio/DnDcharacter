package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"starter_pack/domain"
	"starter_pack/infrastructure"
	"starter_pack/services"
	"strings"
)

func usage() {
	fmt.Printf(`Usage:
  %s create -name CHARACTER_NAME -race RACE -class CLASS -level N -str N -dex N -con N -int N -wis N -cha N
  %s view -name CHARACTER_NAME
  %s list
  %s delete -name CHARACTER_NAME
  %s equip -name CHARACTER_NAME -weapon WEAPON_NAME -slot SLOT
  %s equip -name CHARACTER_NAME -armor ARMOR_NAME
  %s equip -name CHARACTER_NAME -shield SHIELD_NAME
  %s learn-spell -name CHARACTER_NAME -spell SPELL_NAME
  %s prepare-spell -name CHARACTER_NAME -spell SPELL_NAME 
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	charRepo := infrastructure.NewFileCharacterRepo("characters.json")
	spellRepo := infrastructure.NewSpellRepository()
	ctx := context.Background()

	if err := spellRepo.LoadFromCSV("5e-SRD-Spells.csv"); err != nil {
		fmt.Println("Failed to load spells:", err)
		os.Exit(1)
	}

	switch cmd {
	case "create":
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)
		name := createCmd.String("name", "", "character name (required)")
		race := createCmd.String("race", "", "character race")
		class := createCmd.String("class", "", "character class")
		background := createCmd.String("background", "acolyte", "character background")
		level := createCmd.Int("level", 1, "character level")

		str := createCmd.Int("str", 10, "strength")
		dex := createCmd.Int("dex", 10, "dexterity")
		con := createCmd.Int("con", 10, "constitution")
		intel := createCmd.Int("int", 10, "intelligence")
		wis := createCmd.Int("wis", 10, "wisdom")
		cha := createCmd.Int("cha", 10, "charisma")

		if err := createCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Failed to parse flags:", err)
			os.Exit(2)
		}

		if *name == "" {
			fmt.Println("Error: -name is required")
			os.Exit(2)
		}

		skillRepo := domain.NewSkillRepository()

		input := services.CreateCharacterInput{
			Name:       *name,
			Race:       domain.Race(strings.ToLower(*race)),
			Class:      domain.Class(strings.ToLower(*class)),
			Background: *background,
			Level:      *level,
			Str:        *str,
			Dex:        *dex,
			Con:        *con,
			Int:        *intel,
			Wis:        *wis,
			Cha:        *cha,
			Skills:     skillRepo.GetDefaultSkills(*class, *background),
		}
		createService := &services.CreateCharacterService{Repo: charRepo}
		c, err := createService.Execute(ctx, input)
		if err != nil {
			fmt.Println("Error creating character:", err)
			os.Exit(2)
		}

		fmt.Printf("saved character %s\n", c.Name)

	case "list":
		list, err := charRepo.List(ctx)
		if err != nil {
			fmt.Println("Error listing characters:", err)
			os.Exit(2)
		}
		if len(list) == 0 {
			fmt.Println("No characters created yet.")
			return
		}
		fmt.Println("Characters:")
		for _, c := range list {
			fmt.Printf("- %s (Race: %s, Class: %s, Level: %d)\n",
				c.Name, strings.Title(string(c.Race)), strings.Title(string(c.Class)), c.Level)
		}

	case "view":
		viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
		name := viewCmd.String("name", "", "Character name to view")
		if err := viewCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Failed to parse flags:", err)
			os.Exit(2)
		}

		if *name == "" {
			fmt.Println("Please provide a name using -name flag")
			os.Exit(1)
		}

		viewService := &services.ViewCharacterService{Repo: charRepo}

		if err := viewService.Execute(ctx, *name); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		name := deleteCmd.String("name", "", "Character name to delete")

		if err := deleteCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Failed to parse flags:", err)
			os.Exit(2)
		}

		if *name == "" {
			fmt.Println("Error: -name is required")
			os.Exit(2)
		}

		deleteService := &services.DeleteCharacterService{Repo: charRepo}
		if err := deleteService.Execute(ctx, *name); err != nil {
			fmt.Println("Error deleting character:", err)
			os.Exit(2)
		}

		fmt.Printf("deleted %s\n", *name)

	case "equip":
		equipCmd := flag.NewFlagSet("equip", flag.ExitOnError)
		name := equipCmd.String("name", "", "Character name")
		weapon := equipCmd.String("weapon", "", "Weapon name")
		armor := equipCmd.String("armor", "", "Armor name")
		shield := equipCmd.String("shield", "", "Shield name")
		slot := equipCmd.String("slot", "", "Slot (main hand/off hand)")

		if err := equipCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Failed to parse flags:", err)
			os.Exit(2)
		}

		equipService := &services.EquipItemService{Repo: charRepo}
		var output string
		var err error

		switch {
		case *weapon != "":
			output, err = equipService.Execute(ctx, *name, "weapon", *weapon, *slot)
		case *armor != "":
			output, err = equipService.Execute(ctx, *name, "armor", *armor, "")
		case *shield != "":
			output, err = equipService.Execute(ctx, *name, "shield", *shield, "")
		default:
			fmt.Println("Please specify an item to equip (weapon, armor, or shield).")
			os.Exit(1)
		}

		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println(output)

	case "learn-spell":
		learnCmd := flag.NewFlagSet("learn-spell", flag.ExitOnError)
		name := learnCmd.String("name", "", "Character name")
		spell := learnCmd.String("spell", "", "Spell name")

		if err := learnCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Failed to parse flags:", err)
			os.Exit(2)
		}

		if *name == "" || *spell == "" {
			fmt.Println("Error: -name and -spell are required")
			os.Exit(1)
		}

		learnService := &services.LearnSpellService{
			Repo:     charRepo,
			SpellRepo:     spellRepo,
		}
		if output, err := learnService.Execute(ctx, *name, *spell); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		} else {
			fmt.Println(output)
		}

	case "prepare-spell":
		prepareCmd := flag.NewFlagSet("prepare-spell", flag.ExitOnError)
		name := prepareCmd.String("name", "", "Character name")
		spell := prepareCmd.String("spell", "", "Spell name")

		if err := prepareCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Failed to parse flags:", err)
			os.Exit(2)
		}

		if *name == "" || *spell == "" {
			fmt.Println("Error: -name and -spell are required")
			os.Exit(1)
		}

		prepareService := &services.PrepareSpellService{
			Repo:      charRepo,
			SpellRepo: spellRepo,
		}
		if output, err := prepareService.Execute(ctx, *name, *spell); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		} else {
			fmt.Println(output)
		}

	case "enrich":
		services.EnrichData()

	default:
		usage()
		os.Exit(2)
	}
}
