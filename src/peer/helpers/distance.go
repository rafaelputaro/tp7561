package helpers

import (
	"errors"
	"math"
)

const INFINITY_DISTANCE = math.MaxInt32
const MSG_KEY_MUST_HAVE_SAME_LENGTH = "slices must have the same length"

var distanceScale = initDistanceScale()

// Retorna la distancia logarítmica entre dos key's. De esta manera claves
// con una distancia del mismo orden quedan agrupadas bajo la misma distancia
// loragítmica
func GetLogDistance(key1 []byte, key2 []byte) (int, error) {
	if len(key1) != len(key2) {
		return INFINITY_DISTANCE, errors.New(MSG_KEY_MUST_HAVE_SAME_LENGTH)
	}
	result := 0
	for i := range key1 {
		dist := int(math.Log10(math.Abs(float64(key1[i] ^ key2[i]))))
		if dist > 0 {
			result += dist + distanceScale[i]
		}
	}
	return result, nil
}

func initDistanceScale() []int {
	toReturn := []int{}
	for nByte := range LENGTH_KEY_IN_BYTES {
		toReturn = append(toReturn, int(math.Log10(math.Pow(2.0, float64(8*nByte)))))
	}
	return toReturn
}
