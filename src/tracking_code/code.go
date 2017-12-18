package tracking_code

import (
	"github.com/spaco/affiliate/src/duotricenary"
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
	return duotricenary.DecodeUint64(str)
}
