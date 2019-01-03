// Randomly generate names by ethnicity.
package sotdlgen

import (
	"encoding/json"
)

var namesFile = dataDir + "ik_names.json"

type NameList struct {
	Ancestry  string   `json:"ancestry"`
	Ethnicity string   `json:"ethnicity"`
	Type      string   `json:"type"`
	Names     []string `json:"names"`
}

func buildNamesDB() []NameList {
	var db []NameList
	if err := json.Unmarshal(readJson(namesFile), &db); err != nil {
		log.Error(err)
	}
	return db
}

// TODO: Randomly generate.
