package game

import "strings"

func hasTag(tag string, tags ...string) bool {
	for _, t := range tags {
		if strings.Contains(tag, t) {
			return true
		}
	}
	return false
}

func getFirstMatchingTag(tag string, tags ...string) string {
	for _, t := range tags {
		if strings.Contains(tag, t) {
			return t
		}
	}
	return ""
}
