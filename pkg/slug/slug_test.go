package slug_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cannblw/ctx-cli/pkg/slug"
)

// ── FromValue successes ───────────────────────────────────────────────────────

func TestFromValue_SuccessURL(t *testing.T) {
	s := slug.FromValue("https://github.com/org/repo/pull/42")
	assert.Equal(t, "https-github-com-org-repo-pull-42", s)
}

func TestFromValue_SuccessLongString(t *testing.T) {
	long := strings.Repeat("a", 200)
	s := slug.FromValue(long)
	assert.LessOrEqual(t, len(s), 80)
}

func TestFromValue_SuccessSpecialChars(t *testing.T) {
	s := slug.FromValue("Fix #123: Auth Redirect Bug!")
	assert.Equal(t, "fix-123-auth-redirect-bug", s)
}
