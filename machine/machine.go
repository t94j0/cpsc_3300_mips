package machine

import (
	"fmt"
	binary "../binary"
	"io"
	"os"
)

type Machine struct {
	// ir is the current instruction
	ir uint32
	// pc is the next instruction
	pc uint32
	// when halt is true, the program ends
	halt bool

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
}

func (m *Machine) LoadFromStdin() error {
	return m.LoadFromReader(os.Stdin)
}

func (m *Machine) LoadFromReader(r io.Reader) error {
	m.memory = make([]uint32, 0)

	for {
		var line uint32
		_, err := fmt.Fscanf(r, "%x\n", &line)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		m.memory = append(m.memory, line)
	}

	return nil
}

func (m *Machine) PrintMemory() {
	fmt.Println("contents of memory")
	fmt.Println("addr value")
	for i, word := range m.memory {
		fmt.Printf("%03x: %08x\n", i, word)
	}
}

func (m *Machine) PrintInstructionClassCounts() {
	loadStore := m.memoryAccess.load + m.memoryAccess.store
	jumpBranch := m.transferControl.jump + m.transferControl.jumpLink +
		m.transferControl.takenBranch + m.transferControl.untakenBranch
	total := m.instructionClass.alu + loadStore + jumpBranch

	fmt.Println()
	fmt.Printf(`instruction class counts (omits hlt instruction)
  alu ops            %3d
  loads/stores       %3d
  jumps/branches     %3d
total                %3d`, m.instructionClass.alu, loadStore, jumpBranch, total)
	fmt.Println()
	fmt.Println()
}
func (m *Machine) PrintMemoryAccessCounts() {
	iF := m.memoryAccess.instFetch
	load := m.memoryAccess.load
	store := m.memoryAccess.store
	total := iF + load + store

	fmt.Printf(`memory access counts (omits hlt instruction)
  inst. fetches       %3d
  loads               %3d
  stores              %3d
total                 %3d`, iF, load, store, total)
	fmt.Println()
	fmt.Println()
}
func (m *Machine) PrintTransferControlCounts() {
	j := m.transferControl.jump
	jal := m.transferControl.jumpLink
	taken := m.transferControl.takenBranch
	untaken := m.transferControl.untakenBranch
	total := j + jal + taken + untaken

	fmt.Printf(`transfer of control counts
  jumps               %3d
  jump-and-links      %3d
  taken branches      %3d
  untaken branches    %3d
total                 %3d`, j, jal, taken, untaken, total)
	fmt.Println()
}

func (m *Machine) PrintBehavorialSimulation() {
	fmt.Println("\n" + `behavioral simulation of simple MIPS-like machine
  (all values are shown in hexadecimal)`)
	fmt.Println()
}

func (m *Machine) next() {
	m.ir = m.pc
	m.pc++
}

// fetchInstruction parses the current instruction into:
func (m *Machine) fetchOp() uint16 {
	return uint16(m.memory[m.ir] >> 26)
}

func (m *Machine) Execute() {
	fmt.Println("pc   result of instruction at that location")
	for !m.halt {
		m.next()
		m.memoryAccess.instFetch++
		inst := uint16(binary.GetOperation(m.memory[m.ir]))
		if _, ok := opcodeInstructions[inst]; !ok {
			panic(fmt.Sprintf("opcode does not exist: %08x\n", m.memory[m.ir]))
		}
		opcodeInstructions[inst](m, m.memory[m.ir])
	}
	fmt.Println()
}
