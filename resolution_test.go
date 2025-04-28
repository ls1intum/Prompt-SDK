package promptSDK

import (
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ls1intum/prompt-sdk/promptTypes"
	"github.com/stretchr/testify/assert"
)

func TestCoursePhaseParticipationsWithResolutionsValidation(t *testing.T) {
	// Get validator engine
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		t.Fatal("could not get validator engine")
	}

	// Create valid UUID for testing
	validUUID := uuid.MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479")

	// Create test cases
	tests := []struct {
		name        string
		resolutions []Resolution
		wantError   bool
		errorField  string
	}{
		{
			name: "Valid resolutions",
			resolutions: []Resolution{
				{
					DtoName:       "TestDTO",
					BaseURL:       "https://example.com",
					EndpointPath:  "/api/endpoint",
					CoursePhaseID: validUUID,
				},
			},
			wantError: false,
		},
		{
			name: "Missing DtoName",
			resolutions: []Resolution{
				{
					DtoName:       "", // Missing required field
					BaseURL:       "https://example.com",
					EndpointPath:  "/api/endpoint",
					CoursePhaseID: validUUID,
				},
			},
			wantError:  true,
			errorField: "DtoName",
		},
		{
			name: "Invalid BaseURL",
			resolutions: []Resolution{
				{
					DtoName:       "TestDTO",
					BaseURL:       "invalid-url", // Not a valid URL
					EndpointPath:  "/api/endpoint",
					CoursePhaseID: validUUID,
				},
			},
			wantError:  true,
			errorField: "BaseURL",
		},
		{
			name: "Missing EndpointPath",
			resolutions: []Resolution{
				{
					DtoName:       "TestDTO",
					BaseURL:       "https://example.com",
					EndpointPath:  "", // Missing required field
					CoursePhaseID: validUUID,
				},
			},
			wantError:  true,
			errorField: "EndpointPath",
		},
		{
			name: "Invalid CoursePhaseID",
			resolutions: []Resolution{
				{
					DtoName:       "TestDTO",
					BaseURL:       "https://example.com",
					EndpointPath:  "/api/endpoint",
					CoursePhaseID: uuid.UUID{}, // Zero UUID, invalid
				},
			},
			wantError:  true,
			errorField: "CoursePhaseID",
		},
		{
			name: "Multiple invalid fields",
			resolutions: []Resolution{
				{
					DtoName:       "",            // Invalid
					BaseURL:       "invalid-url", // Invalid
					EndpointPath:  "",            // Invalid
					CoursePhaseID: uuid.UUID{},   // Invalid
				},
			},
			wantError: true,
			// Multiple errors expected
		},
		{
			name: "Multiple resolutions with one invalid",
			resolutions: []Resolution{
				{
					DtoName:       "TestDTO1",
					BaseURL:       "https://example.com",
					EndpointPath:  "/api/endpoint1",
					CoursePhaseID: validUUID,
				},
				{
					DtoName:       "", // Invalid
					BaseURL:       "https://example.com",
					EndpointPath:  "/api/endpoint2",
					CoursePhaseID: validUUID,
				},
			},
			wantError:  true,
			errorField: "DtoName",
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Create a struct to validate
			cpp := CoursePhaseParticipationsWithResolutions{
				Participations: []promptTypes.CoursePhaseParticipationWithStudent{},
				Resolutions:    tt.resolutions,
			}

			// Validate the struct
			err := v.Struct(cpp)

			// Check if error matches expectation
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorField != "" {
					// Check if the error contains the expected field
					assert.Contains(t, err.Error(), tt.errorField)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmptyResolutions(t *testing.T) {
	// Get validator engine
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		t.Fatal("could not get validator engine")
	}

	// Test with empty resolutions slice
	cpp := CoursePhaseParticipationsWithResolutions{
		Participations: []promptTypes.CoursePhaseParticipationWithStudent{},
		Resolutions:    []Resolution{},
	}

	// Empty slice should be valid since dive only applies to elements within the slice
	err := v.Struct(cpp)
	assert.NoError(t, err)
}

func TestNilResolutions(t *testing.T) {
	// Get validator engine
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		t.Fatal("could not get validator engine")
	}

	// Test with nil resolutions slice
	cpp := CoursePhaseParticipationsWithResolutions{
		Participations: []promptTypes.CoursePhaseParticipationWithStudent{},
		Resolutions:    nil,
	}

	// Nil slice should be valid since dive only applies to elements within the slice
	err := v.Struct(cpp)
	assert.NoError(t, err)
}
