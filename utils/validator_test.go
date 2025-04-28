package utils

import (
	"testing"
)

// This test checks if the custom validators are registered correctly using the init function.
func TestInitRegistersValidators(t *testing.T) {
	type TestStruct struct {
		Matriculation string `binding:"matriculationNumber"`
		TUMID         string `binding:"tumid"`
	}

	tests := []struct {
		name      string
		input     TestStruct
		wantError bool
	}{
		{"Valid Matriculation and TUMID", TestStruct{"01234567", "ab12cde"}, false},
		{"Invalid Matriculation (no leading 0)", TestStruct{"11234567", "ab12cde"}, true},
		{"Invalid Matriculation (too short)", TestStruct{"0123456", "ab12cde"}, true},
		{"Invalid TUMID (wrong format)", TestStruct{"01234567", "a112cde"}, true},
		{"Invalid TUMID (too short)", TestStruct{"01234567", "ab12cd"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("input: %+v, got error = %v, wantError = %v", tt.input, err, tt.wantError)
			}
		})
	}
}

// TestMatriculationNumberValidator tests the MatriculationNumberValidator function.
func TestMatriculationNumberValidator(t *testing.T) {
	type TestStruct struct {
		MatriculationNumber string `binding:"matriculationNumber"`
	}

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid number", "03745306", false},
		{"Too short", "0374530", true},
		{"Too long", "037453061", true},
		{"Does not start with 0", "13745306", true},
		{"Contains letters", "03A45306", true},
		{"Empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			test := TestStruct{MatriculationNumber: tt.input}
			err := ValidateStruct(test)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

// TestTUMIDValidator tests the TUMIDValidator function.
func TestTUMIDValidator(t *testing.T) {
	type TestStruct struct {
		TUMID string `binding:"tumid"`
	}

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid TUM ID", "ge23das", false},
		{"Too short", "ge23da", true},
		{"Too long", "ge23dass", true},
		{"First letters uppercase", "Ge23das", true},
		{"Digits in wrong place", "g2e3das", true},
		{"Letters in digit spot", "ge2adas", true},
		{"Numbers at the end", "ge23da3", true},
		{"All digits", "1234567", true},
		{"All letters", "abcdefg", true},
		{"Empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			test := TestStruct{TUMID: tt.input}
			err := ValidateStruct(test)
			if (err != nil) != tt.wantErr {
				t.Errorf("input: %q, got error = %v, wantErr = %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
