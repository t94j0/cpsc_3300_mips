package machine

import (
	"fmt"

	"github.com/spf13/cast"
	binary "github.com/t94j0/go-mips-instruction-format"
)

type Pipeline struct {
	issueCycle     uint
	doubleIssue    uint
	controlStop    uint
	structuralStop uint
	dataDepStop    uint
	strData        uint

	m *Machine
}

func NewPipeline(m *Machine) *Pipeline {
	return &Pipeline{m: m}
}

func (p *Pipeline) Schedule() {
	op, funct, inst := p.m.getNextOp()
	p.m.runInstruction()
	if p.shouldRunSecond(op, funct, inst) {
		p.m.runInstruction()
		p.doubleIssue++
		fmt.Printf("  // -- double issue --")
	}
	p.issueCycle++
	p.flush()
}

func (p *Pipeline) shouldRunSecond(oldOp, oldFunct uint16, oldInst uint32) bool {
	printControl := func(s string) { fmt.Printf("%13s%s", " ", s) }
	dep := p.hasDataDep()

	if p.oneHalt(oldInst) {
		printControl("// control stop")
		p.controlStop++
		return false
	}

	if p.bothLS(oldOp, oldFunct) {
		printControl("// structural stop")
		if dep {
			fmt.Printf(" (also data dep.)")
			p.strData++
		}
		p.structuralStop++
		return false
	}

	if p.bothMultiply(oldOp, oldFunct) {
		printControl("// structural stop")
		if dep {
			fmt.Printf(" (also data dep.)")
			p.strData++
		}
		p.structuralStop++
		return false
	}

	if p.firstBranch(oldOp, oldFunct) {
		printControl("// control stop")
		p.controlStop++
		return false
	}

	if dep {
		printControl("// data dependency stop")
		p.dataDepStop++
		return false
	}

	return true
}

func (p *Pipeline) hasDataDep() bool {
	if p.m.writeTo == -1 {
		return false
	}

	write := cast.ToUint16(p.m.writeTo)

	inst, funct, no := p.m.getNextOp()
	rs, rt, rd, shamt, _ := binary.GetRFormat(no)
	is, it, _ := binary.GetIFormat(no)

	switch inst {
	case 0x0:
		switch funct {
		case 0x21:
			return write == rs || write == rt || write == rd
		case 0x24:
			return write == rs || write == rt || write == rd
		case 0x09:
			return write == rs || write == rd
		case 0x08:
			return write == rs
		case 0x27:
			return write == rs || write == rt || write == rd
		case 0x25:
			return write == rs || write == rt || write == rd
		case 0x0:
			return write == shamt || write == rt || write == rd
		case 0x3:
			return write == shamt || write == rt || write == rd
		case 0x2:
			return write == shamt || write == rt || write == rd
		case 0x23:
			return write == rs || write == rt || write == rd
		case 0x26:
			return write == rs || write == rt || write == rd
		}
	case 0x9:
		return write == is || write == it
	case 0x4:
		return write == is || write == it
	case 0x7:
		return write == is
	case 0x6:
		return write == is
	case 0x5:
		return write == is || write == it
	case 0x3:
		return write == 31
	case 0xf:
		return write == it
	case 0x23:
		return write == is || write == it
	case 0x1c:
		return write == rs || write == rt || write == rd
	case 0x0a:
		return write == is || write == it
	case 0x2b:
		return write == is || write == it
	case 0x0e:
		return write == is || write == it
	}

	return false
}

func (p *Pipeline) bothLS(oldOp, oldFunct uint16) bool {
	op, _, _ := p.m.getNextOp()
	isLS := func(op uint16) bool {
		if op == 0x23 || op == 0x2b {
			return true
		}
		return false
	}

	return isLS(oldOp) && isLS(op)
}

func (p Pipeline) oneHalt(oldInst uint32) bool {
	return oldInst == 0x0
}

func (p *Pipeline) firstBranch(oldOp, oldInst uint16) bool {
	isBranch := func(mop, mfunct uint16) bool {
		// j, jal, beq, bne, blez, bgtz
		if mop == 0x02 || mop == 0x03 || mop == 0x05 || mop == 0x06 || mop == 0x07 {
			return true
		}

		// hlt, jalr, jr
		if mop == 0x0 {
			if mfunct == 0x09 || mfunct == 0x08 {
				return true
			}
		}

		return false
	}
	return isBranch(oldOp, oldInst)
}

func (p *Pipeline) bothMultiply(oldOp, oldFunct uint16) bool {
	op, _, _ := p.m.getNextOp()
	return oldOp == 0x1c && op == 0x1c
}

func (p *Pipeline) flush() {
	fmt.Printf("\n")
}
