package peer_metrics

import (
	"strings"
	"tp/common/contact"
)

// Transforma la url de un contacto en un id
func parseContact(contact contact.Contact) string {
	return parseUrlToId(contact.Url)
}

// Transforma una url en el id de un nodo
func parseUrlToId(url string) string {
	splited := strings.Split(url, ":")
	return splited[0]
}
