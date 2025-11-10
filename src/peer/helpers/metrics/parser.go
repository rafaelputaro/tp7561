package peer_metrics

import (
	"strings"
	"tp/common/contact"
)

func parseContact(contact contact.Contact) string {
	return parseUrlToId(contact.Url)
}

func parseUrlToId(url string) string {
	splited := strings.Split(url, ":")
	return splited[0]
}
