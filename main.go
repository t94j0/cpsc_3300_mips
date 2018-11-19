package main

import (
	machine "github.com/t94j0/cpsc_3300_mips/machine"
)

func main() {
	mac := machine.NewMachine()
	if err := mac.LoadFromStdin(); err != nil {
		panic(err)
	}
	mac.PrintMemory()
	mac.PrintBehavorialSimulation()
	mac.Execute()
	mac.PrintInstructionClassCounts()
	mac.PrintMemoryAccessCounts()
	mac.PrintTransferControlCounts()
	mac.PrintInstructionPairing()
}
