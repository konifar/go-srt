package gosrt

// String converts string to *string
func String(s string) *string {
	p := s
	return &p
}
