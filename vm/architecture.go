package vm

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

type LC3 struct {
	Memory    [MEMORY_LOCATIONS]uint16
	Registers [R_COUNT]uint16
}

func New() *LC3 {
	return &LC3{}
}

func (self *LC3) Run(program chan string) []uint16 {
	return []uint16{}
}
