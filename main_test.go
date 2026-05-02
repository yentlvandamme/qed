package main

import (
	"bytes"
	"testing"
)

func TestModeParsing(t *testing.T) {
	tests := []struct {
		name         string
		valueToParse string
		parsedMode   Mode
		hasError     bool
	}{
		{"Parse to Replace", "rePLACe", Replace, false},
		{"Parse to Insert", "inSERt", Insert, false},
		{"Parse to Append", "appEND", Append, false},
		{"Parse invalid value", "invalid", Replace, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedValue, err := ParseMode(tt.valueToParse)
			hasError := err != nil
			if tt.hasError != hasError {
				t.Fatalf("Received %t but expected %t.\n", hasError, tt.hasError)

			}

			if parsedValue != tt.parsedMode && hasError == false {
				t.Fatalf("Received %d but expected %d.\n", parsedValue, tt.parsedMode)
			}
		})
	}
}

func TestInsertLine(t *testing.T) {
	tests := []struct {
		name           string
		mode           Mode
		expectedOutput string
	}{
		{name: "Appending a line", mode: Append, expectedOutput: "Just a line of text.\nfoo\nLet me know if you wanna edit something,\nBut don't be shy\n"},
		{name: "Inserting a line", mode: Insert, expectedOutput: "foo\nJust a line of text.\nLet me know if you wanna edit something,\nBut don't be shy\n"},
		{name: "Replacing a line", mode: Replace, expectedOutput: "foo\nLet me know if you wanna edit something,\nBut don't be shy\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputBuff := new(bytes.Buffer)
			var testString string = "Just a line of text.\nLet me know if you wanna edit something,\nBut don't be shy"

			inputBuff := bytes.NewBufferString(testString)

			args := Arguments{
				Text:       "foo",
				File:       inputBuff,
				LineNumber: 1,
				Mode:       tt.mode,
			}

			UpdateContent(args, outputBuff)

			if outputBuff.String() != tt.expectedOutput {
				t.Fatalf("Output incorrect.\nExpected: %q\nBut got: %q", tt.expectedOutput, outputBuff.String())
			}
		})
	}
}
