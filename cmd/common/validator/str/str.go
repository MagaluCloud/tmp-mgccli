package str

import (
	"fmt"
	"slices"
	"strings"
)

func Validator(value string, validateTag string) error {

	if strings.Contains(validateTag, "oneof=") {
		oneof := strings.Split(validateTag, "oneof=")[1]
		oneofValues := strings.Split(oneof, ",")
		if slices.Contains(oneofValues, value) {
			return nil
		}
		return fmt.Errorf("value must be one of %s", oneofValues)
	}

	return nil
}
