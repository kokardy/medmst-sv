package main

import (
	"strconv"
)

const (
	GS1TYPE_ITEM = "01"
)

type Code interface {
	CheckDigit() (string, string)
	CheckDigitOk() bool
}

func CalcCheckDigit(code string, offset, length int) (string, string) {
	count := 0
	for i := offset; i < length-1; i++ {
		d, _ := strconv.Atoi(string(code[i]))
		//右から数えて偶数位置は3倍
		if (i+length)%2 == 0 {
			d *= 3
		}
		count += d
	}
	count %= 10
	digit := 10 - count
	digit %= 10
	return strconv.Itoa(digit), string(code[length-1])
}

type JAN string

func (j JAN) CheckDigit() (string, string) {
	return CalcCheckDigit(string(j), 0, 13)
}

func (j JAN) CheckDigitOK() bool {
	if d1, d2 := j.CheckDigit(); d1 == d2 {
		return true
	}
	return false
}

type GS1 string

func (g GS1) CheckDigit() (string, string) {
	return CalcCheckDigit(string(g), 2, 16)
}
func (g GS1) CheckDigitOK() bool {
	if d1, d2 := g.CheckDigit(); d1 == d2 {
		return true
	}
	return false
}

func (g GS1) ToJAN() JAN {
	pre_jan := JAN(string(g[3:16]))
	check_digit, _ := pre_jan.CheckDigit()
	return JAN(string(g[3:15]) + check_digit)
}
