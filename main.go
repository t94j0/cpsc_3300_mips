package main

import (
	machine "./machine"
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
