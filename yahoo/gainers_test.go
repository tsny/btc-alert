package yahoo

import (
	"testing"
)

func Test_GetSummary(t *testing.T) {
	sum := GetSummary("MSFT")
	if sum == "" {
		t.Fatal("summary was empty")
	}
}

func TestGetGainers(t *testing.T) {
	arr := GetGainers()
	if len(arr) == 0 {
		t.Fatal("got no gainers")
	}
	t.Logf("%v", arr)
}
