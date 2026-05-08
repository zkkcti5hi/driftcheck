package audit

import "bytes"

// newBytesReader returns a *bytes.Reader for the given slice.
// Extracted so the main file stays free of the bytes import.
func newBytesReader(b []byte) *bytes.Reader {
	return bytes.NewReader(b)
}
