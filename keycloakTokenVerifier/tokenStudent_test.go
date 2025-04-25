package keycloakTokenVerifier

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestGetTokenStudent_NoTokenSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	token, ok := GetTokenStudent(c)
	if ok {
		t.Errorf("expected ok=false, got true")
	}
	if !reflect.DeepEqual(token, TokenStudent{}) {
		t.Errorf("expected zero value TokenStudent, got %#v", token)
	}
}

func TestSetAndGetTokenStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	courseID := uuid.New()
	expected := TokenStudent{
		Roles:                  map[string]bool{"admin": true, "user": false},
		ID:                     uuid.New().String(),
		Email:                  "student@example.com",
		MatriculationNumber:    "MATRIC123",
		UniversityLogin:        "uni_login",
		FirstName:              "Alice",
		LastName:               "Smith",
		IsStudentOfCourse:      true,
		IsStudentOfCoursePhase: true,
		CourseParticipationID:  courseID,
		IsLecturer:             false,
		IsEditor:               true,
		CustomRolePrefix:       "custom_",
	}

	SetTokenStudent(c, expected)
	got, ok := GetTokenStudent(c)
	if !ok {
		t.Fatalf("expected ok=true, got false")
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("mismatch:\n got %#v\nwant %#v", got, expected)
	}
}
