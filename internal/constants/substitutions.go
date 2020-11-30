package constants

import "regexp"

var (
	RegexStatic          = regexp.MustCompile("^/static/$")
	RegexManifestJSON    = regexp.MustCompile("^/manifest.json$")
	RegexDataReactHelmet = regexp.MustCompile(`[ ]*data\-react\-helmet=".*"[ ]*`)
)
