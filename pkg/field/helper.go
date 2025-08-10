package field

func contains(values []string, targetValues ...string) bool {
	for _, s := range values {
		for _, value := range targetValues {
			if s == value {
				return true
			}
		}
	}
	return false
}
