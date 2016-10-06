package pg

import "testing"

func TestGen(t *testing.T) {
	for _, c := range letters {
		t.Logf("%d", int(c))
	}
}
