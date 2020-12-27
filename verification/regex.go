package verification

import "regexp"

type Regex interface {
	MatchEmail(s string) bool
	MatchPhoneIntl(s string) bool
	MatchPhoneLocal(s string) bool
	MatchDomain(s string) bool
	MatchHostname(s string) bool
	MatchRootURL(s string) bool
	Match64BytesHex(s string) bool
	MatchMD5String(s string) bool
	SearchIPv4(s string) (results []string)
	SearchIPv6(s string) (results []string)
	SearchEmail(s string) (results []string)
	SearchPhone(s string) (results []string)
}

type regex struct {
	searchIPv4      *regexp.Regexp
	searchIPv6      *regexp.Regexp
	searchEmail     *regexp.Regexp
	searchPhone     *regexp.Regexp
	matchEmail      *regexp.Regexp
	matchPhoneIntl  *regexp.Regexp
	matchPhoneLocal *regexp.Regexp
	matchDomain     *regexp.Regexp
	matchHostname   *regexp.Regexp
	matchRootURL    *regexp.Regexp
	match32BytesHex *regexp.Regexp
	match64BytesHex *regexp.Regexp
}

func NewRegex() Regex {
	//nolint:lll
	const (
		regexIPv4        = `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`
		regexIPv6        = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
		regexEmail       = `[a-zA-Z0-9-_.+]+@[a-zA-Z0-9-_.]+\.[a-zA-Z]{2,10}`
		regexPhoneSearch = `(\+|( *0 *0)){0,1}[0-9][0-9 ]{7}[0-9 ]*[0-9]`
		regexPhoneIntl   = `(\+|00)[0-9]{9,15}`
		regexPhoneLocal  = `([1-9][0-9]{4,13}|0[1-9][0-9]{3,12})`
		regexDomain      = `(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})`
		regexHostname    = `([a-zA-Z0-9]|[a-zA-Z0-9_][a-zA-Z0-9\-_]{0,61}[a-zA-Z0-9_])(\.([a-zA-Z0-9]|[a-zA-Z0-9_][a-zA-Z0-9\-_]{0,61}[a-zA-Z0-9]))*`
		regexRootURL     = `\/[a-zA-Z0-9\-_/\+]*`
		regex32BytesHex  = `[a-fA-F0-9]{32}`
		regex64BytesHex  = `[a-fA-F0-9]{64}`
	)
	return &regex{
		searchIPv4:      regexp.MustCompile(regexIPv4),
		searchIPv6:      regexp.MustCompile(regexIPv6),
		searchPhone:     regexp.MustCompile(regexPhoneSearch),
		searchEmail:     regexp.MustCompile(regexEmail),
		matchEmail:      matchRegex(regexEmail),
		matchPhoneIntl:  matchRegex(regexPhoneIntl),
		matchPhoneLocal: matchRegex(regexPhoneLocal),
		matchDomain:     matchRegex(regexDomain),
		matchHostname:   matchRegex(regexHostname),
		matchRootURL:    matchRegex(regexRootURL),
		match32BytesHex: matchRegex(regex32BytesHex),
		match64BytesHex: matchRegex(regex64BytesHex),
	}
}

func matchRegex(regex string) *regexp.Regexp {
	return regexp.MustCompile("^" + regex + "$")
}

func (r *regex) MatchEmail(s string) bool {
	return r.matchEmail.MatchString(s)
}
func (r *regex) MatchPhoneIntl(s string) bool {
	return r.matchPhoneIntl.MatchString(s)
}
func (r *regex) MatchPhoneLocal(s string) bool {
	return r.matchPhoneLocal.MatchString(s)
}
func (r *regex) MatchDomain(s string) bool {
	return r.matchDomain.MatchString(s)
}
func (r *regex) MatchHostname(s string) bool {
	return r.matchHostname.MatchString(s)
}
func (r *regex) MatchRootURL(s string) bool {
	return r.matchRootURL.MatchString(s)
}
func (r *regex) Match64BytesHex(s string) bool {
	return r.match64BytesHex.MatchString(s)
}
func (r *regex) MatchMD5String(s string) bool {
	return r.match32BytesHex.MatchString(s)
}

func (r *regex) SearchIPv4(s string) (results []string) {
	return r.searchIPv4.FindAllString(s, -1)
}
func (r *regex) SearchIPv6(s string) (results []string) {
	return r.searchIPv6.FindAllString(s, -1)
}
func (r *regex) SearchEmail(s string) (results []string) {
	return r.searchEmail.FindAllString(s, -1)
}
func (r *regex) SearchPhone(s string) (results []string) {
	return r.searchPhone.FindAllString(s, -1)
}
