package main

import (
	machine "github.com/t94j0/cpsc_3300_mips/machine"
)

func main() {
	machine := machine.Machine{}
	if err := machine.LoadFromStdin(); err != nil {
		panic(err)
	}
	machine.PrintMemory()
	machine.PrintBehavorialSimulation()
	machine.Execute()
	machine.PrintMemory()
	machine.PrintInstructionClassCounts()
	machine.PrintMemoryAccessCounts()
	machine.PrintTransferControlCounts()
}
