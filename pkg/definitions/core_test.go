package definitions_test

import (
	"reflect"
	"testing"

	"github.com/zoul0813/go-prosper/pkg/definitions"
	"github.com/zoul0813/go-prosper/pkg/memory"
)

func TestDivMod(t *testing.T) {
	for _, tt := range []struct {
		left, right, remainder, quotient int
	}{
		{7, 3, 1, 2}, {-7, 3, -1, -2}, {7, -3, 1, -2}, {6, 3, 0, 2},
	} {
		m := memory.NewMemory(16, 16)
		m.Stack().Push(tt.left)
		m.Stack().Push(tt.right)
		if _, err := definitions.Core["/mod"](m); err != nil {
			t.Fatal(err)
		}
		want := []int{tt.remainder, tt.quotient}
		if got := m.Stack().Values[:m.Stack().Index]; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d %d /mod = %v, want %v", tt.left, tt.right, got, want)
		}
	}
}
