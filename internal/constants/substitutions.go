package constants

import "regexp"

var (
	RegexStatic       = regexp.MustCompile("^/static/$")
	RegexManifestJSON = regexp.MustCompile("^/manifest.json$")
)
