package main

import "testing"

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
