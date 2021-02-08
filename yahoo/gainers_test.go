package yahoo

import (
	"strings"
	"testing"
)

func Test_GetTopGainersTickers(t *testing.T) {
	arr := GetTopGainersTickers()
	println(strings.Join(arr, " | "))
}
