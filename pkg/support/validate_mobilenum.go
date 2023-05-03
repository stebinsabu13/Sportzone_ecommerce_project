package support

import "regexp"

func MobileNum_validater(number string) bool {
	re := regexp.MustCompile(`^[0-9]{10}$`)
	if re.MatchString(number) {
		return true
	} else {
		return false
	}
}
