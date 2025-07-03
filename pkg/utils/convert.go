package utils

import "strconv"

func StrToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}

func StrToUint64(str string) uint64 {
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}
	return num
}

func StrToFloat64(str string) float64 {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return num
}

func StrToBool(str string) bool {
	num, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return num
}

func IntToStr(num int) string {
	return strconv.Itoa(num)
}

func Uint64ToStr(num uint64) string {
	return strconv.FormatUint(num, 10)
}
