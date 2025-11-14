package str

import (
	"fmt"
	"strings"
)

func Validator(value string, validateTag string) error {

	if strings.Contains(validateTag, "oneof=") {
		oneof := strings.Split(validateTag, "oneof=")[1]
		oneofValues := strings.Split(oneof, ",")
		for _, oneofValue := range oneofValues {
			if oneofValue == value {
				return nil
			}
		}
		return fmt.Errorf("value %s must be one of %s", value, oneofValues)
	}

	return nil
}
