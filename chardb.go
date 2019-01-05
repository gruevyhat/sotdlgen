// Parses the SotDL core rules PDF and builds a character database.

package sotdlgen

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var db = CharDB{}

// Data filenames
var corebookJSON = dataDir + "Shadow_of_the_Demon_Lord.json"
var namesFile = dataDir + "ik_names.json"

type CharDB struct {
	Paths map[string]Levels `json:"paths"`
	Names []NameList        `json:"names"`
}

type Levels map[int]*Level

type Level struct {
	Strength      int      `json:"strength"`
	Agility       int      `json:"agility"`
	Intellect     int      `json:"intellect"`
	Will          int      `json:"will"`
	PerceptionMod int      `json:"perception_mod"`
	DefenseMod    int      `json:"defense_mod"`
	HealthMod     int      `json:"health_mod"`
	HealingRate   float64  `json:"healing_rate"`
	Speed         int      `json:"speed"`
	Power         int      `json:"power"`
	Damage        int      `json:"damage"`
	Insanity      int      `json:"insanity"`
	Corruption    int      `json:"corruption"`
	Size          string   `json:"size"`
	LangAndProf   []string `json:"lang_and_prof"`
	Talents       []string `json:"talents"`
}

var reWhite = regexp.MustCompile(`(?m:\s+)`)

// Trim whitespace
func trim(text string) string {
	text = reWhite.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	return text
}

var ancestryLevelPatterns = map[int]string{
	0: `(?s)\s*Creating An? %s.*?` +
		`Starting Attribute Scores (?P<Attr>.*?)` +
		`Perception (?P<Perc>.*?)\n` +
		`Defense (?P<Def>.*?)\n` +
		`Health (?P<Hlth>.*?)\n` +
		`Healing Rate (?P<HR>.*?)\n` +
		`Size (?P<Sz>.*?), Speed (?P<Spd>.*?), Power (?P<Pwr>.*?)\n` +
		`Damage (?P<Dmg>.*?), Insanity (?P<Ins>.*?), Corruption (?P<Cor>.*?)\n` +
		`(?P<Desc>.*?)\n\n`,
	4: `(?s)Level 4 Expert %s.*?` +
		`Characteristics (?P<Char>.*?)\n` +
		`(?P<Desc>.*?)\n\n`,
}

var novicePathLevelPatterns = map[int]string{
	1: `(?s)\s*Level 1 %s.*?Attributes (?P<Attr>.*?)\nCharacteristics (?P<Char>.*?)\n(?P<Desc>.*?)\n\n`,
	2: `(?s)Level 2 %s.*?Characteristics (?P<Char>.*?)\n(?P<Desc>.*?)\n\n`,
	5: `(?s)\s*Level 5 Expert %s.*?Characteristics (?P<Char>.*?)\n(?P<Desc>.*?)\n\n`,
	8: `(?s)\s*Level 8\s*Master %s.*?Characteristics (?P<Char>.*?)\n(?P<Desc>.*?)\n\n`,
}

var expertPathLevelPatterns = map[int]string{
	3: `(?s)Level 3 %s.*?Attributes (?P<Attr>.*?)\nCharacteristics (?P<Char>.*?)\n(?P<Desc>.*?)\n\n`,
	6: `(?s)Level 6 %s.*?Characteristics (?P<Char>.*?)\n(?P<Desc>.*?)\n\n`,
	9: `(?sm)\s*Level 9\s*Master %s.*?Characteristics (?P<Char>.*?)\n(?P<Desc>.*?)\n\n`,
}

var masterPathLevelPatterns = map[int]string{
	7:  `(?s)Level 7 %s.*?Attributes (?P<Attr>.*?)\nCharacteristics (?P<Char>.*?)\n(?P<Desc>.*?)\n\n`,
	10: `(?s)Level 10 %s.*?Characteristics (?P<Char>.*?)\n(?P<Desc>.*?)\n\n`,
}

type Patterns map[int]*regexp.Regexp

type NameList struct {
	Ancestry  string   `json:"ancestry"`
	Ethnicity string   `json:"ethnicity"`
	Type      string   `json:"type"`
	Names     []string `json:"names"`
}

func (db *CharDB) buildNames() {
	var names []NameList
	if err := json.Unmarshal(readJson(namesFile), &names); err != nil {
		log.Error(err)
	}
	db.Names = names
}

// Compiles patterns to regular expressions.
func compilePatterns(path string, ptns map[int]string) map[int]*regexp.Regexp {
	reMap := Patterns{}
	for key, ptn := range ptns {
		ptn = fmt.Sprintf(ptn, path)
		reMap[key] = regexp.MustCompile(ptn)
	}
	return reMap
}

// Processes a PDF file with pdftotext. Adapted from <https://github.com/plimble/gika>.
func PDFToText(fn string, out io.Writer) error {
	cmd := exec.Command("pdftotext", "-q", fn, "-")
	stderr := bytes.NewBuffer(nil)
	cmd.Stdout = out

	cmd.Start()
	cmdDone := make(chan error, 1)
	go func() {
		cmdDone <- cmd.Wait()
	}()

	select {
	case <-time.After(time.Duration(500000) * time.Millisecond):
		if err := cmd.Process.Kill(); err != nil {
			return errors.New(err.Error())
		}
		<-cmdDone
		return errors.New("Command timed out")
	case err := <-cmdDone:
		if err != nil {
			return errors.New(stderr.String())
		}
	}

	return nil
}

// Create a new SotDL character database.
func NewCharDB(pdfFn string, analyze bool) (db CharDB, err error) {
	// Build db.
	if _, err = os.Stat(corebookJSON); pdfFn != "" || os.IsNotExist(err) {
		// Build db from PDF.
		log.Info("Extracting DB from PDF.")
		ws := &bytes.Buffer{}
		if err = PDFToText(pdfFn, ws); err != nil {
			log.Error("SotDL Core Rules not found.")
			return db, err
		}
		doc := ws.String()
		db.initialize()
		db.extract(doc, ancestries, ancestryLevelPatterns)
		db.extract(doc, novicePaths, novicePathLevelPatterns)
		db.extract(doc, expertPaths, expertPathLevelPatterns)
		db.extract(doc, masterPaths, masterPathLevelPatterns)
		db.buildNames()
		db.save()
		if analyze {
			db.analyze(doc)
		}
	} else {
		// Load an existing db.
		log.Info("Loading DB from JSON.")
		db.load(corebookJSON)
	}
	return db, nil
}

// Build the nested maps.
func (db *CharDB) initialize() {
	db.Paths = make(map[string]Levels)
	for i := 0; i < len(ancestries); i++ {
		db.Paths[ancestries[i]] = Levels{
			0: &Level{}, 4: &Level{},
		}
	}
	for i := 0; i < len(novicePaths); i++ {
		db.Paths[novicePaths[i]] = Levels{
			1: &Level{}, 2: &Level{}, 5: &Level{}, 8: &Level{},
		}
	}
	for i := 0; i < len(expertPaths); i++ {
		db.Paths[expertPaths[i]] = Levels{
			3: &Level{}, 6: &Level{}, 9: &Level{},
		}
	}
	for i := 0; i < len(masterPaths); i++ {
		db.Paths[masterPaths[i]] = Levels{
			7: &Level{}, 10: &Level{},
		}
	}
}

func (db *CharDB) load(fn string) {
	fn = dataDir + path.Base(fn)
	fi, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	b, err := ioutil.ReadAll(fi)
	if err := json.Unmarshal(b, &db); err != nil {
		panic(err)
	}
}

func (db *CharDB) save() {
	j, _ := json.Marshal(db)
	err := ioutil.WriteFile(corebookJSON, j, 0644)
	if err != nil {
		panic(err)
	}
}

func strToInt(s string) (i int) {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

var attributePatterns = map[string]*regexp.Regexp{
	"Strength":   regexp.MustCompile(`Strength (?P<n>\d+)`),
	"Agility":    regexp.MustCompile(`Agility (?P<n>\d+)`),
	"Intellect":  regexp.MustCompile(`Intellect (?P<n>\d+)`),
	"Will":       regexp.MustCompile(`Will (?P<n>\d+)`),
	"Perception": regexp.MustCompile(`Perception\s*(by|\+)\s*(?P<n>\d+)`),
	"Defense":    regexp.MustCompile(`Defense\s*(by|\+)\s*(?P<n>\d+)`),
	"Health":     regexp.MustCompile(`Health\s*(by|\+)\s*(?P<n>\d+)`),
	"Power":      regexp.MustCompile(`Power\s*(by|\+)\s*(?P<n>\d+)`),
}

func (lvl *Level) parsePrimary(text string) {
	for attr, ptn := range attributePatterns {
		m := ptn.FindStringSubmatch(text)
		if len(m) > 1 {
			var n int
			for i, name := range ptn.SubexpNames() {
				if name == "n" {
					n, _ = strconv.Atoi(m[i])
				}
			}
			switch attr {
			case "Strength":
				lvl.Strength += n
			case "Agility":
				lvl.Agility += n
			case "Intellect":
				lvl.Intellect += n
			case "Will":
				lvl.Will += n
			case "Perception":
				lvl.PerceptionMod += n
			case "Defense":
				lvl.DefenseMod += n
			case "Health":
				lvl.HealthMod += n
			case "Power":
				lvl.Power += n
			}
		}
	}
}

var bonusPattern = regexp.MustCompile(`score\s*\+\s*(\d+)`)

func (lvl *Level) parseDerived(text string) {
	mod := 0
	m := bonusPattern.FindStringSubmatch(text)
	if len(m) > 1 {
		mod, _ = strconv.Atoi(m[1])
	}
	if strings.Contains(text, "equals your Strength") {
		lvl.HealthMod += mod
	}
	if strings.Contains(text, "equals your Agility") {
		lvl.DefenseMod += mod
	}
	if strings.Contains(text, "equals your Intellect") {
		lvl.PerceptionMod += mod
	}
}

func (lvl *Level) parseHealingRate(text string) {
	if strings.Contains(text, "one-quarter your Health") {
		lvl.HealingRate = 0.25
	}
}

var langProfPattern = regexp.MustCompile(`Languages and Professions (.*?\.)`)

func (lvl *Level) parseTalents(text string) {
	m := langProfPattern.FindStringSubmatch(text)
	if len(m) > 1 {
		lvl.LangAndProf = append(lvl.LangAndProf, m[1])
		text = strings.Replace(text, m[0], "", 1)
	}
	if text != "" {
		lvl.Talents = append(lvl.Talents, trim(text))
	}
}

func (db *CharDB) extract(doc string, paths []string, pathPatterns map[int]string) {
	for _, path := range paths {
		reMap := compilePatterns(path, pathPatterns)
		for lvl, re := range reMap {
			m := re.FindStringSubmatch(doc)
			for i, name := range re.SubexpNames() {
				text := trim(m[i])
				switch name {
				case "Attr", "Char":
					db.Paths[path][lvl].parsePrimary(text)
				case "Hlth", "Perc", "Def":
					db.Paths[path][lvl].parseDerived(text)
				case "HR":
					db.Paths[path][lvl].parseHealingRate(text)
				case "Ins":
					db.Paths[path][lvl].Insanity += strToInt(text)
				case "Cor":
					db.Paths[path][lvl].Corruption += strToInt(text)
				case "Pwr":
					db.Paths[path][lvl].Power += strToInt(text)
				case "Spd":
					db.Paths[path][lvl].Speed += strToInt(text)
				case "Sz":
					db.Paths[path][lvl].Size = text
				case "Desc":
					db.Paths[path][lvl].parseTalents(text)
				}
				//	fmt.Println(path, "::", name, "::", lvl, "::", text)
				//	fmt.Println(db.Paths[path][lvl])
			}
		}
	}
}

func (c *CharDB) analyze(doc string) {
	for p := range ancestries {
		path := ancestries[p]
		reMap := compilePatterns(path, ancestryLevelPatterns)
		fmt.Printf("%s :: %d :: %q\n", path, 1, reMap[0].FindAllStringSubmatch(doc, -1)[0][1:13])
		fmt.Printf("%s :: %d :: %q\n", path, 2, reMap[4].FindAllStringSubmatch(doc, -1)[0][1:2])
		fmt.Println()
	}
	for p := range novicePaths {
		path := novicePaths[p]
		reMap := compilePatterns(path, novicePathLevelPatterns)
		fmt.Printf("%s :: %d :: %q\n", path, 1, reMap[1].FindAllStringSubmatch(doc, -1)[0][1:4])
		fmt.Printf("%s :: %d :: %q\n", path, 2, reMap[2].FindAllStringSubmatch(doc, -1)[0][1:3])
		fmt.Printf("%s :: %d :: %q\n", path, 5, reMap[5].FindAllStringSubmatch(doc, -1)[0][1:3])
		fmt.Printf("%s :: %d :: %q\n", path, 8, reMap[8].FindAllStringSubmatch(doc, -1)[0][1:3])
		fmt.Println()
	}
	for p := range expertPaths {
		path := expertPaths[p]
		reMap := compilePatterns(path, expertPathLevelPatterns)
		fmt.Printf("%s :: %d :: %q\n", path, 3, reMap[3].FindAllStringSubmatch(doc, -1)[0][1:4])
		fmt.Printf("%s :: %d :: %q\n", path, 6, reMap[6].FindAllStringSubmatch(doc, -1)[0][1:3])
		fmt.Printf("%s :: %d :: %q\n", path, 9, reMap[9].FindAllStringSubmatch(doc, -1)[0][1:3])
		fmt.Println()
	}
	for p := range masterPaths {
		path := masterPaths[p]
		reMap := compilePatterns(path, masterPathLevelPatterns)
		fmt.Printf("%s :: %d :: %q\n", path, 7, reMap[7].FindAllStringSubmatch(doc, -1)[0][1:4])
		fmt.Printf("%s :: %d :: %q\n", path, 10, reMap[10].FindAllStringSubmatch(doc, -1)[0][1:3])
		fmt.Println()
	}

}
