package utils

import (
	"testing"
)

func TestAnySet(t *testing.T) {
	set := AnySet([]string{"a", "c", "c", "a", "d"})
	t.Logf("%v", set)
	set = AnySet([]string{})
	t.Logf("%v", set)
}
