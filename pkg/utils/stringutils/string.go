package stringutils

import (
	"strconv"
	"strings"
	"unsafe"
)

// S 字符串类型转换
type S string

func (s S) String() string {
	return string(s)
}

// Bytes 转换为[]byte
func (s S) Bytes() []byte {
	return []byte(s)
}

// Bool 转换为bool
func (s S) Bool() (bool, error) {
	b, err := strconv.ParseBool(s.String())
	if err != nil {
		return false, err
	}
	return b, nil
}

// DefaultBool 转换为bool，如果出现错误则使用默认值
func (s S) DefaultBool(defaultVal bool) bool {
	b, err := s.Bool()
	if err != nil {
		return defaultVal
	}
	return b
}

// Int64 转换为int64
func (s S) Int64() (int64, error) {
	i, err := strconv.ParseInt(s.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// DefaultInt64 转换为int64，如果出现错误则使用默认值
func (s S) DefaultInt64(defaultVal int64) int64 {
	i, err := s.Int64()
	if err != nil {
		return defaultVal
	}
	return i
}

// Int 转换为int
func (s S) Int() (int, error) {
	i, err := s.Int64()
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// DefaultInt 转换为int，如果出现错误则使用默认值
func (s S) DefaultInt(defaultVal int) int {
	i, err := s.Int()
	if err != nil {
		return defaultVal
	}
	return i
}

// Uint64 转换为uint64
func (s S) Uint64() (uint64, error) {
	i, err := strconv.ParseUint(s.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// DefaultUint64 转换为uint64，如果出现错误则使用默认值
func (s S) DefaultUint64(defaultVal uint64) uint64 {
	i, err := s.Uint64()
	if err != nil {
		return defaultVal
	}
	return i
}

// Uint 转换为uint
func (s S) Uint() (uint, error) {
	i, err := s.Uint64()
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}

// DefaultUint 转换为uint，如果出现错误则使用默认值
func (s S) DefaultUint(defaultVal uint) uint {
	i, err := s.Uint()
	if err != nil {
		return defaultVal
	}
	return uint(i)
}

// Float64 转换为float64
func (s S) Float64() (float64, error) {
	f, err := strconv.ParseFloat(s.String(), 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

// DefaultFloat64 转换为float64，如果出现错误则使用默认值
func (s S) DefaultFloat64(defaultVal float64) float64 {
	f, err := s.Float64()
	if err != nil {
		return defaultVal
	}
	return f
}

// Float32 转换为float32
func (s S) Float32() (float32, error) {
	f, err := s.Float64()
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

// DefaultFloat32 转换为float32，如果出现错误则使用默认值
func (s S) DefaultFloat32(defaultVal float32) float32 {
	f, err := s.Float32()
	if err != nil {
		return defaultVal
	}
	return f
}

// 转义json特殊字符
func EscapeJson(s string) string {
	specialChars := []string{"\\", "\b", "\f", "\n", "\r", "\t", "\""}
	replaceChars := []string{"\\\\", "\\b", "\\f", "\\n", "\\r", "\\t", "\\\""}

	return ReplaceStrings(s, specialChars, replaceChars)
}

func ReplaceStrings(s string, old []string, replace []string) string {
	if s == "" {
		return s
	}
	if len(old) != len(replace) {
		return s
	}

	for i, v := range old {
		s = strings.Replace(s, v, replace[i], 1000)
	}

	return s
}

// string []byte zero copy
func String2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func Bytes2string(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
