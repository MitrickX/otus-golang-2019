package unpack

import (
	"errors"
	"unicode"
)

// Repeat rune n times and append all in the end of stream (slice of runes)
// If n is 0, no append
// If n is 1, no repeating (append 1 time)
func repeatRune(stream []rune, r rune, n int) []rune {

	// reapeat + append
	for i := 0; i < n; i++ {
		stream = append(stream, r)
	}

	// result
	return stream
}

// Convert rune into digit (as int) and status of converting
func toDigit(symbol rune) (int, bool) {
	if !unicode.IsDigit(symbol) {
		return 0, false
	}
	digit := int(symbol - '0')
	return digit, true
}

//
// Unpack compact input string represented by stream of digits and symbols
//
// Digit has special meaning - it says how much times need to repeat previous symbol in output
// Digit 0 not output previous symbol at all
// Digit 1 not repeat previous symbol
//
//
// Invalid input string when
//  - Two digits go one after another in input
//  - Digit is in the beginning of input
//
// Output will has no digits, only symbols that could repeated by instructions
// On error output is empty string
//
// Examples:
//   "a4bc2d5e" => "aaaabccddddde"
//   "abcd" => "abcd"
//   "45" => "" (invalid input)
//
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

	// first symbol is digit => error
	if unicode.IsDigit(stream[0]) {
		return "", errors.New("Invalid input: digit in the beginning of input")
	}

	// will unpacking into this slice of rune, in the end it is our result
	var result []rune

	// symbol that need to repeat
	var candidate rune

	// how much times need to repeat symbol, default is 1 (without repeating)
	var count int = 1

	// is previous symbol digit
	var isDigit = false

	// go through stream of symbols and digits (e.g. a14bc2d5e)
	// if consume digit it would be instruction for repeating
	// otherwise repeat previous candidate and set current symbol as a new candidate
	for _, symbol := range stream {

		// read symbol and try to convert to digit (int)
		digit, ok := toDigit(symbol)

		// it is digit
		if ok {

			// there was consumed digit on previous step - two digits case => error
			if isDigit {
				return "", errors.New("Invalid input: two digits go one after another")
			}

			// remeber digit and a fact that there is a digit
			count = digit
			isDigit = true

			// go to next symbol
			continue
		}

		// repeating symbol (ignore inital case, when noting to repeat yet)
		if candidate != 0 {
			result = repeatRune(result, candidate, count)
		}

		// next candidate to repeat
		candidate = symbol

		// reset to defaults
		count = 1
		isDigit = false

	}

	return string(result), nil
}
