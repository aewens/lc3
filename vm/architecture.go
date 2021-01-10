package vm

import (
	"fmt"
)

const MEMORY_LOCATIONS = 65536

const (
	R_R0 int = iota
	R_R1
	R_R2
	R_R3
	R_R4
	R_R5
	R_R6
	R_R7
	R_PC   // Program Counter Register
	R_COND // Condition Register
	R_COUNT
)

const (
	OP_BR   int = iota // Branch
	OP_ADD             // Add
	OP_LD              // Load
	OP_ST              // Store
	OP_JSR             // Jump Register
	OP_AND             // Bitwise And
	OP_LDR             // Load Register
	OP_STR             // Store Register
	OP_RTI             // Return From Interruption (unused)
	OP_NOT             // Bitwise Not
	OP_LDI             // Load Indirect
	OP_STI             // Store Indirect
	OP_JMP             // Jump
	OP_RES             // Reserved (unused)
	OP_LEA             // Load Effective Address
	OP_TRAP            // Execute Trap
	OP_COUNT
)

const (
	FL_P = 1 << 0 // Positive
	FL_Z = 1 << 1 // Zero
	FL_N = 1 << 2 // Negative
)

const PC_START = 0x3000 // Default magic number

type LC3 struct {
	Memory    [MEMORY_LOCATIONS]int
	Registers [R_COUNT]int
	Halted    bool
}

func New() *LC3 {
	self := &LC3{Halted: false}
	self.Registers[R_PC] = PC_START
	return self
}

func (self *LC3) checkLocation(location int) {
	/*
		Using uint16 could avoid the need for this function
		However that would make it more difficult to extend the VM later on
	*/
	if location < 0 || location >= MEMORY_LOCATIONS {
		panic(fmt.Sprintf("Location '%d' is out of bounds", location))
	}
}

func (self *LC3) checkRegister(index int) {
	if index >= R_COUNT {
		panic(fmt.Sprintf("Register '%d' is invalid", index))
	}
}

func (self *LC3) readMemory(location int) int {
	self.checkLocation(location)
	return self.Memory[location]
}

func (self *LC3) signExtend(value int, bits int) int {
	// Check if value is negative / if first bit is a 1 using 2's complement
	if (value>>(bits-1))&1 == 1 {
		// Left pad of ones to fit uint16 assuming 2's complement bit
		leftPad := (0xFFFF << bits) & 0xFFFF
		value = value | leftPad
	}
	return value
}

func (self *LC3) updateFlags(index int) {
	if index < 0 || index > R_R7 {
		panic(fmt.Sprintf("Invalid register for flag update: %d", index))
	}

	value := self.Registers[index]
	if value == 0 {
		self.Registers[R_COND] = FL_Z
		return
	}

	if (value >> 15) == 1 {
		self.Registers[R_COND] = FL_N
		return
	}

	self.Registers[R_COND] = FL_P
}

func (self *LC3) Step() {
	if self.Halted {
		return
	}

	counter := self.Registers[R_PC]
	self.Registers[R_PC] = counter + 1

	instruction := self.readMemory(counter)
	op := instruction >> 12
	/*
		Okay, so instructions in LC-3 are 16 bits wide used as follows:
		* 0x0000-0x00FF: Trap Vector Table
		* 0x0100-0x01FF: Interrupt Vector Table
		* 0x0200-0x02FF: Operating System and Supervisor Stack
		* 0x0300-0xFDFF: Available for user programs
		* 0xFE00-0xFFFF: Device Register Addresses (i.e. the opcode)

		So bits 12 through 15 has the opcode, which is why we bit shift by 12
	*/

	/*
		Terminology:
		* DR: Destination Register (from R0 to R7)
		* SR1: Source Register #1 (from R0 to R7)
		* SR2: Source Register #2 (from R0 to R7)
		* IM5: 5-bit Immediate Value; uint5 (0 to 32) or int4 (-16 to 15)
	*/
	switch op {
	case OP_ADD:
		/*
			ADD Encoding:
			a. Register Mode - FEDC|BA9|876|5|43|210
				* Bits F-C: 0001 (opcode)
				* Bits B-9: DR
				* Bits 8-6: SR1
				* Bits _-5: 0 (Register Mode)
				* Bits 4-3: 00 (unused)
				* Bits 2-0: SR2
			b. Immediate Mode - FEDC|BA9|876|5|43210
				* Bits F-C: 0001 (opcode)
				* Bits B-9: DR
				* Bits 8-6: SR1
				* Bits _-5: 1 (Immediate Mode)
				* Bits 4-0: IM5
		*/

		dr := (instruction >> 0x9) & 0b111
		self.checkRegister(dr)

		sr1 := (instruction >> 0x6) & 0b111
		self.checkRegister(sr1)

		mode := (instruction >> 0x5) & 0b1
		if mode == 1 {
			im5 := self.signExtend(instruction&0b11111, 5)
			self.Registers[dr] = self.Registers[sr1] + im5
		} else {
			sr2 := instruction & 0b111
			self.Registers[dr] = self.Registers[sr1] + self.Registers[sr2]
		}

		self.updateFlags(dr)
	//case OP_AND:
	//case OP_NOT:
	//case OP_BR:
	//case OP_JMP:
	//case OP_JSR:
	//case OP_LD:
	//case OP_LDI:
	//case OP_LDR:
	//case OP_LEA:
	//case OP_ST:
	//case OP_STI:
	//case OP_STR:
	//case OP_TRAP:
	case OP_RTI:
		fallthrough
	case OP_RES:
		fallthrough
	default:
		panic(fmt.Sprintf("Invalid opcode: %d", op))
	}
}

func (self *LC3) Run(program chan string) []int {
	for !self.Halted {
		self.Step()
	}

	return []int{}
}
