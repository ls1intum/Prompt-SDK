package keycloakTokenVerifier

import "testing"

func TestOnlyContainsAdminAndLecturer(t *testing.T) {
	tests := []struct {
		name string
		set  map[string]struct{}
		want bool
	}{
		{
			name: "nil map",
			set:  nil,
			want: true,
		},
		{
			name: "empty map",
			set:  map[string]struct{}{},
			want: true,
		},
		{
			name: "only admin",
			set:  map[string]struct{}{PromptAdmin: {}},
			want: true,
		},
		{
			name: "only lecturer",
			set:  map[string]struct{}{PromptLecturer: {}},
			want: true,
		},
		{
			name: "admin and lecturer",
			set:  map[string]struct{}{PromptAdmin: {}, PromptLecturer: {}},
			want: true,
		},
		{
			name: "contains other role",
			set:  map[string]struct{}{PromptAdmin: {}, "SomeOtherRole": {}},
			want: false,
		},
		{
			name: "only other role",
			set:  map[string]struct{}{"SomeOtherRole": {}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := onlyContainsAdminAndLecturer(tt.set)
			if got != tt.want {
				t.Errorf("onlyContainsAdminAndLecturer(%v) = %v, want %v", tt.set, got, tt.want)
			}
		})
	}
}

func TestContainsCustomRoleName(t *testing.T) {
	tests := []struct {
		name  string
		roles []string
		want  bool
	}{
		{"no roles", []string{}, false},
		{"only built-in", []string{PromptAdmin, PromptLecturer, CourseLecturer, CourseEditor, CourseStudent}, false},
		{"single custom", []string{"CustomRole"}, true},
		{"mix built-in and custom", []string{PromptAdmin, "X"}, true},
		{"multiple customs", []string{"X", "Y"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := containsCustomRoleName(tt.roles...)
			if got != tt.want {
				t.Errorf("containsCustomRoleName(%v) = %v; want %v", tt.roles, got, tt.want)
			}
		})
	}
}

func TestRequiresLecturerOrCustom(t *testing.T) {
	tests := []struct {
		name       string
		allowedSet map[string]struct{}
		roles      []string
		want       bool
	}{
		{"empty set, no roles", nil, []string{}, false},
		{"only lecturer allowed", map[string]struct{}{CourseLecturer: {}}, []string{}, true},
		{"only editor allowed", map[string]struct{}{CourseEditor: {}}, []string{}, true},
		{"lecturer and editor", map[string]struct{}{CourseLecturer: {}, CourseEditor: {}}, []string{}, true},
		{"no lec/editor but custom role in roles", map[string]struct{}{}, []string{"Custom"}, true},
		{"no lec/editor and no custom", map[string]struct{}{}, []string{PromptAdmin, CourseStudent}, false},
		{"editor but custom roles ignored", map[string]struct{}{CourseEditor: {}}, []string{"Whatever"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := requiresLecturerOrCustom(tt.allowedSet, tt.roles)
			if got != tt.want {
				t.Errorf("requiresLecturerOrCustom(%v, %v) = %v; want %v", tt.allowedSet, tt.roles, got, tt.want)
			}
		})
	}
}
