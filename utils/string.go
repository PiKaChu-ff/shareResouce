package utils

func ConcatWithPart(part string, s... string) (result string) {
	for i, v := range s {
		if i != 0 {
			result += part
		}
		result += v
	}
	return result
}
