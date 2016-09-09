package askUser

import (
	"fmt"
	"strings"
)

// Supports a question and returns a bool as true if they said 'yes' or false
// if 'no'
func Args(question string, qtype bool) (bool, string) {
	var answer string

	fmt.Print(question)
	fmt.Scanln(&answer)

	if qtype {
		answer = strings.ToLower(answer)

		if answer != "y" {
			return false, ""
		}
		return true, ""
	}

	return true, answer
}
