package goeinstein

import (
	"fmt"
	"strconv"
	"strings"
)

func ToString(v interface{}) string { return fmt.Sprintf("%v", v) }

func ToLowerCase(s string) string { return strings.ToLower(s) }
func ToUpperCase(s string) string { return strings.ToUpper(s) }

func NumToStr(num int) string { return strconv.Itoa(num) }

func StrToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		panic(fmt.Sprintf("Invalid integer '%s': %v", str, err))
	}
	return i
}

func StrToDouble(str string) float32 {
	f, err := strconv.ParseFloat(str, 32) //nolint:gomnd
	if err != nil {
		panic(fmt.Sprintf("Invalid double '%s': %v", str, err))
	}
	return float32(f)
}
