package core

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const CRLF string = "\r\n"

var RespNil = []byte("$-1\r\n")

// +OK\r\n => OK, 5
func readSimpleString(data []byte) (string, int, error) {
	pos := 1
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	return string(data[1:pos]), pos + 2, nil
}

// :123\r\n => 123
func readInt64(data []byte) (int64, int, error) {
	pos := 1
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	value, _ := strconv.ParseInt(string(data[1:pos]), 10, 64)
	return value, pos + 2, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

// $5\r\nhello\r\n => 5, 4
func readLen(data []byte) (int, int) {
	res, pos, _ := readInt64(data)
	return int(res), pos
}

// $5\r\nhello\r\n => "hello"
func readBulkString(data []byte) (string, int, error) {
	length, pos := readLen(data)
	return string(data[pos: pos+length]), pos + length + 2, nil
}

// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n => {"hello", "world"}
func readArray(data []byte) (any, int, error) {
	length, pos := readLen(data)
	var res []any = make([]any, length)
	// implement start
	for i := 0; i < length; i++ {
		val, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		res[i] = val
		pos += delta
	}
	// implement end
	return res, pos, nil
}

func DecodeOne(data []byte) (any, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case ':':
		return readInt64(data)
	case '-':
		return readError(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}
	return nil, 0, nil
}

func Decode(data []byte) (any, error) {
	res, _, err := DecodeOne(data)
	return res, err
}

func encodeString(s string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
}

func encodeStringArray(sa []string) []byte {
	var b []byte
	buf := bytes.NewBuffer(b)
	for _, s := range sa {
		buf.Write(encodeString(s))
	}
	return []byte(fmt.Sprintf("*%d\r\n%s", len(sa), buf.Bytes()))
}

func Encode(value any, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimpleString {
			return []byte(fmt.Sprintf("+%s%s", v, CRLF))
			// return fmt.Appendf([]byte{}, "+%s%s", v, CRLF)
		}
		return []byte(fmt.Sprintf("$%d%s%s%s", len(v), CRLF, v, CRLF))
	case int64, int32, int16, int8, int:
		return []byte(fmt.Sprintf(":%d\r\n", v))
	case error:
		return []byte(fmt.Sprintf("-%s\r\n", v))
	case []string:
		return encodeStringArray(value.([]string))
	case [][]string:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, sa := range value.([][]string) {
			buf.Write(encodeStringArray(sa))
		}
		return []byte(fmt.Sprintf("*%d\r\n%s", len(value.([][]string)), buf.Bytes()))
	case []any:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, x := range value.([]any) {
			buf.Write(Encode(x, false))
		}
		return []byte(fmt.Sprintf("*%d\r\n%s", len(value.([]any)), buf.Bytes()))
	default:
		return RespNil
	}
}

func ParseCmd(data []byte) (*Command, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}

	array := value.([]any)
	tokens := make([]string, len(array))
	for i := range tokens {
		tokens[i] = array[i].(string)
	}
	res := &Command{Cmd: strings.ToUpper(tokens[0]), Args: tokens[1:]}
	return res, nil
}