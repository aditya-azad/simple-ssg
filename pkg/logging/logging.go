package logging

import (
	"fmt"
	"os"
)

func Error(format string, args ...any) {
	fmt.Printf("ERROR: %s\n", fmt.Sprintf(format, args...))
	os.Exit(1)
}
