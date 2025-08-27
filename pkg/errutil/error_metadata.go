package errutil

// ErrorMetadata is a key-value map that can be used to store additional information about the error.
type ErrorMetadata struct {
	metadata map[string]interface{}
}

// IsEmpty returns true if the metadata is empty.
func (m ErrorMetadata) IsEmpty() bool {
	return len(m.metadata) == 0
}

// Get returns the value for the given key.
func (m ErrorMetadata) Get(key string) interface{} {
	if m.metadata == nil {
		return nil
	}

	return m.metadata[key]
}

// Set sets the value for the given key.
func (m ErrorMetadata) Set(key string, value interface{}) ErrorMetadata {
	if m.metadata == nil {
		m.metadata = make(map[string]interface{})
	}

	m.metadata[key] = value
	return m
}

// AsMap returns the metadata as a map.
func (m ErrorMetadata) AsMap() map[string]interface{} {
	return m.metadata
}

// NewErrorMetadata creates a new empty ErrorMetadata instance.
func NewErrorMetadata() ErrorMetadata {
	return ErrorMetadata{metadata: make(map[string]interface{})}
}

// ErrorMetadataFromMap creates a new ErrorMetadata from the given map.
func ErrorMetadataFromMap(metadata map[string]interface{}) ErrorMetadata {
	return ErrorMetadata{metadata: metadata}
}
