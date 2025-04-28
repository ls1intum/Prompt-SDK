package utils

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Register the custom validators for the gin binding
func init() {
	validate = validator.New()
	validate.SetTagName("binding")
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

// MatriculationNumberValidator checks if a string is exactly 8 digits and starts with 0
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
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

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
