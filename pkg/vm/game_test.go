package vm

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func loadGame(t *testing.T, path string) *VM {
	t.Helper()
	m := loadCore(t)
	source, err := os.ReadFile("../../" + path)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Load(string(source)); err != nil {
		t.Fatal(err)
	}
	assertStack(t, m)
	return m
}

func gameAddress(t *testing.T, m *VM, name string) int {
	t.Helper()
	processTest(t, m, name)
	return m.Memory.Stack().Pop()
}

// A compact, independent 2048 rule implementation: compact each line, merge
// adjacent equal pairs once, then compact again in the direction of movement.
func expectedMove(grid []int, direction string) ([]int, int) {
	result := make([]int, 16)
	score := 0
	for line := 0; line < 4; line++ {
		indices := make([]int, 4)
		for p := range indices {
			switch direction {
			case "left":
				indices[p] = line*4 + p
			case "right":
				indices[p] = line*4 + 3 - p
			case "up":
				indices[p] = p*4 + line
			case "down":
				indices[p] = (3-p)*4 + line
			}
		}
		var tiles []int
		for _, i := range indices {
			if grid[i] != 0 {
				tiles = append(tiles, grid[i])
			}
		}
		out := 0
		for p := 0; p < len(tiles); p++ {
			n := tiles[p]
			if p+1 < len(tiles) && n == tiles[p+1] {
				n *= 2
				score += n
				p++
			}
			result[indices[out]] = n
			out++
		}
	}
	return result, score
}

func TestGameMoves(t *testing.T) {
	for _, file := range []string{"2048.fs"} {
		t.Run(file, func(t *testing.T) {
			m := loadGame(t, file)
			grid := gameAddress(t, m, "grid")
			score := gameAddress(t, m, "score")
			moved := gameAddress(t, m, "moved")
			values := []int{0, 2, 4, 8}
			for _, direction := range []string{"left", "right", "up", "down"} {
				for pattern := 0; pattern < 256; pattern++ {
					board := make([]int, 16)
					for r := 0; r < 4; r++ {
						for c := 0; c < 4; c++ {
							board[r*4+c] = values[((pattern>>(2*c))+r)%4]
						}
					}
					copy(m.Memory.DataSpace().Values[grid:grid+16], board)
					m.Memory.DataSpace().Put(score, 0)
					want, wantScore := expectedMove(board, direction)
					processTest(t, m, direction)
					got := m.Memory.DataSpace().Values[grid : grid+16]
					if !reflect.DeepEqual(got, want) || m.Memory.DataSpace().Get(score) != wantScore {
						t.Fatalf("%s %v: got %v score %d; want %v score %d", direction, board, got, m.Memory.DataSpace().Get(score), want, wantScore)
					}
					if (m.Memory.DataSpace().Get(moved) != 0) != !reflect.DeepEqual(board, want) {
						t.Fatalf("incorrect changed flag for %s %v", direction, board)
					}
					assertStack(t, m)
					if m.Memory.ReturnStack().Index != 0 {
						t.Fatal("movement leaked loop parameters")
					}
				}
			}
		})
	}
}

func TestGameRandomAndFullBoard(t *testing.T) {
	m := loadGame(t, "2048.fs")
	grid := gameAddress(t, m, "grid")
	processTest(t, m, "1 seed !")
	for i := 0; i < 1000; i++ {
		processTest(t, m, "16 random")
		n := m.Memory.Stack().Pop()
		if n < 0 || n >= 16 {
			t.Fatalf("random index %d", n)
		}
		assertStack(t, m)
	}
	processTest(t, m, "init-grid")
	for count := 1; count <= 16; count++ {
		processTest(t, m, "new")
		occupied := 0
		for _, n := range m.Memory.DataSpace().Values[grid : grid+16] {
			if n == 2 {
				occupied++
			} else if n != 0 {
				t.Fatalf("unexpected tile %d", n)
			}
		}
		if occupied != count {
			t.Fatalf("occupied = %d, want %d", occupied, count)
		}
		assertStack(t, m)
	}
	processTest(t, m, "new") // A full grid must return, not search forever.
	assertStack(t, m)
}

func TestGamePlayAndReplay(t *testing.T) {
	m := loadGame(t, "2048.fs")
	score := gameAddress(t, m, "score")
	for _, input := range []string{"wasdxq", "q"} {
		m.Memory.DataSpace().Put(score, 999)
		keys := strings.NewReader(input)
		m.ReadKey = func() (rune, error) { r, _, err := keys.ReadRune(); return r, err }
		processTest(t, m, "play")
		assertStack(t, m)
		if m.Memory.ReturnStack().Index != 0 {
			t.Fatal("play leaked return stack")
		}
		if m.Memory.DataSpace().Get(score) == 999 {
			t.Fatal("score was not reset")
		}
		if keys.Len() != 0 {
			t.Fatal("play returned before quit")
		}
	}
	m.ReadKey = func() (rune, error) { return 0, io.EOF }
	if err := m.Execute("play"); err != io.EOF {
		t.Fatalf("input error = %v", err)
	}
	assertStack(t, m)
	if m.Memory.ReturnStack().Index != 0 {
		t.Fatal("input failure leaked return stack")
	}
}

func TestOriginalGameLoadsAndQuits(t *testing.T) {
	m := loadGame(t, "2048_raw.fs")
	m.ReadKey = func() (rune, error) { return 'q', nil }
	processTest(t, m, "play")
	assertStack(t, m)
	if m.Memory.ReturnStack().Index != 0 {
		t.Fatal("original game leaked return stack")
	}
}

func TestGameDoesNotSpawnAfterNoOp(t *testing.T) {
	m := loadGame(t, "2048.fs")
	grid := gameAddress(t, m, "grid")
	processTest(t, m, "1 seed !") // The initial tile is in the top-left corner.
	reads := 0
	m.ReadKey = func() (rune, error) {
		reads++
		if reads == 1 {
			return 'a', nil
		}
		occupied := 0
		for _, tile := range m.Memory.DataSpace().Values[grid : grid+16] {
			if tile != 0 {
				occupied++
			}
		}
		if occupied != 1 {
			t.Fatalf("no-op move spawned a tile: %v", m.Memory.DataSpace().Values[grid:grid+16])
		}
		return 'q', nil
	}
	processTest(t, m, "play")
	assertStack(t, m)
}
