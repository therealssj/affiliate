package duotricenary

var encodeTable = []byte("0123456789abcdefghijklmnopqrstuv")

func Encode(num uint64) string {
	buf := make([]byte, 0, 16)
	for {
		buf = append(buf, encodeTable[num&31])
		num >>= 5
		if num == 0 {
			break
		}
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

func DecodeUint64(str string) uint64 {
	if len(str) == 0 {
		panic("blank string")
	}
	slice := []byte(str)
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	var res uint64 = 0
	for i := 0; i < len(slice); i++ {
		b := slice[i]
		if b < 58 {
			if b < 48 {
				panic("参数str中" + string(b) + "为非法字符")
			}
			res += uint64(b-48) << uint(5*i)
		} else if b > 96 {
			if b > 118 {
				panic("参数str中" + string(b) + "为非法字符")
			}
			res += uint64(b-87) << uint(5*i)
		} else {
			panic("参数str中" + string(b) + "为非法字符")
		}
	}
	return res
}
