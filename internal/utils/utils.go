package utils

import "os"

// ContainsInt checks if a given slice contains a given integer
func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// PathExists checks if a given directory exists
func PathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		// Invalid folder
		return false
	}

	return true
}
