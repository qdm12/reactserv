package constants

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Regex(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		regex *regexp.Regexp
		s     string
		match bool
	}{
		"RegexStatic mismatch": {
			regex: RegexStatic,
			s:     "/static/ ",
		},
		"RegexStatic match": {
			regex: RegexStatic,
			s:     "/static/",
			match: true,
		},
		"RegexManifestJSON mismatch": {
			regex: RegexManifestJSON,
			s:     "/manifest.json/",
		},
		"RegexManifestJSON match": {
			regex: RegexManifestJSON,
			s:     "/manifest.json",
			match: true,
		},
		"RegexDataReactHelmet mismatch": {
			regex: RegexDataReactHelmet,
			s:     "data-react-helmet=",
		},
		"RegexDataReactHelmet match": {
			regex: RegexDataReactHelmet,
			s:     ` data-react-helmet="true"`,
			match: true,
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			match := testCase.regex.MatchString(testCase.s)
			assert.Equal(t, testCase.match, match)
		})
	}
}
