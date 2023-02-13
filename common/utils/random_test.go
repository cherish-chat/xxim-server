package utils

import "testing"

func TestRandomUtil_String(t *testing.T) {
	s := Random.String([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}, 6)
	t.Log(s)
}

func BenchmarkRandomUtil_String(b *testing.B) {
	strings := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	for i := 0; i < b.N; i++ {
		Random.String(strings, 6)
	}
}
