package main

import (
	"fmt"
	"strconv"
	"strings"
)

var asciiEscapes = []byte{'\\', '"', 'b', 'f', 'n', 'r', 't'}
var binaryEscapes = []byte{'\\', '"', '\b', '\f', '\n', '\r', '\t'}

func escapeByte(b *strings.Builder, c byte) {
	for i, esc := range binaryEscapes {
		if esc == c {
			b.WriteByte('\\')
			b.WriteByte(asciiEscapes[i])
			return
		}
	}
	if c < 0x20 {
		b.WriteString(fmt.Sprintf("\\u%04x", c))
		return
	}
	b.WriteByte(c)
}

func DecorateJsonStr(str string) string {
	var result strings.Builder
	result.WriteByte('"')
	for _, c := range []byte(str) {
		escapeByte(&result, c)
	}
	result.WriteByte('"')
	return result.String()
}

func writeStrArrToJson(result *strings.Builder, arr []string) {
	result.WriteString("[")
	for i, value := range arr {
		if i > 0 {
			result.WriteString(",")
		}
		result.WriteString(fmt.Sprintf("\"%s\"", value))
	}
	result.WriteString("]")
}

func writeIntArrToJson(result *strings.Builder, arr []int) {
	result.WriteString("[")
	for i, value := range arr {
		if i > 0 {
			result.WriteString(",")
		}
		result.WriteString(strconv.Itoa(value))
	}
	result.WriteString("]")
}
