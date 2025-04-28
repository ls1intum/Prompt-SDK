package utils

// Package utils provides validation utilities for the Prompt SDK.
// It implements custom validators for TUM-specific identifiers
// that integrate with the gin-gonic/gin framework and go-playground/validator.
//
// Available validators:
//   - matriculationNumber: Validates TUM matriculation numbers (8 digits starting with '0')
//   - tumid: Validates TUM IDs (format: aa00aaa where 'a' is a lowercase letter and '0' is a digit)
//
// Usage with struct tags:
//   import "github.com/your-org/Prompt-SDK/utils"
//
//   type Student struct {
//     MatriculationNumber string `binding:"required,matriculationNumber"`
//     UniversityLogin     string `binding:"required,tumid"`
//   }
//
// Usage with Gin framework:
//   // Validators are automatically registered with Gin during package initialization
//
//   func CreateStudent(c *gin.Context) {
//     var student Student
//     if err := c.ShouldBindJSON(&student); err != nil {
//       c.JSON(400, gin.H{"error": err.Error()})
//       return
//     }
//     // Process valid student...
//   }
//
// Direct validation:
//   student := Student{...}
//   if err := utils.ValidateStruct(student); err != nil {
//     // Handle validation error
//   }

import (
	"fmt"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Register the custom validators for the gin binding
func init() {
	validate = validator.New()
	validate.SetTagName("binding")

	// Register with local validate instance first
	if err := validate.RegisterValidation("matriculationNumber", MatriculationNumberValidator); err != nil {
		panic(fmt.Sprintf("Failed to register local matriculationNumber validator: %v", err))
	}
	if err := validate.RegisterValidation("tumid", TUMIDValidator); err != nil {
		panic(fmt.Sprintf("Failed to register local tumid validator: %v", err))
	}

	// Also register with Gin's validator engine
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("matriculationNumber", MatriculationNumberValidator); err != nil {
			panic(fmt.Sprintf("Failed to register matriculationNumber validator: %v", err))
		}
		if err := v.RegisterValidation("tumid", TUMIDValidator); err != nil {
			panic(fmt.Sprintf("Failed to register tumid validator: %v", err))
		}
	}
}

// ValidateStruct validates a struct using the shared validator instance
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// MatriculationNumberValidator checks if a string is a valid matriculation number.
// A valid matriculation number must:
//   - Be exactly 8 characters long
//   - Start with the digit '0'
//   - Contain only numeric digits (0-9)
//
// This function is designed to be used as a custom validator with the validator package.
func MatriculationNumberValidator(fl validator.FieldLevel) bool {
	matriculationNumber := fl.Field().String()

	// Must be exactly 8 characters, start with '0', and all digits
	if len(matriculationNumber) != 8 {
		return false
	}
	if matriculationNumber[0] != '0' {
		return false
	}
	for _, r := range matriculationNumber {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// TUMIDValidator validates that a string follows the TUM ID format.
//
// A valid TUM ID consists of exactly 7 characters with the pattern:
// - First 2 characters: lowercase letters (a-z)
// - Next 2 characters: digits (0-9)
// - Last 3 characters: lowercase letters (a-z)
//
// Example of valid TUM ID: "ab12xyz"
//
// This function is designed to be used with the validator package as a custom validation function.
//
// Parameters:
//   - fl: validator.FieldLevel interface that provides access to the field being validated
//
// Returns:
//   - bool: true if the field contains a valid TUM ID, false otherwise
func TUMIDValidator(fl validator.FieldLevel) bool {
	tumID := fl.Field().String()

	if len(tumID) != 7 {
		return false
	}

	for i := 0; i < 2; i++ { // first two letters
		if tumID[i] < 'a' || tumID[i] > 'z' {
			return false
		}
	}
	for i := 2; i < 4; i++ { // two digits
		if tumID[i] < '0' || tumID[i] > '9' {
			return false
		}
	}
	for i := 4; i < 7; i++ { // last three letters
		if tumID[i] < 'a' || tumID[i] > 'z' {
			return false
		}
	}

	return true
}
