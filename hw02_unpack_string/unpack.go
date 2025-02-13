package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func isNumb(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func getNumb(ch rune) (int, error) {
	num, err := strconv.Atoi(string(ch))
	if err != nil {
		return 0, err
	}
	return num, nil
}

func repeatChar(resStr *strings.Builder, char rune, count int) {
	if count > 0 {
		str := strings.Repeat(string(char), count)
		resStr.WriteString(str)
	}
}

func processLastSimb(resStr *strings.Builder, char rune, isEcraning bool) error {
	if isEcraning {
		resStr.WriteRune(char)
	}
	if !isNumb(char) && !isEcraning {
		if char == '\\' {
			return ErrInvalidString
		}
		resStr.WriteRune(char)
	}
	return nil
}

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	var resStr strings.Builder
	var prevChar rune
	isEcraning := false

	for i, char := range s {
		if i == 0 {
			if isNumb(char) {
				return "", ErrInvalidString
			}
			prevChar = char
			continue
		}

		if prevChar == '\\' && !isEcraning {
			if !isNumb(char) && char != '\\' {
				return "", ErrInvalidString
			}
			isEcraning = true
			prevChar = char
			continue
		}

		if isEcraning && !isNumb(char) && char != '\\' {
			return "", ErrInvalidString
		}

		if isEcraning {
			isEcraning = false
			if isNumb(char) {
				num, _ := getNumb(char)

				repeatChar(&resStr, prevChar, num)

				prevChar = char
				continue
			}
			resStr.WriteRune(prevChar)
			prevChar = char
			continue
		}

		if isNumb(char) {
			if isNumb(prevChar) {
				return "", ErrInvalidString
			}
			num, _ := getNumb(char)

			repeatChar(&resStr, prevChar, num)

			prevChar = char
			continue
		}

		if !isNumb(prevChar) {
			resStr.WriteRune(prevChar)
		}
		prevChar = char
	}

	if err := processLastSimb(&resStr, prevChar, isEcraning); err != nil {
		return "", err
	}

	return resStr.String(), nil
}
