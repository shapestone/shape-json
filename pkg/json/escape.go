package json

// escapeTable maps ASCII bytes to their JSON escape character.
// 0 means no escape needed. Non-zero is the byte to write after backslash.
// Control characters (0x00-0x1F) that don't have named escapes use 0x01
// as a sentinel to indicate \u00XX encoding is needed.
var escapeTable [256]byte

const hexDigits = "0123456789abcdef"

func init() {
	// Named escapes
	escapeTable['"'] = '"'
	escapeTable['\\'] = '\\'
	escapeTable['/'] = '/'
	escapeTable['\b'] = 'b'
	escapeTable['\f'] = 'f'
	escapeTable['\n'] = 'n'
	escapeTable['\r'] = 'r'
	escapeTable['\t'] = 't'

	// Control characters without named escapes use sentinel value 0x01
	for i := byte(0); i < 0x20; i++ {
		if escapeTable[i] == 0 {
			escapeTable[i] = 0x01 // sentinel: needs \u00XX encoding
		}
	}
}

// appendEscapedString appends a JSON-escaped string to buf (without surrounding quotes).
// This is a zero-allocation function: it writes directly to the provided buffer.
func appendEscapedString(buf []byte, s string) []byte {
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 0x20 && c != '"' && c != '\\' && c != '/' {
			continue
		}

		// Flush unescaped run
		buf = append(buf, s[start:i]...)

		esc := escapeTable[c]
		if esc == 0x01 {
			// Control character: \u00XX
			buf = append(buf, '\\', 'u', '0', '0', hexDigits[c>>4], hexDigits[c&0x0F])
		} else {
			buf = append(buf, '\\', esc)
		}
		start = i + 1
	}
	// Flush remaining
	buf = append(buf, s[start:]...)
	return buf
}
