package utils

// Sptr returns a pointer to a string.
func Sptr(s string) *string {
	return &s
}

// Iptr - int to *int
func Iptr(i int) *int {
	return &i
}

// Fptr - float64 to *float64
func Fptr(f float64) *float64 {
	return &f
}

// Bptr - bool to *bool
func Bptr(b bool) *bool {
	return &b
}

// StringFromSptr returns either the dereferenced string or "<nil>"
func StringFromSptr(sp *string) string {
	if sp == nil {
		return "<nil>"
	}
	return "&`" + *sp + "`"
}
