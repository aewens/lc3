package vm

import (
	"testing"
)

func TestAddRegister(t *testing.T) {
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

func TestAddImmediate(t *testing.T) {
	computer := New()
	computer.Registers[R_0] = 1
	computer.Memory[PC_START] = 0b0001010000111110
	computer.Step()
	value := computer.Registers[R_2]

	expect := 0xFFFF
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

func TestAndRegister(t *testing.T) {
	computer := New()
	computer.Registers[R_0] = 0xFF
	computer.Registers[R_1] = 0xF0
	computer.Memory[PC_START] = 0b0101010000000001
	computer.Step()
	value := computer.Registers[R_2]

	expect := 0xF0
	if value != expect {
		t.Fatalf("Expected %d: %d", expect, value)
	}
}

func TestAndImmediate(t *testing.T) {
	computer := New()
	computer.Registers[R_0] = 0xFF
	computer.Memory[PC_START] = 0b0101010000101111
	computer.Step()
	value := computer.Registers[R_2]

	expect := 0xF
	if value != expect {
		t.Fatalf("Expected %d: %d", expect, value)
	}
}

func TestNot(t *testing.T) {
	computer := New()
	computer.Registers[R_0] = 0xF0
	computer.Memory[PC_START] = 0b1001001000111111
	computer.Step()
	value := computer.Registers[R_1]

	expect := 0xF
	if value != expect {
		t.Fatalf("Expected %d: %d", expect, value)
	}
}

func TestBranchN(t *testing.T) {
	computer := New()
	computer.Registers[R_0] = 0xFFFE
	computer.Registers[R_1] = 1
	computer.Memory[PC_START] = 0b0001010000000001
	computer.Memory[PC_START+1] = 0b0000100000000001
	computer.Memory[PC_START+2] = 3
	computer.Step()
	computer.Step()
	value := computer.readMemoryNext()

	expect := 3
	if value != expect {
		t.Fatalf("Expected %d: %d", expect, value)
	}
}

func TestBranchZ(t *testing.T) {
	computer := New()
	computer.Registers[R_0] = 0
	computer.Registers[R_1] = 0
	computer.Memory[PC_START] = 0b0001010000000001
	computer.Memory[PC_START+1] = 0b0000100000000001
	computer.Memory[PC_START+2] = 3
	computer.Step()
	computer.Step()
	value := computer.readMemoryNext()

	expect := 3
	if value != expect {
		t.Fatalf("Expected %d: %d", expect, value)
	}
}

func TestBranchP(t *testing.T) {
	computer := New()
	computer.Registers[R_0] = 1
	computer.Registers[R_1] = 1
	computer.Memory[PC_START] = 0b0001010000000001
	computer.Memory[PC_START+1] = 0b0000100000000001
	computer.Memory[PC_START+2] = 3
	computer.Step()
	computer.Step()
	value := computer.readMemoryNext()

	expect := 3
	if value != expect {
		t.Fatalf("Expected %d: %d", expect, value)
	}
}
