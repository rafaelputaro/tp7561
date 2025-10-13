package keys

import (
	"testing"
)

func TestDistance(t *testing.T) {
	key1 := GetKey("dfadfasdfjads9fi32fdasfasdfadfs")
	key2 := GetKey("dfadfasdfjads9fi32fdasfasdfadfs")
	key3 := GetKey("2")
	dist, _ := GetLogDistance(key1, key2)
	println(dist)
	for i := range LENGTH_KEY_IN_BYTES {
		key2[i] = key3[i]
		dist, err := GetLogDistance(key1, key2)
		if err != nil {
			t.Errorf("Error must be null")
		}
		println(dist)
	}
}

func TestDistanceScales(t *testing.T) {
	for i := range distanceScale {
		println(distanceScale[i])
	}

}
