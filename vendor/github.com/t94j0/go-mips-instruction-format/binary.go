package binary

import (
	"math"
)

// GetBits gets the bits between `start` index and `end` index
//
// 0b100000ssssstttttiiiiiiiiiiiiiiii; start = 0; end = 5
// return 0b100000
func getBits(target uint32, start, end uint8) (uint32) {
	shifted := target >> (31 - end)
	mask := uint32(math.Pow(2, float64(end-start+1))) - 1

	return  shifted & mask
}

// GetOperation gets the first 6 bits of a 32 bit binary string
func GetOperation(target uint32) uint32 {
	return getBits(target, 0, 5)
}

// GetFunct gets the last 6 bytes
func GetFunct(target uint32) uint32 {
	return getBits(target, 26, 31)
}

func GetRFormat(target uint32) (uint16, uint16, uint16, uint16, uint16) {
	s := getBits(target, 6, 10)
	t := getBits(target, 11, 15)
	d := getBits(target, 16, 20)
	shift := getBits(target, 21, 25)
	f := getBits(target, 26, 31)
	return uint16(s), uint16(t), uint16(d), uint16(shift), uint16(f)
}

func GetIFormat(target uint32) (uint16, uint16, uint16) {
	s := getBits(target, 6, 10)
	t := getBits(target, 11, 15)
	imm := getBits(target, 16, 31)
	return uint16(s), uint16(t), uint16(imm)
}

func GetJFormat(target uint32) (uint32) {
	return getBits(target, 6, 31)
}
