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

	expect := 3
	if value != expect {
		t.Fatalf("Expected %d: %d", expect, value)
	}
}

func TestLoadIndirect(t *testing.T) {
	computer := New()
	computer.Memory[PC_START] = 0b1010001000000001
	computer.Memory[PC_START+1] = 2
	computer.Memory[PC_START+2] = PC_START+1
	computer.Step()
	value := computer.Registers[R_1]

	expect := 2
	if value != expect {
		t.Fatalf("Expected %d: %d", expect, value)
	}
}
