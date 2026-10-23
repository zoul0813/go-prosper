package vm

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "output")
	if err != nil {
		t.Fatal(err)
	}
	previous := os.Stdout
	os.Stdout = f
	defer func() { os.Stdout = previous; f.Close() }()
	fn()
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	output, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return string(output)
}

func TestSourceCommentsAndStrings(t *testing.T) {
	m := loadCore(t)
	output := captureOutput(t, func() {
		processTest(t, m, ": c 1 \\ if ; loop ignored\n ( ; if loop) 2 + ; c")
		processTest(t, m, `." café (text) \"quoted\""`)
	})
	assertStack(t, m, 3)
	if output != `café (text) "quoted"` {
		t.Fatalf("output = %q", output)
	}
	processTest(t, m, ": multiline")
	processTest(t, m, "4 5 + ; multiline")
	assertStack(t, m, 3, 9)
}

func TestSourceFailuresRecover(t *testing.T) {
	for _, source := range []string{
		`( unclosed`, `." unclosed`, `: bad abort" unclosed`, `: bad if ;`,
		`: unfinished 1`, `: bad +loop ;`, `: bad i ; bad`, `: bad j ; bad`,
	} {
		t.Run(source, func(t *testing.T) {
			m := loadCore(t)
			if err := m.Load(source); err == nil {
				t.Fatal("invalid source accepted")
			}
			processTest(t, m, "1 2 +")
			assertStack(t, m, 3)
		})
	}
}

func TestPlusLoop(t *testing.T) {
	for _, tt := range []struct {
		code string
		want []int
	}{
		{": c 5 0 do i 2 +loop ; c", []int{0, 2, 4}},
		{": c 0 2 do i -1 +loop ; c", []int{2, 1, 0}},
		{": c 0 3 do i -2 +loop ; c", []int{3, 1}},
		{": c 0 0 do i -1 +loop ; c", []int{0}},
		{": c 2 0 do 2 0 do j i loop loop ; c", []int{0, 0, 0, 1, 1, 0, 1, 1}},
		// A zero increment repeats the same index; the body changes the next step.
		{"variable step 0 step ! : c 2 0 do i step @ 1 step ! +loop ; c", []int{0, 0, 1}},
	} {
		t.Run(tt.code, func(t *testing.T) {
			m := loadCore(t)
			processTest(t, m, tt.code)
			assertStack(t, m, tt.want...)
			if m.Memory.ReturnStack().Index != 0 {
				t.Fatal("loop parameters remain")
			}
		})
	}
	// Crossing the signed boundary must not terminate a loop prematurely.
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	m := loadCore(t)
	processTest(t, m, ": c "+strconv.Itoa(minInt+1)+" "+strconv.Itoa(maxInt)+" do i loop ; c")
	assertStack(t, m, maxInt, minInt)
}

func TestAbortAndKey(t *testing.T) {
	m := loadCore(t)
	processTest(t, m, `: check abort" excess increment " ; 42 0 check`)
	assertStack(t, m, 42)
	if err := m.Execute(`: fail 3 0 do 1 check loop ; fail`); err == nil || !strings.Contains(err.Error(), "excess increment") {
		t.Fatalf("abort error = %v", err)
	}
	assertStack(t, m)
	if m.Memory.ReturnStack().Index != 0 {
		t.Fatal("abort retained loop state")
	}
	m.ReadKey = func() (rune, error) { return 'w', nil }
	processTest(t, m, "key")
	assertStack(t, m, int('w'))
	m.ReadKey = func() (rune, error) { return 0, ErrBye }
	if err := m.Execute("key"); !errors.Is(err, ErrBye) {
		t.Fatalf("key error = %v", err)
	}
	assertStack(t, m)
}

func TestUnsignedOutput(t *testing.T) {
	m := loadCore(t)
	output := captureOutput(t, func() { processTest(t, m, "42 5 u.r hex 255 5 u.r decimal -1 0 u.r") })
	want := "   42  255" + strconv.FormatUint(uint64(^uint(0)), 10)
	if output != want {
		t.Fatalf("output = %q, want %q", output, want)
	}
	assertStack(t, m)
}
