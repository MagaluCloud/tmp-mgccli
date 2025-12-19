package cmdutils

import "strings"

func StringToSlice(s, sep string, shouldTrim bool) []string {
	entries := strings.Split(s, sep)

	result := make([]string, 0)
	if shouldTrim {
		for _, str := range entries {
			newValue := strings.TrimSpace(str)
			if newValue == "" {
				continue
			}
			result = append(result, newValue)
		}
	}

	return result
}
