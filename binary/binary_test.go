package binary

import (
	"fmt"
	"testing"
)

func ExampleGetIFormat() {
	s, t, i := GetIFormat(0xac030003)
	fmt.Printf("%d %d %d", s, t, i)
	// output:
	// 0 3 3
}

func TestGetIFormat_lw(test *testing.T) {
	s, t, i := GetIFormat(0x8c010001)
	if s != 0x0 || i != 0x1 || t != 0x1 {
		test.Logf("results: s = %d, t = %d, i = %d", s, t, i)
		test.Error("failed GetIFormat")
	}
}

func TestGetJFormat(t *testing.T) {
	loc := GetJFormat( 0x08000004)
	if loc != 0x4 {
		t.Logf("result: loc = %d\n", loc)
		t.Error("failed GetJFormat")
	}
}

func TestGetIFormat_sw(test *testing.T) {
	s, t, i := GetIFormat(0xac030003)
	if s != 0x0 || i != 0x3 || t != 0x3 {
		test.Logf("results: s = %d, t = %d, i = %d", s, t, i)
		test.Error("failed GetIFormat")
	}
}

func TestGetRFormat(test *testing.T) {
	s, t, d, shift, f := GetRFormat(0x00432023)
	if f != 0x23 || s != 0x2  || t != 0x3 || d != 0x4 || shift != 0x0 {
		test.Logf("results: s = %x, t = %x, d = %x, shift = %x, funct = 0x%x",
			s, t, d, shift, f)
		test.Log("expected: s = 4, t = 3, d = 2, shift = 0, " + "f = 0x23")
		test.Error("failed GetRFormat")
	}
}
