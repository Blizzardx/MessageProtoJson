package common

import "strings"

func FirstLetterToUpper(letter string) string {
	size := strings.Count(letter, "") - 1
	if size > 1 {
		firstL := string(letter[0])
		otherL := string(letter[1:])
		return strings.ToUpper(firstL) + otherL
	} else {
		return strings.ToUpper(letter)
	}
}
func FirstLetterToLower(letter string) string {
	size := strings.Count(letter, "") - 1
	if size > 1 {
		firstL := string(letter[0])
		otherL := string(letter[1:])
		return strings.ToLower(firstL) + otherL
	} else {
		return strings.ToLower(letter)
	}
}
