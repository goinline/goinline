// Copyright 2021 冯立强 mr.fengliqiang@gmail.com.  All rights reserved.

package goinline

//IsAlpha
//go:nosplit
func IsAlpha(ch uint8) bool {
	t := (ch | ('A' ^ 'a')) - 'a'
	return t <= 'z'-'a'
}

//IsWhite
//go:nosplit
func IsWhite(ch uint8) bool {
	return ch != 0 && ch <= 0x20
}

//ToLower
//go:nosplit
func ToLower(ch uint8) uint8 {

	if t := ch | ('A' ^ 'a'); t-'a' <= 'z'-'a' {
		return t
	}
	return ch
}

//ToUpper
//go:nosplit
func ToUpper(ch uint8) uint8 {

	if t := ch &^ ('A' ^ 'a'); t-'A' <= 'z'-'a' {
		return t
	}
	return ch
}

//go:nosplit
func min(a, b int) int {

	if a < b {
		return a
	}
	return b
}

//StrIsdigit : check digit string
//go:nosplit
func StrIsdigit(value string) bool {
	for i := 0; i < len(value); i++ {
		if value[i] < '0' || '9' < value[i] {
			return false
		}
	}
	return len(value) > 0
}

//StrCasecmp : return int; a < b => -1; a > b => 1; a == b => 0
//go:nosplit
func StrCasecmp(a, b string) int {
	size := min(len(a), len(b))
	for i := 0; i < size; i++ {
		lowerA := ToLower(a[i])
		lowerB := ToLower(b[i])
		if lowerA == lowerB {
			continue
		}
		if lowerB < lowerA {
			return 1
		}
		return -1
	}
	if len(a) == len(b) {
		return 0
	}
	if len(a) < len(b) {
		return -1
	}
	return 1
}

//StrCmp : return int; a < b => -1; a > b => 1; a == b => 0
//go:nosplit
func StrCmp(a, b string) int {
	size := min(len(a), len(b))
	for i := 0; i < size; i++ {
		if a[i] != b[i] {
			if b[i] < a[i] {
				return 1
			}
			return -1
		}
	}
	if len(a) == len(b) {
		return 0
	}
	if len(a) < len(b) {
		return -1
	}
	return 1
}

//StrSubstr
//go:nosplit
func StrSubstr(str string, begin, size int) string {
	begin = min(begin, len(str))
	return str[begin:min(len(str), begin+size)]
}

//StrRsubstr
//go:nosplit
func StrRsubstr(str string, size int) string {
	size = min(len(str), size)
	return str[len(str)-size:]
}

//StrChr
//go:nosplit
func StrChr(str string, ch uint8) int {
	for offset := 0; offset < len(str); offset++ {
		if str[offset] == ch {
			return offset
		}
	}
	return -1
}

//StrRchr : search charter from right
//go:nosplit
func StrRchr(str string, ch uint8) int {
	for offset := len(str); offset > 0; offset-- {
		if str[offset-1] == ch {
			return offset - 1
		}
	}
	return -1
}

//StrStr : search target from source begin, return first index
//go:nosplit
func StrStr(source, target string) int {
	if len(target) == 0 {
		return 0
	}
	matchSize := len(source) - len(target)
	if matchSize < 0 {
		return -1
	}
	for off := 0; off <= matchSize; off++ {
		i := 0
		for i < len(target) && source[off+i] == target[i] {
			i++
		}
		if i == len(target) {
			return off
		}
	}
	return -1
}

//StrRstr : search target from source end, return first index
//go:nosplit
func StrRstr(source, target string) int {
	if len(target) == 0 {
		return len(source)
	}
	matchBegin := len(source) - len(target)
	if matchBegin < 0 {
		return -1
	}
	for off := matchBegin; off >= 0; off-- {
		i := 0
		for i < len(target) && source[off+i] == target[i] {
			i++
		}
		if i == len(target) {
			return off
		}
	}
	return -1
}

//StrSlice : split string by issep
//go:nosplit
func StrSlice(source string, issep func(byte) bool) []string {
	result := []string{}
	begin := -1
	for i := 0; i < len(source); i++ {
		if issep(source[i]) {
			if begin > -1 {
				result = append(result, source[begin:i])
				begin = -1
			}
		} else {
			if begin == -1 {
				begin = i
			}
		}
	}
	if begin > -1 && begin < len(source) {
		result = append(result, source[begin:])
	}
	return result
}

//StrSplit : split string by issep, handler handle the str, if handler return true, then match end.
//go:nosplit
func StrSplit(source string, issep func(byte) bool, handler func(string) bool) {
	if handler == nil || issep == nil {
		return
	}
	begin := -1
	for i := 0; i < len(source); i++ {
		if issep(source[i]) {
			if begin > -1 {
				if handler(source[begin:i]) {
					return
				}
				begin = -1
			}
		} else {
			if begin == -1 {
				begin = i
			}
		}
	}
	if begin > -1 && begin < len(source) {
		handler(source[begin:])
	}
	return
}

//StrSplitS : split string by sep, handler handle the str, if handler return true, then match end.
//go:nosplit
func StrSplitS(source, sep string, handler func(string) bool) {
	if handler == nil || len(source) == 0 {
		return
	}
	if len(sep) == 0 {
		handler(source)
		return
	}
	for {
		if id := StrStr(source, sep); id != -1 {
			if handler(source[0:id]) {
				break
			}
			source = source[id+len(sep):]
		} else {
			handler(source)
			break
		}
	}
}

//StrPaire : split paire by issep, parse like: name = value  or name : value
//go:nosplit
func StrPaire(source string, issep func(byte) bool) []string {
	if len(source) == 0 {
		return []string{}
	}
	for i := 0; i < len(source); i++ {
		if issep(source[i]) {
			return []string{source[0:i], source[i+1:]}
		}
	}
	return []string{source}
}

//go:nosplit
func StrPaireS(source, sep string, handler func(name, value string)) {
	if handler == nil || len(source) == 0 {
		return
	}
	if len(sep) == 0 {
		handler(source, "")
		return
	}
	if id := StrStr(source, sep); id != -1 {
		handler(source[0:id], source[id+len(sep):])
	} else {
		handler(source, "")
	}
}

//StrTrim : remove string left and string right white charter
//go:nosplit
func StrTrim(source string, issep func(byte) bool) string {
	begin := 0
	for begin < len(source) && issep(source[begin]) {
		begin++
	}
	source = source[begin:]
	end := len(source)
	for end > 0 && issep(source[end-1]) {
		end--
	}
	return source[0:end]
}
