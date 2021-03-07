//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package vt100

import (
	"bytes"
	"regexp"
)

var reLine = regexp.MustCompilePOSIX(`^[[:xdigit:]]{8}(([[:blank:]]+[[:xdigit:]]{2}){1,16}).*$`)
var reByte = regexp.MustCompilePOSIX(`[[:blank:]]+([[:xdigit:]]{2})`)

// ParseHexDump parses data from the encoding/hex.Dump formatted
// output.
func ParseHexDump(data []byte) ([]byte, error) {
	var result bytes.Buffer

	for {
		match := reLine.FindSubmatchIndex(data)
		if match == nil {
			return result.Bytes(), nil
		}
		bytes := data[match[2]:match[3]]
		data = data[match[1]:]

		for {
			m := reByte.FindSubmatchIndex(bytes)
			if m == nil {
				break
			}
			result.WriteByte(hex2bin(bytes[m[2]])<<4 | hex2bin(bytes[m[2]+1]))
			bytes = bytes[m[1]:]
		}
	}
}

func hex2bin(h byte) byte {
	switch h {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return h - '0'

	case 'a', 'b', 'c', 'd', 'e', 'f':
		return h - 'a' + 10

	case 'A', 'B', 'C', 'D', 'E', 'F':
		return h - 'A' + 10

	default:
		return 0
	}
}
