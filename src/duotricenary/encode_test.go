package duotricenary

import (
	"testing"
)

func TestEncode(t *testing.T) {
	for i := 0; i < 10000; i++ {
		assertUint64(uint64(i), DecodeUint64(Encode(uint64(i))), t)
	}
	assertString("0", Encode(0), t)
	assertString("1", Encode(1), t)
	assertString("2", Encode(2), t)
	assertString("3", Encode(3), t)
	assertString("4", Encode(4), t)
	assertString("5", Encode(5), t)
	assertString("6", Encode(6), t)
	assertString("7", Encode(7), t)
	assertString("8", Encode(8), t)
	assertString("9", Encode(9), t)
	assertString("a", Encode(10), t)
	assertString("b", Encode(11), t)
	assertString("c", Encode(12), t)
	assertString("d", Encode(13), t)
	assertString("e", Encode(14), t)
	assertString("f", Encode(15), t)
	assertString("g", Encode(16), t)
	assertString("h", Encode(17), t)
	assertString("i", Encode(18), t)
	assertString("j", Encode(19), t)
	assertString("k", Encode(20), t)
	assertString("l", Encode(21), t)
	assertString("m", Encode(22), t)
	assertString("n", Encode(23), t)
	assertString("o", Encode(24), t)
	assertString("p", Encode(25), t)
	assertString("q", Encode(26), t)
	assertString("r", Encode(27), t)
	assertString("s", Encode(28), t)
	assertString("t", Encode(29), t)
	assertString("u", Encode(30), t)
	assertString("v", Encode(31), t)
	assertString("10", Encode(32), t)
	assertString("1a", Encode(42), t)
	assertString("1v", Encode(63), t)
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
