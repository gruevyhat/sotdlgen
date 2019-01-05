// Helper functions.
package sotdlgen

import (
	"io/ioutil"
	"os"
	"strings"
)

var dataDir = os.Getenv("GOPATH") + "/src/github.com/gruevyhat/sotdlgen/assets/"

func readJson(filename string) []byte {
	raw, _ := ioutil.ReadFile(filename)
	return raw
}

func arrayContains(arr []string, s string) bool {
	for _, a := range arr {
		if a == s || strings.HasPrefix(a, s) || strings.HasSuffix(s, a) {
			return true
		}
	}
	return false
}

func arrayRemove(s string, a []string) []string {
	for i, x := range a {
		if x == "" || x == s {
			a = append(a[:i], a[i+1:]...)
		}
	}
	return a
}
