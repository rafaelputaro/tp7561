package peer_metrics

import (
	"regexp"
	"strconv"
	"strings"
	"tp/common/contact"
)

// Retorna el segundo número que aparece en el nombre de un archivo
func parseFileNumber(fileName string) float64 {
	patron := `(\d+)`
	ren := regexp.MustCompile(patron)
	found := ren.FindAllString(fileName, -1)
	//common.Log.Debugf("%v", found)
	if len(found) > 1 {
		converted, err := strconv.Atoi(found[1])
		if err == nil {
			return float64(converted)
		}
	}
	return -1
}

func parseContact(contact contact.Contact) string {
	return parseUrlToId(contact.Url)
}

func parseUrlToId(url string) string {
	splited := strings.Split(url, ":")
	return splited[0]
}
