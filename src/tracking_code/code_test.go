package tracking_code

import (
	"testing"
)

func TestGenerateCode(t *testing.T) {
	for i := 0; i < 10000; i++ {
		assertUint64(uint64(i), GetId(GenerateCode(uint64(i))), t)
	}
	assertString("wxyz0", GenerateCode(0), t)
	assertString("wxyz9", GenerateCode(9), t)
	assertString("wxyza", GenerateCode(10), t)
	assertString("wxyzv", GenerateCode(31), t)
	assertString("wxy10", GenerateCode(32), t)
	assertString("wxy1a", GenerateCode(42), t)
	assertString("wxy1v", GenerateCode(63), t)
	assertString("wx10c", GenerateCode(32*32+12), t)
	assertString("w100h", GenerateCode(32*32*32+17), t)
	assertString("1000n", GenerateCode(32*32*32*32+23), t)
	assertString("10000r", GenerateCode(32*32*32*32*32+27), t)
	assertString("100000b", GenerateCode(32*32*32*32*32*32+11), t)
}

func assertUint64(expect uint64, actual uint64, t *testing.T) {
	if expect != actual {
		t.Errorf("Failed. Got %d, expected %d.", actual, expect)
	}
}

func assertString(expect string, actual string, t *testing.T) {
	if expect != actual {
		t.Errorf("Failed. Got %s, expected %s.", actual, expect)
	}
}
