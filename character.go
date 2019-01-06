// Package sotdlgen implements a character generator for the SotDL RPG.
package sotdlgen

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"

	logging "github.com/op/go-logging"
)

// Declare logger.
var (
	log       = logging.MustGetLogger("sotdl")
	logLevels = map[string]logging.Level{
		"INFO":    logging.INFO,
		"ERROR":   logging.ERROR,
		"WARNING": logging.WARNING,
	}
)

// Declare various character data lists.
var (
	ancestries = []string{
		"Human", "Dwarf", "Goblin", "Orc", "Changeling", "Clockwork",
	}
	novicePaths = []string{
		"Priest", "Magician", "Warrior", "Rogue",
	}
	expertPaths = []string{
		"Artificer", "Assassin", "Berserker", "Cleric", "Druid", "Fighter",
		"Oracle", "Paladin", "Ranger", "Scout", "Sorcerer", "Spellbinder", "Thief",
		"Warlock", "Witch", "Wizard",
	}
	masterPaths = []string{
		"Abjurer", "Acrobat", "Aeromancer", "Apocalyptist", "Arcanist", "Astromancer",
		"Avenger", "Bard", "Beastmaster", "Blade", "Brute", "Cavalier", "Champion",
		"Chaplain", "Chronomancer", "Conjurer", "Conqueror", "Death Dealer", "Defender",
		"Dervish", "Destroyer", "Diplomat", "Diviner", "Dreadnaught", "Duelist",
		"Enchantment", "Engineer", "Executioner", "Exorcist", "Explorer", "Geomancer",
		"Gladiator", "Gunslinger", "Healer", "Hexer", "Hydromancer", "Illusionist",
		"Infiltrator", "Inquisitor", "Jack-of-all-Trades", "Mage Knight", "Magus",
		"Marauder", "Miracle Worker", "Myrmidon", "Necromancer", "Poisoner", "Pyromancer",
		"Runesmith", "Savant", "Sentinel", "Shapeshifter", "Sharpshooter", "Stormbringer",
		"Technomancer", "Templar", "Tenebrist", "Thaumaturge", "Theurge", "Transmuter",
		"Traveler", "Weapon Master", "Woodwose", "Zealot",
	}
	genders   = []string{"Male", "Female"}
	languages = []string{
		"Common Tongue", "Dark Speech", "Dwarfish", "Elvish", "High Archaic", "Trollish",
		"Secret Language",
	}
)

// Character represents the primary features of the character.
type Character struct {
	Name        string     `json:"name"`
	Gender      string     `json:"gender"`
	Ancestry    string     `json:"ancestry"`
	LangAndProf []string   `json:"languages_and_professions"`
	NovicePath  string     `json:"novice_path"`
	ExpertPath  string     `json:"expert_path"`
	MasterPath  string     `json:"master_path"`
	Talents     []string   `json:"talents"`
	Level       int        `json:"level"`
	Attributes  Attributes `json:"attributes"`
	Seed        string     `json:"seed"`
	//Background  string     `json:"background"`
	//Description string     `json:"description"`
	//Magic       []Spell    `json:"magic"`
	//Weapons     []Weapon   `json:"weapons"`
	//Armor       []Armor    `json:"armor"`
	//Equipment   []string   `json:"equipment"`
}

// Attributes represents character statistics.
type Attributes struct {
	Strength              int    `json:"strength"`
	Agility               int    `json:"agility"`
	Intellect             int    `json:"intellect"`
	Will                  int    `json:"will"`
	Speed                 int    `json:"speed"`
	Power                 int    `json:"power"`
	Health                int    `json:"health"`
	Size                  string `json:"size"`
	Insanity              int    `json:"insanity"`
	Corruption            int    `json:"corruption"`
	Defense               int    `json:"defense"`
	Perception            int    `json:"perception"`
	HealingRate           int    `json:"healing_rate"`
	healthMod             int
	defenseMod            int
	perceptionMod         int
	healingRateMultiplier float64
}

// Spell represents properties of a given spell.
type Spell struct {
	// TODO: Implement spells.
	Name         string `json:"name"`
	Type         string `json:"type"`
	Rank         int    `json:"rank"`
	Target       string `json:"target"`
	Area         string `json:"area"`
	Duration     string `json:"duration"`
	Triggered    bool   `json:"triggered"`
	Sacrifice    bool   `json:"sacrifice"`
	Permanence   bool   `json:"permanence"`
	AttackRoll20 string `json:"attack_20+"`
	Description  string `json:"description"`
}

// Weapon represents properties of a given weapon.
type Weapon struct {
	// TODO: Implement weapons.
	Name         string `json:"name"`
	Type         string `json:"type"`
	Hands        string `json:"hands"`
	Cumbersome   bool   `json:"cumbersome"`
	Finesse      bool   `json:"finesse"`
	DefenseBonus int    `json:"defense_bonus"`
	Misfire      bool   `json:"misfire"`
	Range        string `json:"range"`
	Reach        int    `json:"reach"`
	Reload       bool   `json:"reload"`
	Size         int    `json:"size"`
	Uses         string `json:"uses"`
	Thrown       bool   `json:"thrown"`
	Damage       Die    `json:"damage"`
}

// Armor represents properties of a given suit of armor.
type Armor struct {
	// TODO: Implement armor.
	Name         string `json:"name"`
	Type         string `json:"type"`
	DefenseBonus int    `json:"defense"`
}

// Randomly increments n character attributes by 1.
func (c *Character) incrRandomAttrs(n int) {
	for i := 0; i < n; i++ {
		switch randomInt(0, 4) {
		case 0:
			c.Attributes.Strength++
		case 1:
			c.Attributes.Agility++
		case 2:
			c.Attributes.Intellect++
		case 3:
			c.Attributes.Will++
		}
	}
}

// Sets the random seed from a hex hash string.
func (c *Character) setCharSeed(charSeed string) {
	var err error
	c.Seed, err = setSeed(charSeed)
	if err != nil {
		log.Error("Failed to set character hash:", err)
	}
}

// Pick a random level in [0..10] if non supplied.
func (c *Character) setLevel(level string) {
	if level != "" {
		c.Level, _ = strconv.Atoi(level)
	} else {
		c.Level = randomInt(0, 10)
	}
}

func (c *Character) calcHealingRate() {
	hr := float64(c.Attributes.Health) * c.Attributes.healingRateMultiplier
	c.Attributes.HealingRate = int(math.Floor(hr))
}

func (c *Character) increaseAttributes(i int, path string) {
	if i == 0 && path == "Human" {
		c.incrRandomAttrs(1)
	}
	switch i {
	case 1, 3:
		c.incrRandomAttrs(2)
	case 7:
		c.incrRandomAttrs(3)
	}
}

func (c *Character) calcDerived() {
	c.Attributes.Perception = c.Attributes.Intellect + c.Attributes.perceptionMod
	c.Attributes.Defense = c.Attributes.Agility + c.Attributes.defenseMod
	c.Attributes.Health = c.Attributes.Strength + c.Attributes.healthMod
}

//
func (c *Character) setPath(path string) {
	// Set the path.
	if c.Ancestry == "" {
		if path == "" {
			path = randomChoice(ancestries)
		}
		c.Ancestry = path
	} else if c.NovicePath == "" {
		if path == "" {
			path = randomChoice(novicePaths)
		}
		c.NovicePath = path
	} else if c.ExpertPath == "" {
		if path == "" {
			path = randomChoice(expertPaths)
		}
		c.ExpertPath = path
	} else if c.MasterPath == "" {
		if path == "" {
			path = randomChoice(masterPaths)
		}
		c.MasterPath = path
	} else {
		return
	}
	// Add attributes, etc.
	var keys []int
	for k := range db.Paths[path] {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, i := range keys {
		lvl := db.Paths[path][i]
		if i <= c.Level {
			// Attributes
			c.Attributes.Strength += lvl.Strength
			c.Attributes.Agility += lvl.Agility
			c.Attributes.Intellect += lvl.Intellect
			c.Attributes.Will += lvl.Will

			// Characteristics
			c.Attributes.perceptionMod += lvl.PerceptionMod
			c.Attributes.defenseMod += lvl.DefenseMod
			c.Attributes.healthMod += lvl.HealthMod

			c.Attributes.Speed += lvl.Speed
			c.Attributes.Power += lvl.Power

			if lvl.Size != "" {
				c.Attributes.Size = lvl.Size
			}

			c.Attributes.Insanity += lvl.Insanity
			c.Attributes.Corruption += lvl.Corruption

			if lvl.HealingRate != 0.0 {
				c.Attributes.healingRateMultiplier = lvl.HealingRate
			}

			// Talents
			c.Talents = append(c.Talents, lvl.Talents...)

			// Languages and Professions
			if i == 0 {
				c.LangAndProf = append(c.LangAndProf,
					"Two professions of your choice; you may trade one for a language.")
			}
			c.LangAndProf = append(c.LangAndProf, lvl.LangAndProf...)
		}
		// Attribute increases
		c.increaseAttributes(i, path)
	}
	// Recalc
	c.calcDerived()
	c.calcHealingRate()
}

// Randomly sample from name db.
func (c *Character) setName(name string) {
	if name != "" {
		c.Name = name
	} else {
		firstNames := []string{}
		surnames := []string{}
		ethnicities := []string{}
		for _, nl := range db.Names {
			if nl.Ancestry == c.Ancestry && !arrayContains(ethnicities, nl.Ethnicity) {
				ethnicities = append(ethnicities, nl.Ethnicity)
			}
		}
		ethnicity := ""
		if len(ethnicities) > 0 {
			ethnicity = randomChoice(ethnicities)
		} else {
			ethnicity = db.Names[randomInt(0, len(db.Names))].Ethnicity
		}
		for _, nl := range db.Names {
			if nl.Ethnicity == ethnicity {
				switch nl.Type {
				case c.Gender:
					firstNames = append(firstNames, nl.Names...)
				case "Surname":
					surnames = append(surnames, nl.Names...)
				}
			}
		}
		firstName := ""
		surname := ""
		if len(firstNames) > 0 {
			firstName = randomChoice(firstNames)
		}
		if len(surnames) > 0 {
			surname = randomChoice(surnames)
		}
		c.Name = firstName + " " + surname
	}
}

// Randomly sample from gender list.
func (c *Character) setGender(gender string) {
	if gender != "" {
		c.Gender = gender
	} else {
		c.Gender = randomChoice(genders)
	}
}

// TODO: Additional character data functions.
func (c *Character) setMagic()                         {}
func (c *Character) setWeapons()                       {}
func (c *Character) setArmor()                         {}
func (c *Character) setEquipment()                     {}
func (c *Character) setDescription(description string) {}
func (c *Character) setBackground(background string)   {}
func (c *Character) setProfessions(professions string) {}
func (c *Character) setLanguages(languages string)     {}

// Print writes tab-delimited character details to STDOUT.
func (c Character) Print() {
	fmt.Println("Name\t" + c.Name)
	fmt.Println("Gender\t" + c.Gender)
	fmt.Println("Level\t", c.Level)
	fmt.Println("Character Seed\t", c.Seed)
	fmt.Println()
}

// ToJSON writes JSON character details to STDOUT.
func (c Character) ToJSON(pretty bool) string {
	var j []byte
	if pretty {
		j, _ = json.MarshalIndent(c, "  ", "  ")
	} else {
		j, _ = json.Marshal(c)
	}
	fmt.Println(string(j))
	//err := ioutil.WriteFile(fn, j, 0644)
	//if err != nil {
	//	panic(err)
	//}
	return string(j)
}

// Opts contains user input optionsr; used in CLI implementations.
type Opts struct {
	Age         string
	Ancestry    string `docopt:"--ancestry"`
	Background  string
	Description string
	ExpertPath  string `docopt:"--expert-path"`
	Gender      string `docopt:"--gender"`
	Languages   string
	Level       string `docopt:"--level"`
	LogLevel    string `docopt:"--log-level"`
	MasterPath  string `docopt:"--master-path"`
	Name        string `docopt:"--name"`
	NovicePath  string `docopt:"--novice-path"`
	Professions string
	Seed        string `docopt:"--seed"`
	DataFile    string `docopt:"--data-file"`
}

// NewCharacter generates a SotDL character given a set of user options.
func NewCharacter(opts Opts) (c Character, err error) {

	logging.SetLevel(logLevels[opts.LogLevel], "")

	// Load the character db if empty.
	if len(db.Paths) == 0 {
		log.Info("Loading Character DB.")
		db, err = NewCharDB(opts.DataFile, false)
		if err != nil {
			return c, err
		}
	}

	// Initialize character and set random seed from hash
	c.setCharSeed(opts.Seed)

	// Generate base characteristics
	log.Info("Generating attributes and characteristics.")
	c.setGender(opts.Gender)
	c.setName(opts.Name)
	c.setLevel(opts.Level)
	c.setPath(opts.Ancestry)
	if c.Level > 0 {
		c.setPath(opts.NovicePath)
	}
	if c.Level > 2 {
		c.setPath(opts.ExpertPath)
	}
	if c.Level > 6 {
		c.setPath(opts.MasterPath)
	}

	// Generate stuff
	//c.setMagic()
	//c.setWeapons()
	//c.setArmor()
	//c.setEquipment()

	// Generate fluff
	//c.setDescription(opts.Description)
	//c.setBackground(opts.Background)
	//c.setProfessions(opts.Professions)
	//c.setLanguages(opts.Languages)

	return c, nil
}
