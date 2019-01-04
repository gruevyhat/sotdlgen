package sotdlgen

import (
	"fmt"
	"os"
	"testing"
)

func TestNewCharDB(t *testing.T) {

	pdfFn := "Shadow_of_the_Demon_Lord.pdf"
	jsonFn := "Shadow_of_the_Demon_Lord.json"

	err := os.Remove(dataDir + jsonFn)
	if err != nil {
		fmt.Println("No JSON file found.")
	}

	db, err = NewCharDB("./assets/"+pdfFn, false)
	if err != nil {
		t.Error("Failed to create DB from pdfFn.")
	}

	if _, err = os.Stat(dataDir + jsonFn); os.IsNotExist(err) {
		t.Error("Failed to create JSON file.")
	}

	db, err = NewCharDB("", false)
	if err != nil {
		t.Error("Failed to create DB from empty pdfFn.")
	}

	l := len(db.Paths)
	if l != 90 {
		t.Errorf("DB incorrect size. Expected %d, got %d.", 90, l)
	}

	if !arrayContains(ancestries, db.Names[0].Ancestry) {
		t.Errorf("Cannot build names database.")
	}

}
