// Unit tests for name generator.
package sotdlgen

import (
	"testing"
)

func TestBuildNamesDB(t *testing.T) {
	db := buildNamesDB()
	if !arrayContains(ancestries, db[0].Ancestry) {
		t.Errorf("Cannot build names database.")
	}
}
