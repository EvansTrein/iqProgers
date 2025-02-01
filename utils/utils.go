package utils

import "regexp"

func IsGUID(s string) bool {
	guidRegex := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	ok, _ := regexp.MatchString(guidRegex, s)
	return ok
}
