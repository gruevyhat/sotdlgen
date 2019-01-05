// Implements character generation logic.
package sotdlgen

import (
	"encoding/json"
	"fmt"
	"math"
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
	genders   = []string{"Male", "Female", "Other"}
	languages = []string{
		"Common Tongue", "Dark Speech", "Dwarfish", "Elvish", "High Archaic", "Trollish",
		"Secret Language",
	}
)

// Declare primary character data structure.
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
	Hash        string
	//Background  string     `json:"background"`
	//Description string     `json:"description"`
	//Magic       []Spell    `json:"magic"`
	//Weapons     []Weapon   `json:"weapons"`
	//Armor       []Armor    `json:"armor"`
	//Equipment   []string   `json:"equipment"`
}

// Declare character statistics struct.
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

// TODO: Implement spells.
type Spell struct {
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

// TODO: Implement weapons.
type Weapon struct {
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

// TODO: Implement armor.
type Armor struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	DefenseBonus int    `json:"defense"`
}

// Randomly increments n character attributes by 1.
func (c *Character) incrRandomAttrs(n int) {
	for i := 0; i < n; i++ {
		switch randomInt(0, 4) {
		case 0:
			c.Attributes.Strength += 1
		case 1:
			c.Attributes.Agility += 1
		case 2:
			c.Attributes.Intellect += 1
		case 3:
			c.Attributes.Will += 1
		}
	}
}

// Sets the random seed from a hex hash string.
func (c *Character) setCharHash(charHash string) {
	var err error
	c.Hash, err = setSeed(charHash)
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
	c.Attributes.HealingRate = int(math.Floor(float64(c.Attributes.Health) * c.Attributes.healingRateMultiplier))
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
	for i, lvl := range db.Paths[path] {
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
			for _, tal := range lvl.Talents {
				c.Talents = append(c.Talents, tal)
			}
			// Languages and Professions
			for _, lp := range lvl.LangAndProf {
				c.LangAndProf = append(c.LangAndProf, lp)
			}
		}
		// Attribute increases
		c.increaseAttributes(i, path)
	}
	// Recalc
	c.calcHealingRate()
	c.calcDerived()
}

// Randomly sample from name db.
func (c *Character) setName(name string) {
	if name != "" {
		c.Name = name
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

// Write tab-delimited character details to STDOUT.
func (c Character) Print() {
	fmt.Println("Name\t" + c.Name)
	fmt.Println("Gender\t" + c.Gender)
	fmt.Println("Level\t", c.Level)
	fmt.Println("Character Hash\t", c.Hash)
	fmt.Println()
}

// Write JSON character details to STDOUT.
func (c Character) ToJSON() string {
	j, _ := json.MarshalIndent(c, "  ", "  ")
	fmt.Println(string(j))
	//err := ioutil.WriteFile(fn, j, 0644)
	//if err != nil {
	//	panic(err)
	//}
	return string(j)
}

// Character options.
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

// Generates a new character given a set of options.
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
	c.setCharHash(opts.Seed)

	// Generate base characteristics
	log.Info("Generating attributes and characteristics.")
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
	c.setMagic()
	c.setWeapons()
	c.setArmor()
	c.setEquipment()

	// Generate fluff
	log.Info("Generating fluff.")
	c.setName(opts.Name)
	c.setGender(opts.Gender)
	c.setDescription(opts.Description)
	c.setBackground(opts.Background)
	c.setProfessions(opts.Professions)
	c.setLanguages(opts.Languages)

	return c, nil
}
