package yahoo

import (
	"strings"
	"testing"
)

func Test_GetTopGainersTickers(t *testing.T) {
	arr := GetTopMoversTickers(true)
	println(strings.Join(arr, " | "))
	arr = GetTopMoversTickers(false)
	println(strings.Join(arr, " | "))
}
