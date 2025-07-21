package promptTypes

// StudyDegree represents the academic degree level that a student is pursuing.
// This type is used to categorize students based on their current academic level
// and helps determine course eligibility, academic requirements, and progression rules.
type StudyDegree string

// Study degree constants defining the available academic levels in the system.
const (
	// StudyDegreeBachelor represents undergraduate bachelor's degree programs.
	// Typically the first university degree, usually taking 3-4 years to complete.
	StudyDegreeBachelor StudyDegree = "bachelor"

	// StudyDegreeMaster represents graduate master's degree programs.
	// Advanced degree that typically follows a bachelor's degree, usually taking 1-2 years to complete.
	StudyDegreeMaster StudyDegree = "master"
)
