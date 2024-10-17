package utils

import (
	"regexp"
)

func IsValidDateFormat(dateString string) bool {
	const dateFormat = `^\d{4}-\d{2}-\d{2}$`
	re := regexp.MustCompile(dateFormat)
	return re.MatchString(dateString)
}
