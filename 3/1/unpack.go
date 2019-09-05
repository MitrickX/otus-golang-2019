package unpack

import (
	"errors"
	"unicode"
)

// Repeat rune at least one time and append all in the end of stream (slice of runes)
func repeatRune(stream []rune, r rune, n int) []rune {

	// at least one time
	if n <= 0 {
		n = 1
	}

	// reapeat + append
	for i := 0; i < n; i++ {
		stream = append(stream, r)
	}

	// result
	return stream
}

// Unpack compact input string represented by stream of digits [2-9] and symbols
// Each digit say how much times need to repeat previous symbol in output
//
// Invalid input string when
//  - Digit 0 or 1 in input
//  - Two digits one after another in input
//  - Digit in beginning of input
//
// Output has no digits, only symbols that could be repeated
// On error output is empty string
//
// Examples:
//   "a4bc2d5e" => "aaaabccddddde"
//   "abcd" => "abcd"
//   "45" => "" (invalid input)
func unpack(input string) (string, error) {

	// boundary case
	if len(input) == 0 {
		return "", nil
	}

	// represent input as slice of runes (to append terminate symbol)
	var stream = []rune(input)

	// append terminate symbol in the end to prevent code dupblicate
	// without terminate symbol we have to make a repeat logic one more time after loop (flushing in-between states)
	stream = append(stream, 0)

	// first symbol is digit, error
	if unicode.IsDigit(stream[0]) {
		return "", errors.New("Invalid input: digit in beginning of input")
	}

	// will unpacking into this slice of rune, in the end it is our result
	var result []rune

	// symbol that need to repeat
	var candidate rune

	// how much times need to repeat symbol
	var count int

	// go through stream of symbols and digits (e.g. a14bc2d5e)
	// if consume digit it would be count for repeating
	// if consume symbol repeat previous candidate and set symbol as a new candidate
	for _, symbol := range stream {

		is_digit := unicode.IsDigit(symbol)

		if is_digit {

			// there is already count - two digits case, error
			if count > 0 {
				return "", errors.New("Invalid input: two digits one after another")
			}

			// digit symbol to int
			count = int(symbol - '0')

			// 0 in input stream, error
			if count == 0 {
				return "", errors.New("Invalid input: digit 0 in input")
			}
			// 1 in input stream, error
			if count == 1 {
				return "", errors.New("Invalid input: digit 1 in input")
			}

			// go to next symbol
			continue
		}

		// ignore "emptiness" (initial value of rune)
		// "emptiness" means we has not candidate to repeat yet
		if candidate != 0 {
			result = repeatRune(result, candidate, count)
		}

		// next candidate to repeat
		candidate = symbol

		// reset counter
		count = 0
	}

	return string(result), nil
}

func main() {

}
