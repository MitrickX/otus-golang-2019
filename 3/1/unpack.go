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

// Convert rune into instruction int and status of converting
// Instruction digit is digit in interval [2-9]
func toInstructionDigit(symbol rune) (int, bool) {
	if !unicode.IsDigit(symbol) {
		return 0, false
	}
	digit := int(symbol - '0')
	if digit > 1 {
		return digit, true
	} else {
		return 0, false
	}
}

// Is rune instruction digit
// Instruction digit is digit in interval [2-9]
func isInstructionDigit(symbol rune) bool {
	_, ok := toInstructionDigit(symbol)
	return ok
}

//
// Unpack compact input string represented by stream of digits [2-9] and symbols
// Further in here, digit [2-9] would be called as instruction digit
//
// Instruction has special meaning - it says how much times need to repeat previous symbol in output
// Other digits - 0 and 1 - interpreted just as a symbol, without special meanings
//
// Invalid input string when
//  - Two instruction digits go one after another in input
//  - Instruction digit is in the beginning of input
//
// Output will has no instruction digits, only symbols that could repeated by instructions
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

	// first symbol is instruction digit, error
	if isInstructionDigit(stream[0]) {
		return "", errors.New("Invalid input: instruction digit in the beginning of input")
	}

	// will unpacking into this slice of rune, in the end it is our result
	var result []rune

	// symbol that need to repeat
	var candidate rune

	// instruction how much times need to repeat symbol
	var instruction int

	// go through stream of symbols and digits (e.g. a14bc2d5e)
	// if consume digit [2-9] it would be instruction for repeating
	// otherwis repeat previous candidate and set symbol as a new candidate
	for _, symbol := range stream {

		// read symbol and try to convert to instruction
		digit, ok := toInstructionDigit(symbol)

		// it is instruction
		if ok {

			// there was instruction on previous step - two instruction digits case => error
			if instruction > 0 {
				return "", errors.New("Invalid input: two instruction digits go one after another")
			}

			// remeber instruction
			instruction = digit

			// go to next symbol
			continue
		}

		// ignore "empty" candidate - nothing to reapeat yet
		if candidate != 0 {
			result = repeatRune(result, candidate, instruction)
		}

		// next candidate to repeat
		candidate = symbol

		// reset instruction
		instruction = 0
	}

	return string(result), nil
}
