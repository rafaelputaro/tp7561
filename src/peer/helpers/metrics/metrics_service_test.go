package peer_metrics

import (
	"testing"
	"tp/common/contact"
	common_metrics "tp/common/metrics"
)

func TestParseFileName(t *testing.T) {
	files := []string{"filec-1-2.txt.part3", "filec-2-12.txt.part4", "filec-1-2.txt"}
	numbers := []float64{2, 12, 2}
	for index, fileName := range files {
		found := common_metrics.ParseFileNumber(fileName)
		if found != numbers[index] {
			t.Errorf("Error found %v must be: %v", found, numbers[index])
		}
	}
}

func TestParseContact(t *testing.T) {
	contacts := []contact.Contact{
		*contact.NewContact([]byte{1}, "peer-1:8080"),
		*contact.NewContact([]byte{2}, "peer-20:9020"),
	}
	contactNames := []string{"peer-1", "peer-20"}
	for index, aContact := range contacts {
		found := parseContact(aContact)
		if found != contactNames[index] {
			t.Errorf("Error found %v must be: %v", found, contactNames[index])
		}
	}
}
