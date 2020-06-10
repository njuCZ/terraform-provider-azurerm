package validate

import (
	"fmt"
	"regexp"
)

func EventHubClusterName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// The name attribute rules are :
	// 1. can contain only letters, numbers and hyphens.
	// 2. The first character must be a letter.
	// 3. The last character must be a letter or number
	// 3. The value must be between 6 and 50 characters long
	if !regexp.MustCompile(`^([a-zA-Z])([a-zA-Z\d-]{4,48})([a-zA-Z\d])$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s can contain only letters, numbers and hyphens. must start with a letter, and it must end with a letter or number. The value must be between 6 and 50 characters long", k))
	}

	return
}
