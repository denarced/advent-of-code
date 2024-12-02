package shared

import (
	"fmt"
	"os"
)

func Abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func Die(err error, message string) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Die. Message: \"%s\". Error: %s.\n", message, err)
	//revive:disable-next-line:deep-exit
	os.Exit(2)
}
