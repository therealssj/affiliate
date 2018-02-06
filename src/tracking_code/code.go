package tracking_code

import (
	"github.com/spolabs/affiliate/src/duotricenary"
)

func GenerateCode(id uint64) string {
	str := duotricenary.Encode(id)
	switch len(str) {
	case 1:
		return "wxyz" + str
	case 2:
		return "wxy" + str
	case 3:
		return "wx" + str
	case 4:
		return "w" + str
	default:
		return str
	}
}

func GetId(str string) (id uint64) {
	defer func() {
		if err := recover(); err != nil {
			id = 0
		}
	}()
	if len(str) < 5 {
		return 0
	}
	if str[0] == 'w' {
		if str[1] == 'x' {
			if str[2] == 'y' {
				if str[3] == 'z' {
					str = str[4:]
				} else {
					str = str[3:]
				}
			} else {
				str = str[2:]
			}
		} else {
			str = str[1:]
		}
	}
	return duotricenary.DecodeUint64(str)
}
