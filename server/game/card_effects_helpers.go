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

func findAnyUnit(ps *PlayerState) *CardInstance {
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			if ps.Units[col][row] != nil {
				return ps.Units[col][row]
			}
		}
	}
	return nil
}

func findFrontRowUnit(ps *PlayerState) *CardInstance {
	frontRow := ps.GetFrontRow()
	if frontRow < 0 {
		return nil
	}
	for col := 0; col < 3; col++ {
		if ps.Units[col][frontRow] != nil {
			return ps.Units[col][frontRow]
		}
	}
	return nil
}
