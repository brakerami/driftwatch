package reporter

import "fmt"

// String returns the string representation of a Format.
func (f Format) String() string {
	return string(f)
}

// ParseFormat converts a raw string into a Format value.
// It returns an error if the string does not match a known format.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatText, FormatJSON:
		return Format(s), nil
	default:
		return "", fmt.Errorf("reporter: unknown format %q; valid values are \"text\", \"json\"", s)
	}
}

// MarshalText implements encoding.TextMarshaler so Format serialises
// correctly when embedded in JSON structs.
func (f Format) MarshalText() ([]byte, error) {
	return []byte(f), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (f *Format) UnmarshalText(b []byte) error {
	parsed, err := ParseFormat(string(b))
	if err != nil {
		return err
	}
	*f = parsed
	return nil
}
