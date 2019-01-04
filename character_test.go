package sotdlgen

import (
	"strings"
	"testing"
)

func TestNewCharacter(t *testing.T) {
	opts := []Opts{
		Opts{
			Name:     "Xev",
			LogLevel: "INFO",
		},
		Opts{
			DataFile: "./assets/Shadow_of_the_Demon_Lord.pdf",
			LogLevel: "INFO",
		},
		Opts{
			Name:       "Borkenhekenaken",
			Gender:     "Male",
			Ancestry:   "Goblin",
			NovicePath: "Magician",
			ExpertPath: "Wizard",
			Seed:       "1575d911f49e59ee",
			Level:      "3",
			LogLevel:   "INFO",
		},
	}
	for _, o := range opts {
		c, _ := NewCharacter(o)
		if c.Name != o.Name {
			t.Errorf("Incorrect name. Expected '%s'. Found '%s'.", c.Name, o.Name)
		}
		if c.Hash == "" {
			t.Error("Incorrect Hash. No value assigned")
		}
		if !arrayContains(genders, c.Gender) {
			g := strings.Join(genders, ", ")
			t.Errorf("Incorrect gender. Expected '%s' in '%s'.", c.Gender, g)
		}
		if c.Level < 0 {
			t.Errorf("Incorrect Level. '%d' is less than zero.", c.Level)
		}
	}
}
