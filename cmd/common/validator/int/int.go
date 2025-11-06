package integer

import (
	"fmt"
	"strconv"
	"strings"
)

func Validator(value int, validateTag string) error {
	minimumValue, maximumValue := 0, 0
	var err error

	values := strings.Split(validateTag, ",")

	for _, value := range values {
		if strings.Contains(value, "minimum=") {
			minimum := strings.Split(value, "minimum=")[1]
			minimumValue, err = strconv.Atoi(minimum)
			if err != nil {
				return err
			}
			continue
		}

		if strings.Contains(value, "maximum=") {
			maximum := strings.Split(value, "maximum=")[1]
			maximumValue, err = strconv.Atoi(maximum)
			if err != nil {
				return err
			}
			continue
		}
	}

	if err != nil {
		return err
	}

	if minimumValue > 0 && value < minimumValue {
		return fmt.Errorf("value %d must be greater than %d", value, minimumValue)
	}

	if maximumValue > 0 && value > maximumValue {
		return fmt.Errorf("value %d must be less than %d", value, maximumValue)
	}

	return nil
}
