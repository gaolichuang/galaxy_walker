package url_parser

import (
        LOG "galaxy_walker/internal/gcodebase/log"
        "galaxy_walker/internal/github.com/PuerkitoBio/purell"
	"net/url"
)

func GetHost(_url string) string {
	u := GetURLObj(_url)
	if u == nil {
		return ""
	}
	return u.Host
}
func GetURLObj(_url string) *url.URL {
	var u *url.URL
	var e error
	if u, e = url.Parse(_url); e != nil {
		LOG.Error("Parse Url Fail")
		return nil
	}
	return u
}
func NormalizeUrl(_url string) string {
	normal_url, err := purell.NormalizeURLString(_url, purell.FlagsSafe|purell.FlagRemoveDotSegments|purell.FlagRemoveFragment)
	if err != nil {
		return ""
	}
	return normal_url
}
