package sotdlgen

import (
	"strings"
	"testing"
)

type opt map[string]string

func TestNewCharacter(t *testing.T) {
	opts := []opt{
		opt{
			"name":        "Borkenhekenaken",
			"gender":      "Male",
			"ancestry":    "Goblin",
			"novice-path": "Magician",
			"expert-path": "Wizard",
			"seed":        "1575d911f49e59ee",
			"level":       "3",
		},
		opt{
			"name":     "Xev",
			"gender":   "",
			"seed":     "",
			"ancestry": "",
		},
	}
	for _, o := range opts {
		c := NewCharacter(o)
		if c.Name != o["name"] {
			t.Errorf("Incorrect name. Expected '%s'. Found '%s'.", c.Name, o["name"])
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
