package main

import (
	"fmt"
	"strings"
)

func banner(str string) {
	b := strings.Repeat("-", len(str))
	fmt.Printf("%s\n%s\n%s\n", b, str, b)
}

func bannerf(str string, args ...interface{}) {
	str = sf(str, args...)
	banner(str)
}
