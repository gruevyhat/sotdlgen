package sotdlgen

import (
	"encoding/binary"
	"encoding/hex"
	"hash/fnv"
	"math/rand"
	"strconv"
	"time"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func setSeed(charHash string) (string, error) {
	if charHash == "" {
		defaultSeed := time.Now().UTC().UnixNano()
		charHash = strconv.FormatInt(defaultSeed, 16)
	}
	h, err := hex.DecodeString(charHash)
	if err != nil {
		return charHash, err
	}
	seed := binary.BigEndian.Uint64(h)
	rand.Seed(int64(seed))
	log.Info("Set new seed:", seed)
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
