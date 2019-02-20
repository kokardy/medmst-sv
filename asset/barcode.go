package main

import (
	"strconv"
)

const (
	//GS1TYPEITEM is GS1コードのアイテム接頭番号
	GS1TYPEITEM = "01"
)

//Code has checkdisit
type Code interface {
	CheckDigit() (string, string)
	CheckDigitOk() bool
}

//CalcCheckDigit calcs checkdigit
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

//JAN is JAN CODE
type JAN string

//CheckDigit calc checkdigit
func (j JAN) CheckDigit() (string, string) {
	return CalcCheckDigit(string(j), 0, 13)
}

//CheckDigitOK validates checkdigit
func (j JAN) CheckDigitOK() bool {
	if d1, d2 := j.CheckDigit(); d1 == d2 {
		return true
	}
	return false
}

//GS1 is GS1 code
type GS1 string

//CheckDigit calc checkdigit
func (g GS1) CheckDigit() (string, string) {
	return CalcCheckDigit(string(g), 2, 16)
}

//CheckDigitOK validates checkdigit
func (g GS1) CheckDigitOK() bool {
	if d1, d2 := g.CheckDigit(); d1 == d2 {
		return true
	}
	return false
}

//ToJAN converts GS1 to JAN
func (g GS1) ToJAN() JAN {
	preJan := JAN(string(g[3:16]))
	checkDigit, _ := preJan.CheckDigit()
	return JAN(string(g[3:15]) + checkDigit)
}
