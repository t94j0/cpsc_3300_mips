package machine

import (
	"fmt"
	"math"

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

func printInstruction(instruction string, ir uint32) {
	f := "%03x: %-6s"
	fmt.Printf(f, ir, instruction)
}

func zeroOpcode(m *Machine, inst uint32) {
	if inst == 0x0 {
		m.memoryAccess.instFetch--
		printInstruction("hlt", m.ir)
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
	printInstruction("j", m.ir)
	m.pc = dst
}
func jal(m *Machine, inst uint32) {
	m.transferControl.jumpLink++
	dst := binary.GetJFormat(inst)
	printInstruction("jal", m.ir)
	m.registers[31] = m.pc
	m.pc = dst
}
func bne(m *Machine, inst uint32) {
	sr, tr, immu := binary.GetIFormat(inst)
	s, t := m.registers[sr], m.registers[tr]

	printInstruction("bne", m.ir)
	if s != t {
		m.transferControl.takenBranch++
		loc := int16(m.pc) + cast.ToInt16(immu)

		m.pc = uint32(loc)
	} else {
		m.transferControl.untakenBranch++
	}
}
func blez(m *Machine, inst uint32) {
	s, _, immu := binary.GetIFormat(inst)
	val := int32(m.registers[s])
	printInstruction("blez", m.ir)
	if val <= 0 {
		m.transferControl.takenBranch++
		loc := int16(m.ir) + int16(immu)
		m.pc = uint32(loc)
	} else {
		m.transferControl.untakenBranch++

	}
}
func bgtz(m *Machine, inst uint32) {
	s, _, immu := binary.GetIFormat(inst)
	valu := m.registers[s]
	val := cast.ToInt32(valu)
	printInstruction("bgtz", m.ir)
	if val > 0 {
		m.transferControl.takenBranch++
		loc := int16(m.pc) + int16(immu)
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

	printInstruction("addiu", m.ir)
	m.writeTo = int(t)
	m.registers[t] = sum
}
func slti(m *Machine, inst uint32) {
	m.instructionClass.alu++
	su, tu, immu := binary.GetIFormat(inst)
	s, imm := int32(m.registers[su]), int32(immu)
	printInstruction("slti", m.ir)
	if s < imm {

		m.registers[tu] = 1
	} else {
		m.registers[tu] = 0
	}
	m.writeTo = int(tu)
}
func lui(m *Machine, inst uint32) {
	m.instructionClass.alu++
	_, tu, imm := binary.GetIFormat(inst)
	val := cast.ToUint32(imm) << 16
	printInstruction("lui", m.ir)
	m.registers[tu] = uint32(val)
	m.writeTo = int(tu)
}
func xori(m *Machine, inst uint32) {
	m.instructionClass.alu++
	su, tu, imm := binary.GetIFormat(inst)
	s := m.registers[su]
	xoriVal := s ^ uint32(imm)
	printInstruction("xori", m.ir)
	m.registers[tu] = xoriVal
	m.writeTo = int(tu)
}
func mul(m *Machine, inst uint32) {
	m.instructionClass.alu++
	su, tu, du, _, funct := binary.GetRFormat(inst)
	if funct != 0x2 {
		panic(fmt.Sprintf("%x is not a valid instruction", inst))
	}
	s, t := m.registers[su], m.registers[tu]
	prod := s * t
	printInstruction("mul", m.ir)
	m.registers[du] = prod
	m.writeTo = int(du)
}
func lw(m *Machine, inst uint32) {
	m.memoryAccess.load++
	s, t, imm := binary.GetIFormat(inst)
	printInstruction("lw", m.ir)
	m.registers[t] = m.memory[s+imm]
	m.writeTo = int(t)
}
func sw(m *Machine, inst uint32) {
	m.memoryAccess.store++
	s, t, imm := binary.GetIFormat(inst)
	printInstruction("sw", m.ir)
	m.memory[s+imm] = uint32(uint16(m.registers[t]))
}
func beq(m *Machine, inst uint32) {
	sr, tr, immu := binary.GetIFormat(inst)
	s, t := m.registers[sr], m.registers[tr]
	printInstruction("beq", m.ir)
	if s == t {
		m.transferControl.takenBranch++
		loc := int16(m.ir) + int16(immu)

		m.pc = uint32(loc)
	} else {
		m.transferControl.untakenBranch++
	}
}
func addu(m *Machine, inst uint32) {
	m.instructionClass.alu++
	s, t, d, _, _ := binary.GetRFormat(inst)
	sum := m.registers[s] + m.registers[t]
	printInstruction("addu", m.ir)
	m.registers[d] = sum
	m.writeTo = int(d)
}
func and(m *Machine, inst uint32) {
	m.instructionClass.alu++
	su, tu, du, _, _ := binary.GetRFormat(inst)
	s, t := m.registers[su], m.registers[tu]
	andVal := s & t
	printInstruction("and", m.ir)
	m.registers[du] = andVal
	m.writeTo = int(du)
}
func jalr(m *Machine, inst uint32) {
	m.transferControl.jumpLink++
	s, _, d, _, _ := binary.GetRFormat(inst)
	m.registers[d] = m.pc
	printInstruction("jalr", m.ir)
	m.pc = m.registers[s]
	m.writeTo = int(d)
}
func jr(m *Machine, inst uint32) {
	m.transferControl.jump++
	s, _, _, _, _ := binary.GetRFormat(inst)
	printInstruction("jr", m.ir)
	m.pc = m.registers[s]
}
func nor(m *Machine, inst uint32) {
	m.instructionClass.alu++
	su, tu, du, _, _ := binary.GetRFormat(inst)
	s, t := m.registers[su], m.registers[tu]
	norVal := ^(s | t)
	printInstruction("nor", m.ir)
	m.registers[du] = norVal
	m.writeTo = int(du)
}
func or(m *Machine, inst uint32) {
	m.instructionClass.alu++
	su, tu, du, _, _ := binary.GetRFormat(inst)
	s, t := m.registers[su], m.registers[tu]
	orVal := s | t
	printInstruction("or", m.ir)
	m.registers[du] = orVal
	m.writeTo = int(du)
}
func sll(m *Machine, inst uint32) {
	m.instructionClass.alu++
	_, tu, du, hu, _ := binary.GetRFormat(inst)
	t := m.registers[tu]
	sllVal := t << hu
	printInstruction("sll", m.ir)
	m.registers[du] = sllVal
	m.writeTo = int(du)
}
func sra(m *Machine, inst uint32) {
	m.instructionClass.alu++
	_, tu, du, hu, _ := binary.GetRFormat(inst)
	t := m.registers[tu]
	mask := uint32(math.Pow(2, float64(hu))-1) << (32 - hu)

	sraVal := (t >> hu) | mask

	printInstruction("sra", m.ir)
	m.registers[du] = sraVal
	m.writeTo = int(du)
}
func srl(m *Machine, inst uint32) {
	m.instructionClass.alu++
	_, tu, du, h, _ := binary.GetRFormat(inst)
	t := m.registers[tu]
	srlVal := t >> h
	printInstruction("srl", m.ir)
	m.registers[du] = srlVal
	m.writeTo = int(du)
}
func subu(m *Machine, inst uint32) {
	m.instructionClass.alu++
	s, t, d, _, _ := binary.GetRFormat(inst)
	diff := m.registers[s] - m.registers[t]
	printInstruction("subu", m.ir)
	m.registers[d] = diff
	m.writeTo = int(d)
}
func xor(m *Machine, inst uint32) {
	m.instructionClass.alu++
	su, tu, du, _, _ := binary.GetRFormat(inst)
	s, t := m.registers[su], m.registers[tu]
	xorVal := s ^ t
	printInstruction("xor", m.ir)
	m.registers[du] = xorVal
	m.writeTo = int(du)
}
