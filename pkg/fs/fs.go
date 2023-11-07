package fs

import "os"

func IsDir(path *string) bool {
	info, err := os.Stat(*path)
	if err != nil || !info.IsDir() {
		return false
	}
	return true
}

