package sotdlgen

import (
	"encoding/hex"
	"hash/fnv"
	"math/rand"
	"strconv"
	"time"
)

var defaultSeed = time.Now().UTC().UnixNano()

func init() {
	rand.Seed(defaultSeed)
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func setSeed(charHash string) (string, error) {
	if charHash == "" {
		charHash = strconv.FormatInt(defaultSeed, 16)
	}
	src := []byte(charHash)
	dst := make([]byte, hex.DecodedLen(len(src)))
	seed, err := hex.Decode(dst, src)
	if err != nil {
		return charHash, err
	}
	rand.Seed(int64(seed))
	return charHash, nil
}

func sampleWithoutReplacement(choices []string, n int) []string {
	samples := []string{}
	idxs := rand.Perm(len(choices))
	for i := 0; i < n; i++ {
		samples = append(samples, choices[idxs[i]])
	}
	return samples
}

func randomChoice(choices []string) string {
	r := rand.Intn(len(choices))
	return choices[r]
}

func randomInt(min, max int) int {
	// Returns an int in [min,max).
	return rand.Intn(max-min) + min
}

func weightedRandomChoice(choices []string, weights []float64) string {
	sum := 0.0
	for _, w := range weights {
		sum += w
	}
	r := rand.Float64()*sum - 1.0
	total := 0.0
	for i, w := range weights {
		total += w
		if r <= total {
			return choices[i]
		}
	}
	return choices[0]
}
