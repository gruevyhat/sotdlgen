// Implements character generation logic.
package sotdlgen

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("sotdl")

var logLevels = map[string]logging.Level{
	"INFO":    logging.INFO,
	"ERROR":   logging.ERROR,
	"WARNING": logging.WARNING,
}

const (
	startErr = "\033[31m"
	endErr   = "\033[0m"
)

var ancestries = []string{
	"Human", "Dwarf", "Goblin", "Orc", "Changeling", "Clockwork",
}

var novicePaths = []string{
	"Priest", "Magician", "Warrior", "Rogue",
}
var expertPaths = []string{
	"Artificer", "Assassin", "Berserker", "Cleric", "Druid", "Fighter",
	"Oracle", "Paladin", "Ranger", "Scout", "Sorcerer", "Spellbinder", "Thief",
	"Warlock", "Witch", "Wizard",
}

var masterPaths = []string{
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

var genders = []string{"Male", "Female", "Other"}

var languages = []string{"Common Tongue", "Dark Speech", "Dwarfish", "Elvish", "High Archaic", "Trollish", "Secret Language"}

type Character struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Gender      string      `json:"gender"`
	Ancestry    string      `json:"ancestry"`
	Background  string      `json:"background"`
	LangAndProf []string    `json:"languages_and_professions"`
	NovicePath  string      `json:"novice_path"`
	ExpertPath  string      `json:"expert_path"`
	MasterPath  string      `json:"master_path"`
	Talents     []string    `json:"talents"`
	Magic       []Spell     `json:"magic"`
	Weapons     []Weapon    `json:"weapons"`
	Armor       []Armor     `json:"armor"`
	Equipment   []Equipment `json:"equipment"`
	Level       int         `json:"level"`
	Attributes  Attributes  `json:"attributes"`
	Hash        string
}

type Profession string
type Language string
type Path string
type Equipment string

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

type Armor struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	DefenseBonus int    `json:"defense"`
}

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

type Die struct {
	code int
	pips int
}

func (d Die) toStr() string {
	var dieStr string
	if d.pips > 0 {
		dieStr = strconv.Itoa(d.code) + "d6+" + strconv.Itoa(d.pips)
	} else {
		dieStr = strconv.Itoa(d.code) + "d6"
	}
	return dieStr
}

func (c *Character) setCharHash(charHash string) {
	if charHash != "" {
		c.Hash = charHash
		log.Info("NEW SEED:", c.Hash)
	} else {
		var err error
		c.Hash, err = setSeed(charHash)
		if err != nil {
			log.Error(err)
		} else {
			log.Info("OLD SEED:", c.Hash)
		}
	}
}

func (c *Character) getLevel(level string) {
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

func (c *Character) getPath(path string) {
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

func (c *Character) getMagic()     {}
func (c *Character) getWeapons()   {}
func (c *Character) getArmor()     {}
func (c *Character) getEquipment() {}

func (c *Character) getName(name string) {
	if name != "" {
		c.Name = name
	}
}

func (c *Character) getGender(gender string) {
	if gender != "" {
		c.Gender = gender
	} else {
		c.Gender = randomChoice(genders)
	}
}

func (c *Character) getDescription(description string) {}
func (c *Character) getBackground(background string)   {}
func (c *Character) getProfessions(professions string) {}
func (c *Character) getLanguages(languages string)     {}

func (c Character) Print() {
	fmt.Println("Name\t" + c.Name)
	fmt.Println("Gender\t" + c.Gender)
	fmt.Println("Level\t", c.Level)
	fmt.Println("Character Hash\t", c.Hash)
	fmt.Println()
}

func (c Character) ToJSON() {
	j, _ := json.MarshalIndent(c, "  ", "  ")
	fmt.Println(string(j))
	//err := ioutil.WriteFile(fn, j, 0644)
	//if err != nil {
	//	panic(err)
	//}
}

func NewCharacter(opts map[string]string) Character {

	logging.SetLevel(logLevels[opts["log-level"]], "")

	c := Character{}

	// Generate base characteristics
	c.setCharHash(opts["seed"])
	c.getLevel(opts["level"])
	c.getPath(opts["ancestry"])
	if c.Level > 0 {
		c.getPath(opts["novice-path"])
	}
	if c.Level > 2 {
		c.getPath(opts["expert-path"])
	}
	if c.Level > 6 {
		c.getPath(opts["master-path"])
	}

	// Generate stuff
	c.getMagic()
	c.getWeapons()
	c.getArmor()
	c.getEquipment()

	// Generate fluff
	c.getName(opts["name"])
	c.getGender(opts["gender"])
	c.getDescription(opts["description"])
	c.getBackground(opts["background"])
	c.getProfessions(opts["professions"])
	c.getLanguages(opts["languages"])

	c.ToJSON()

	return c
}
