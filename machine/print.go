package machine

import (
	"fmt"
)

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

	fmt.Printf(`instruction class counts (omits hlt instruction)
  alu ops           %3d
  loads/stores      %3d
  jumps/branches    %3d
total               %3d`, m.instructionClass.alu, loadStore, jumpBranch, total)
	fmt.Println()
	fmt.Println()
}
func (m *Machine) PrintMemoryAccessCounts() {
	iF := m.memoryAccess.instFetch
	load := m.memoryAccess.load
	store := m.memoryAccess.store
	total := iF + load + store

	fmt.Printf(`memory access counts (omits hlt instruction)
  inst. fetches     %3d
  loads             %3d
  stores            %3d
total               %3d`, iF, load, store, total)
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
  jumps             %3d
  jump-and-links    %3d
  taken branches    %3d
  untaken branches  %3d
total               %3d`, j, jal, taken, untaken, total)
	fmt.Println()
}

func (m *Machine) PrintBehavorialSimulation() {
	fmt.Println("\n" + `simple MIPS-like machine with instruction pairing
  (all values are shown in hexadecimal)`)
	fmt.Println()
}

func (m *Machine) PrintInstructionPairing() {
	pipe := m.pipeline
	var perc float64
	if pipe.issueCycle != 0 {
		perc = (float64(pipe.doubleIssue) / float64(pipe.issueCycle)) * 100.0
	}

	fmt.Printf("\ninstruction pairing counts (includes hlt instruction)\n")
	fmt.Printf("  issue cycles     %4d\n", pipe.issueCycle)
	fmt.Printf("  double issues    %4d ( %.1f percent of issue cycles)\n", pipe.doubleIssue, perc)
	fmt.Printf("  control stops    %4d\n", pipe.controlStop)
	fmt.Printf("  structural stops %4d (%d of which would also stop on a data dep.)\n", pipe.structuralStop, pipe.strData)
	fmt.Printf("  data dep. stops  %4d\n", pipe.dataDepStop)
}
