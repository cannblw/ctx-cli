package slug

import gosimple "github.com/gosimple/slug"

// FromValue generates a URL-friendly slug from a value string.
// Truncates to 80 chars max.
func FromValue(value string) string {
	s := gosimple.Make(value)
	if len(s) > 80 {
		s = s[:80]
	}
	return s
}
