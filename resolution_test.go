package promptSDK

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetEndpointPath(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple path", "path", "path"},
		{"leading slash", "/path", "path"},
		{"trailing slash", "path/", "path"},
		{"leading and trailing slash", "/path/", "path"},
		{"nested with slashes", "//nested/path//", "nested/path"},
		{"empty string", "", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got := getEndpointPath(c.input)
			assert.Equal(t, c.expected, got)
		})
	}
}

func TestBuildURL_NoExtraPaths(t *testing.T) {
	id := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	res := Resolution{
		BaseURL:       "https://example-prompt.com/api",
		CoursePhaseID: id,
		EndpointPath:  "/my-endpoint/",
	}
	got := buildURL(res)
	want := "https://example-prompt.com/api/course_phase/123e4567-e89b-12d3-a456-426614174000/my-endpoint"
	assert.Equal(t, want, got)
}

func TestBuildURL_WithExtraPaths(t *testing.T) {
	id := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	res := Resolution{
		BaseURL:       "http://localhost:8080/v1",
		CoursePhaseID: id,
		EndpointPath:  "endpoint",
	}
	got := buildURL(res, "p1", "details")
	want := "http://localhost:8080/v1/course_phase/00000000-0000-0000-0000-000000000000/endpoint/p1/details"
	assert.Equal(t, want, got)
}

func TestBuildURL_WithInvalidBaseURL(t *testing.T) {
	// Test with an invalid URL that would cause issues
	res := Resolution{
		BaseURL:       ":%invalid",
		CoursePhaseID: uuid.New(),
		EndpointPath:  "endpoint",
	}
	got := buildURL(res)
	// Verify the function gracefully handles invalid URLs
	assert.Empty(t, got)
}
