package machine

import (
	"fmt"

	binary "github.com/t94j0/go-mips-instruction-format"
)

type Machine struct {
	// ir is the current instruction
	ir uint32
	// pc is the next instruction
	pc uint32
	// when halt is true, the program ends
	halt bool

	writeTo int

	memory    []uint32
	registers [32]uint32

	instructionClass struct {
		alu uint64
	}

	memoryAccess struct {
		instFetch uint64
		load      uint64
		store     uint64
	}

	transferControl struct {
		jump          uint64
		jumpLink      uint64
		takenBranch   uint64
		untakenBranch uint64
	}

	pipeline *Pipeline
}

func NewMachine() *Machine {
	return &Machine{writeTo: -1}
}

func (m *Machine) cycle() {
	if m.registers[0] != 0 {
		m.registers[0] = 0
	}

	m.writeTo = -1

	m.ir = m.pc
	m.pc++
}

func (m *Machine) getOperations() (uint16, uint16, uint32) {
	op := m.memory[m.ir]
	inst := uint16(binary.GetOperation(op))
	funct := uint16(binary.GetFunct(op))
	return inst, funct, op
}

func (m *Machine) getNextOp() (uint16, uint16, uint32) {
	op := m.memory[m.pc]
	inst := uint16(binary.GetOperation(op))
	funct := uint16(binary.GetFunct(op))
	return inst, funct, op
}

func (m *Machine) runInstruction() {
	m.cycle()
	m.memoryAccess.instFetch++
	inst, _, _ := m.getOperations()
	if _, ok := opcodeInstructions[inst]; !ok {
		panic(fmt.Sprintf("opcode does not exist: %08x\n", m.memory[m.ir]))
	}
	opcodeInstructions[inst](m, m.memory[m.ir])
}

func (m *Machine) Execute() {
	fmt.Println("instruction pairing analysis")
	m.pipeline = NewPipeline(m)
	for !m.halt {
		m.pipeline.Schedule()
	}
	fmt.Println()
}
