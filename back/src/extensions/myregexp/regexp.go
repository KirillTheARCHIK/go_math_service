package myregexp

import "regexp"

type MyRegexp struct {
	regexp.Regexp
}

func (regex MyRegexp) ValueByGroupName(match []string, groupName string) string {
	for groupIndex, _groupName := range regex.SubexpNames() {
		if groupName == _groupName {
			return match[groupIndex]
		}
	}
	panic("Regexp group with name " + groupName + "does not exist")
}
