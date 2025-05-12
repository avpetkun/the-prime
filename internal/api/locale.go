package api

import "strings"

func parseAndGetLocale(textWithLocals, userLanguageCode string) (localized string) {
	userLanguageCode += ":"
	lines := strings.Split(textWithLocals, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, userLanguageCode) {
			return strings.TrimSpace(line[len(userLanguageCode):])
		}
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "en:") {
			return strings.TrimSpace(line[3:])
		}
	}
	return textWithLocals
}
