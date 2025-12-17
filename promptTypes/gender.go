package promptTypes

// Gender represents the gender identity options available in the Prompt system.
// This type is used for demographic data collection and analysis while respecting
// privacy and inclusivity requirements.
type Gender string

// Gender constants defining the available gender options.
// These values are used in forms, databases, and API responses throughout the system.
const (
	// GenderMale represents male gender identity.
	GenderMale Gender = "male"

	// GenderFemale represents female gender identity.
	GenderFemale Gender = "female"

	// GenderDiverse represents non-binary, third gender, or other diverse gender identities.
	GenderDiverse Gender = "diverse"

	// GenderPreferNotToSay allows users to opt out of providing gender information
	// while still maintaining data completeness for optional demographic tracking.
	GenderPreferNotToSay Gender = "prefer_not_to_say"
)
