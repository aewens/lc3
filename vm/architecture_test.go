package vm

import (
	"testing"
)

func TestAdd(t *testing.T) {
	computer := New()
	computer.Registers[R_0] = 1
	computer.Registers[R_1] = 2
	computer.Memory[PC_START] = 0b0001010000000001
	computer.Step()
	value := computer.Registers[R_2]

	if value != 3 {
		t.Fatalf("Expected 3: %d", value)
	}
}
