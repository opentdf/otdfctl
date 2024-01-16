package cli

func CommaSeparated(values []string) string {
	result := ""
	for i, v := range values {
		if i != 0 {
			result += ", "
		}
		result += v
	}
	return result
}
