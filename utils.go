// Helper functions.
package sotdlgen

import (
	"io/ioutil"
	"strings"
)

var dataDir = getDataDir()

func getDataDir() string {
	dir := ""
	//if dir = os.Getenv("GOPATH"); dir != "" {
	//	dir += "/src/github.com/gruevyhat/sotdlgen/assets/"
	//}
	return dir + "./assets/"
}

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
