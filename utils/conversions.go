package utils

import (
	"strconv"
	"unsafe"
)

func StrToI64(str string) int64 {
	if len(str) <= 0 {
		return 0
	}
	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return res
}

func I64ToStr(integer int64) string {
	return strconv.FormatInt(integer, 10)
}

func IToStr(integer int) string {
	return strconv.Itoa(integer)
}

func StrToI(str string) int {
	if len(str) <= 0 {
		return 0
	}
	res, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return res
}

// https://josestg.medium.com/140x-faster-string-to-byte-and-byte-to-string-conversions-with-zero-allocation-in-go-200b4d7105fc
func BytesToString(b []byte) string {
	// Ignore if your IDE shows an error here; it's a false positive.
	p := unsafe.SliceData(b)
	return unsafe.String(p, len(b))
}
func StringToBytes(s string) []byte {
	p := unsafe.StringData(s)
	b := unsafe.Slice(p, len(s))
	return b
}
