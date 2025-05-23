package bencode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"bytes"
)

func Parse(reader *bufio.Reader) (interface{}, error) {
	// Peek at the first byte to determine the type
	char, err := reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("failed to read initial byte: %w", err)
	}

	switch char {
	case 'i': // Integer: i<integer>e
		return parseInteger(reader)
	case 'l': // List: l<bencoded values>e
		return parseList(reader)
	case 'd': // Dictionary: d<bencoded string><bencoded value>...e
		return parseDictionary(reader)
	default: // String: <length>:<string data>
		if char >= '0' && char <= '9' {
			// Put the digit back to be read as part of the length
			if err := reader.UnreadByte(); err != nil {
				return nil, fmt.Errorf("failed to unread byte for string length: %w", err)
			}
			return parseString(reader)
		}
		return nil, fmt.Errorf("unexpected bencode token: '%c' (ASCII %d)", char, char)
	}
}


// Parsing Integer

func parseInteger(reader *bufio.Reader) (int64, error) {
	// Read until 'e'
	strBytes, err := reader.ReadBytes('e')
	if err != nil {
		return 0, fmt.Errorf("integer not terminated by 'e' or read error: %w", err)
	}

	// Convert bytes to string, removing the trailing 'e'
	numStr := string(strBytes[:len(strBytes)-1])

	if len(numStr) == 0 {
		return 0, errors.New("empty integer string")
	}

	// Special cases: "i-0e" is invalid. "i0e" is valid. "i03e" is invalid.
	if numStr[0] == '-' {
		if len(numStr) > 1 && numStr[1] == '0' {
			return 0, errors.New("invalid negative integer format (e.g., i-0e)")
		}
	} else {
		if len(numStr) > 1 && numStr[0] == '0' {
			return 0, errors.New("invalid integer format (leading zero, e.g., i03e)")
		}
	}

	val, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse integer string '%s': %w", numStr, err)
	}
	return val, nil
}


// Parsing String

func parseString(reader *bufio.Reader) (string, error) {
	// Read length part until ':'
	lenBytes, err := reader.ReadBytes(':')
	if err != nil {
		return "", fmt.Errorf("string length delimiter ':' not found or read error: %w", err)
	}

	// Convert length bytes to string, removing the trailing ':'
	lenStr := string(lenBytes[:len(lenBytes)-1])
	if len(lenStr) == 0 {
		return "", errors.New("empty string length")
	}

	length, err := strconv.Atoi(lenStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse string length '%s': %w", lenStr, err)
	}
	if length < 0 {
		return "", fmt.Errorf("invalid string length: %d", length)
	}

	// Read the string data itself
	buf := make([]byte, length)
	n, err := io.ReadFull(reader, buf)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return "", fmt.Errorf("unexpected end of input while reading string data, expected %d bytes, got %d: %w", length, n, err)
		}
		return "", fmt.Errorf("failed to read string data (length %d): %w", length, err)
	}

	return string(buf), nil
}


// Parsing List

func parseList(reader *bufio.Reader) ([]interface{}, error) {
	var list []interface{}
	for {
		// Peek at the next character
		char, err := reader.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("failed to read byte for list item or end: %w", err)
		}

		if char == 'e' { // End of list
			return list, nil
		}

		// Put the character back to be parsed by BencodeParse
		if err := reader.UnreadByte(); err != nil {
			return nil, fmt.Errorf("failed to unread byte for list item: %w", err)
		}

		item, err := Parse(reader) // Recursive call
		if err != nil {
			return nil, fmt.Errorf("failed to parse list item: %w", err)
		}
		list = append(list, item)
	}
}


// Parsing Dictionary

func parseDictionary(reader *bufio.Reader) (map[string]interface{}, error) {
	dict := make(map[string]interface{})
	var lastKey string // To check if keys are sorted (optional validation)

	for {
		// Peek at the next character for end of dictionary or key
		char, err := reader.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("failed to read byte for dict key or end: %w", err)
		}

		if char == 'e' { // End of dictionary
			return dict, nil
		}

		// Put the character back to be parsed as a string key
		if err := reader.UnreadByte(); err != nil {
			return nil, fmt.Errorf("failed to unread byte for dict key: %w", err)
		}

		// Keys in dictionaries must be strings
		key, err := parseString(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to parse dictionary key: %w", err)
		}

		if len(dict) > 0 && key <= lastKey {
		}
		lastKey = key

		// Parse the value associated with the key
		value, err := Parse(reader) // Recursive call
		if err != nil {
			return nil, fmt.Errorf("failed to parse dictionary value for key '%s': %w", key, err)
		}
		dict[key] = value
	}
}

// ParseString is a convenience helper to parse a bencoded string directly.
func ParseString(bencodedString string) (interface{}, error) {
	reader := bufio.NewReader(bytes.NewReader([]byte(bencodedString)))
	return Parse(reader)
}

