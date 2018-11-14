package machine

import (
	"fmt"

	"github.com/spf13/cast"

	binary "github.com/t94j0/go-mips-instruction-format"
)

type InstructionFunc func(m *Machine, inst uint32)

var opcodeInstructions = map[uint16]InstructionFunc{
	0x00: zeroOpcode, 0x02: j, 0x03: jal, 0x04: beq, 0x05: bne, 0x06: blez,
	0x07: bgtz, 0x09: addiu, 0x0a: slti, 0x0f: lui, 0x0e: xori, 0x1c: mul,
	0x23: lw, 0x2b: sw,
}
var zeroInstructions = map[uint16]InstructionFunc{
	0x21: addu, 0x24: and, 0x09: jalr, 0x08: jr, 0x27: nor, 0x25: or, 0x00: sll,
	0x03: sra, 0x02: srl, 0x23: subu, 0x26: xor,
}

func printInstruction(instruction string, register uint16, value, pc uint32) {
	f := "%03x: %-5s - register r[%d] now contains 0x%08x\n"
	fmt.Printf(f, pc, instruction, register, value)
}
func printJump(ir uint32, instruction string, register uint32) {
	f := "%03x: %-5s - jump to 0x%08X\n"
	fmt.Printf(f, ir, instruction, register)
}
func printBranch(pc uint32, instruction string, register uint32) {
	f := "%03x: %-5s - branch taken to 0x%08x\n"
	fmt.Printf(f, pc, instruction, register)
}

func zeroOpcode(m *Machine, inst uint32) {
	if inst == 0x0 {
		m.memoryAccess.instFetch--
		fmt.Printf("%03x: hlt\n", m.ir)
		m.halt = true
		return
	}
	funct := uint16(binary.GetFunct(inst))
	if _, ok := zeroInstructions[funct]; !ok {
		panic(fmt.Sprintf("opcode does not exist: %08x", inst))
	}
	zeroInstructions[funct](m, inst)
}
func j(m *Machine, inst uint32) {
	m.transferControl.jump++
	dst := binary.GetJFormat(inst)
	printJump(m.ir, "j", dst)
	m.pc = dst
}
func jal(m *Machine, inst uint32) {
	m.transferControl.jumpLink++
	dst := binary.GetJFormat(inst)
	printJump(m.ir, "jal", dst)
	m.registers[30] = m.pc
	m.pc = dst
}
func bne(m *Machine, inst uint32) {
	sr, tr, immu := binary.GetIFormat(inst)
	s, t := m.memory[sr], m.memory[tr]
	if s != t {
		m.transferControl.takenBranch++
		loc := int16(m.ir) + int16(immu)
		printBranch(m.pc, "bne", uint32(loc))
		m.pc = uint32(loc)
	} else {
		m.transferControl.untakenBranch++
	}
}
func blez(m *Machine, inst uint32) {
	s, _, immu := binary.GetIFormat(inst)
	val := int32(m.registers[s])

	if val <= 0 {
		m.transferControl.takenBranch++
		loc := int16(m.ir) + int16(immu)
		printBranch(m.pc, "blez", uint32(loc))
		m.pc = uint32(loc)
	} else {
		m.transferControl.untakenBranch++
	}
}
func bgtz(m *Machine, inst uint32) {
	s, _, immu := binary.GetIFormat(inst)
	valu := m.registers[s]
	val := int32(valu)

	if val > 0 {
		m.transferControl.takenBranch++
		loc := int16(m.ir) + int16(immu)
		printBranch(m.pc, "bgtz", uint32(loc))
		m.pc = uint32(loc)
	} else {
		m.transferControl.untakenBranch++
	}
}
func addiu(m *Machine, inst uint32) {
	m.instructionClass.alu++
	s, t, imm := binary.GetIFormat(inst)
	immCast := cast.ToInt16(imm)

	sum := uint32(int32(m.registers[s]) + int32(immCast))

	printInstruction("addiu", t, sum, m.ir)
	m.registers[t] = sum
}
func slti(m *Machine, inst uint32) {
	m.instructionClass.alu++
	su, tu, immu := binary.GetIFormat(inst)
	s, imm := int32(m.registers[su]), int32(immu)
	if s < imm {
		printInstruction("slti", tu, 1, m.ir)
		m.registers[tu] = 1
	} else {
		printInstruction("slti", tu, 0, m.ir)
		m.registers[tu] = 0
	}
}
func lui(m *Machine, inst uint32) {
	m.instructionClass.alu++
	_, tu, imm := binary.GetIFormat(inst)
	val := cast.ToUint32(imm) << 16
	printInstruction("lui", tu, uint32(val), m.ir)
	m.registers[tu] = uint32(val)
}
func xori(m *Machine, inst uint32) {
	m.instructionClass.alu++
	su, tu, imm := binary.GetIFormat(inst)
	s := m.registers[su]
	xoriVal := s ^ uint32(imm)
	printInstruction("xori", tu, xoriVal, m.ir)
	m.registers[tu] = xoriVal
}
func mul(m *Machine, inst uint32) {
	m.instructionClass.alu++
	su, tu, du, _, funct := binary.GetRFormat(inst)
	if funct != 0x2 {
		panic(fmt.Sprintf("%x is not a valid instruction", inst))
	}
	s, t := m.registers[su], m.registers[tu]
	prod := s * t
	printInstruction("mul", du, prod, m.ir)
	m.registers[du] = prod
}
func lw(m *Machine, inst uint32) {
	m.memoryAccess.load++
	s, t, imm := binary.GetIFormat(inst)
	value := m.memory[s+imm]
	fmt.Printf("%03x: %-5s - register r[%d] now contains 0x%08x\n", m.ir,
		"lw", t, value)
	m.registers[t] = m.memory[s+imm]
}
func sw(m *Machine, inst uint32) {
	m.memoryAccess.store++
	s, t, imm := binary.GetIFormat(inst)
	fmt.Printf("%03x: %-5s - register r[%d] value stored in memory\n", m.ir, "sw", t)
	m.memory[s+imm] = m.registers[t]

}
func beq(m *Machine, inst uint32) {
	sr, tr, immu := binary.GetIFormat(inst)
	s, t := m.memory[sr], m.memory[tr]
	if s == t {
		m.transferControl.takenBranch++
		loc := int16(m.ir) + int16(immu)
		printBranch(m.pc, "beq", uint32(loc))
		m.pc = uint32(loc)
	} else {
		m.transferControl.untakenBranch++
	}
}
func addu(m *Machine, inst uint32) {
	m.instructionClass.alu++
	s, t, d, _, _ := binary.GetRFormat(inst)
	sum := m.registers[s] + m.registers[t]
	printInstruction("addu", d, sum, m.ir)
	m.registers[d] = sum
}
func and(m *Machine, inst uint32) {
	m.instructionClass.alu++
	su, tu, du, _, _ := binary.GetRFormat(inst)
	s, t := m.registers[su], m.registers[tu]
	andVal := s & t
	printInstruction("and", du, andVal, m.ir)
	m.registers[du] = andVal
}
func jalr(m *Machine, inst uint32) {
	m.transferControl.jumpLink++
	s, _, d, _, _ := binary.GetRFormat(inst)
	m.registers[d] = m.pc
	printJump(m.ir, "jalr", m.registers[s])
	m.pc = m.registers[s]
}
func jr(m *Machine, inst uint32) {
	m.transferControl.jump++
	s, _, _, _, _ := binary.GetRFormat(inst)
	printJump(m.ir, "jr", m.registers[s])
	m.pc = m.registers[s]
}
func nor(m *Machine, inst uint32) {
	m.instructionClass.alu++
	su, tu, du, _, _ := binary.GetRFormat(inst)
	s, t := m.registers[su], m.registers[tu]
	norVal := ^(s | t)
	printInstruction("nor", du, norVal, m.ir)
	m.registers[du] = norVal
}
func or(m *Machine, inst uint32) {
	m.instructionClass.alu++
	su, tu, du, _, _ := binary.GetRFormat(inst)
	s, t := m.registers[su], m.registers[tu]
	orVal := s | t
	printInstruction("or", du, orVal, m.ir)
	m.registers[du] = orVal
}
func sll(m *Machine, inst uint32) {
	m.instructionClass.alu++
	_, tu, du, hu, _ := binary.GetRFormat(inst)
	t := m.registers[tu]
	sllVal := t << hu
	printInstruction("sll", du, sllVal, m.ir)
	m.registers[du] = sllVal
}
func sra(m *Machine, inst uint32) {
	m.instructionClass.alu++
	_, tu, du, hu, _ := binary.GetRFormat(inst)
	sraVal := m.registers[tu] >> m.registers[hu]
	printInstruction("sll", du, sraVal, m.ir)
	m.registers[du] = sraVal
}
func srl(m *Machine, inst uint32) {
	m.instructionClass.alu++
	_, tu, du, h, _ := binary.GetRFormat(inst)
	t := m.registers[tu]
	srlVal := t >> h
	printInstruction("srl", du, srlVal, m.ir)
	m.registers[du] = srlVal
}
func subu(m *Machine, inst uint32) {
	m.instructionClass.alu++
	s, t, d, _, _ := binary.GetRFormat(inst)
	diff := m.registers[s] - m.registers[t]
	printInstruction("subu", d, diff, m.ir)
	m.registers[d] = diff
}
func xor(m *Machine, inst uint32) {
	m.instructionClass.alu++
	su, tu, du, _, _ := binary.GetRFormat(inst)
	s, t := m.registers[su], m.registers[tu]
	xorVal := s ^ t
	printInstruction("xor", du, xorVal, m.ir)
	m.registers[du] = xorVal
}
