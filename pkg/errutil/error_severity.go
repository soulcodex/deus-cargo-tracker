package errutil

// ErrorSeverity represents the severity of an error.
type ErrorSeverity int

// Error severity levels.
const (
	SeverityInfo ErrorSeverity = iota
	SeverityWarning
	SeverityFatal
)

// IsValid reports whether s is a known severity.
func (s ErrorSeverity) IsValid() bool {
	return s >= SeverityInfo && s <= SeverityFatal
}

// String returns the string representation of the error severity.
// String returns a lowercase name for s, or "unknown".
func (s ErrorSeverity) String() string {
	switch s {
	case SeverityInfo:
		return "info"
	case SeverityWarning:
		return "warning"
	case SeverityFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

// ErrorWithSeverity is an error that has a severity level.
type ErrorWithSeverity interface {
	Error

	// Severity returns the severity of the error.
	Severity() ErrorSeverity
}
