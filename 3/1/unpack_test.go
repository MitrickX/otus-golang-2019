package unpack

import "testing"

func TestUnpackEmpty(t *testing.T) {
	input := ""
	output, err := unpack(input)

	if err != nil {
		t.Errorf("Error not nil: %s\n", err)
	}

	if output != "" {
		t.Errorf("Output is not empty: %s\n", output)
	}
}

func TestUnpackHasDigitInBeginning(t *testing.T) {
	input := "4abe2zu2"
	output, err := unpack(input)

	if err == nil {
		t.Errorf("Error is nil\n")
	} else if err.Error() != "Invalid input: instruction digit in the beginning of input" {
		t.Errorf("Unexpected error: %s\n", err)
	}

	if output != "" {
		t.Errorf("Output is not empty when error occurs: %s\n", output)
	}
}

func TestUnpackHasTwoDigitsOneAfterAnother(t *testing.T) {
	input := "q4abe23zu2"
	output, err := unpack(input)

	if err == nil {
		t.Errorf("Error is nil\n")
	} else if err.Error() != "Invalid input: two instruction digits go one after another" {
		t.Errorf("Unexpected error: %s\n", err)
	}

	if output != "" {
		t.Errorf("Output is not empty when error occurs: %s\n", output)
	}
}

func TestUnpackValidMixed(t *testing.T) {
	input := "a4bc2d5e"
	output, err := unpack(input)

	if err != nil {
		t.Errorf("Error is not nil: %s\n", err)
		return
	}

	expected := "aaaabccddddde"
	if output != expected {
		t.Errorf("Output is incorrect, must be \"%s\" instread of \"%s\"\n", expected, output)
	}
}

func TestUnpackValidOnlySymbols(t *testing.T) {
	input := "abcd"
	output, err := unpack(input)

	if err != nil {
		t.Errorf("Error is not nil: %s\n", err)
		return
	}

	expected := "abcd"
	if output != expected {
		t.Errorf("Output is incorrect, must be \"%s\" instread of \"%s\"\n", expected, output)
	}
}

func TestUnpackHasNoInstructionDigits(t *testing.T) {
	input := "ab11c02qw12d"
	output, err := unpack(input)

	if err != nil {
		t.Errorf("Error is not nil: %s\n", err)
		return
	}

	expected := "ab11c00qw11d"
	if output != expected {
		t.Errorf("Output is incorrect, must be \"%s\" instread of \"%s\"\n", expected, output)
	}
}
