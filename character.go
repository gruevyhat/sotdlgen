// Implements character generation logic.
package sotdlgen

import (
	"fmt"
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

type Character struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Gender      string          `json:"gender"`
	Ancestry    string          `json:"ancestry"`
	Background  string          `json:"background"`
	Professions []Profession    `json:"professions"`
	Languages   []Language      `json:"languages"`
	Paths       map[string]Path `json:"paths"`
	Talents     []Talent        `json:"talents"`
	Magic       []Spell         `json:"magic"`
	Weapons     []Weapon        `json:"weapons"`
	Armor       []Armor         `json:"armor"`
	Equipment   []Equipment     `json:"equipment"`
	Level       int             `json:"level"`
	Attributes  Attributes      `json:"attributes"`
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
	Name    string `json:"name"`
	Type    string `json:"type"`
	Defense int    `json:"defense"`
}

type Talent struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
}

type Attributes struct {
	Strength    int `json:"strength"`
	Agility     int `json:"agility"`
	Intellect   int `json:"intellect"`
	Will        int `json:"will"`
	Speed       int `json:"speed"`
	Health      int `json:"health"`
	Size        int `json:"size"`
	Insanity    int `json:"insanity"`
	Corruption  int `json:"corruption"`
	Defense     int `json:"defense"`
	Perception  int `json:"perception"`
	HealingRate int `json:"healing_rate"`
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

func (c *Character) getPaths(paths string) {}

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
}

func NewCharacter(opts map[string]string) Character {

	logging.SetLevel(logLevels[opts["log-level"]], "")

	c := Character{}

	// Generate base characteristics
	c.setCharHash(opts["seed"])
	c.getLevel(opts["level"])
	c.getPaths(opts["ancestry"])
	if c.Level > 0 {
		c.getPaths(opts["novice-path"])
	}
	if c.Level > 2 {
		c.getPaths(opts["expert-path"])
	}
	if c.Level > 6 {
		c.getPaths(opts["master-path"])
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

	c.Print()

	return c
}
