package peer_metrics

import (
	"testing"
)

func TestParseFileName(t *testing.T) {
	files := []string{"filec-1-2.txt.part3", "filec-2-12.txt.part4", "filec-1-2.txt"}
	numbers := []int{2, 12, 2}
	for index, fileName := range files {
		found := parseFileNumber(fileName)
		if found != numbers[index] {
			t.Errorf("Error found %v must be: %v", found, numbers[index])
		}
	}
}
