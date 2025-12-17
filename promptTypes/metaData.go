package promptTypes

// MetaData represents flexible key-value storage for dynamic data throughout the Prompt system.
// This type allows modules to store and retrieve arbitrary structured data without requiring
// schema changes to core types. Common use cases include:
//   - Application answers and form data
//   - Configuration settings
//   - Analytics and tracking information
//   - Module-specific state and preferences
//
// The interface{} values can contain strings, numbers, booleans, arrays, or nested objects,
// making this type suitable for JSON serialization and flexible data exchange between modules.
type MetaData map[string]interface{}
