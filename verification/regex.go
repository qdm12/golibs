package verification

import "regexp"

//nolint:lll
const (
	regex32BytesHex  = `[a-fA-F0-9]{32}`
	regex64BytesHex  = `[a-fA-F0-9]{64}`
	regexDomain      = `(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})`
	regexEmail       = `[a-zA-Z0-9-_.+]+@[a-zA-Z0-9-_.]+\.[a-zA-Z]{2,10}`
	regexHostname    = `([a-zA-Z0-9]|[a-zA-Z0-9_][a-zA-Z0-9\-_]{0,61}[a-zA-Z0-9_])(\.([a-zA-Z0-9]|[a-zA-Z0-9_][a-zA-Z0-9\-_]{0,61}[a-zA-Z0-9]))*`
	regexIPv4        = `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`
	regexIPv6        = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	regexPhoneIntl   = `(\+|00)[0-9]{9,15}`
	regexPhoneLocal  = `([1-9][0-9]{4,13}|0[1-9][0-9]{3,12})`
	regexPhoneSearch = `(\+|( *0 *0)){0,1}[0-9][0-9 ]{7}[0-9 ]*[0-9]`
	regexRootURL     = `\/[a-zA-Z0-9\-_/\+]*`
)

//nolint:gochecknoglobals
var (
	regexMatch32BytesHex = matchRegex(regex32BytesHex)
	regexMatch64BytesHex = matchRegex(regex64BytesHex)
	regexMatchDomain     = matchRegex(regexDomain)
	regexMatchEmail      = matchRegex(regexEmail)
	regexMatchHostname   = matchRegex(regexHostname)
	regexMatchPhoneIntl  = matchRegex(regexPhoneIntl)
	regexMatchPhoneLocal = matchRegex(regexPhoneLocal)
	regexMatchRootURL    = matchRegex(regexRootURL)
	regexSearchEmail     = regexp.MustCompile(regexEmail)
	regexSearchIPv4      = regexp.MustCompile(regexIPv4)
	regexSearchIPv6      = regexp.MustCompile(regexIPv6)
	regexSearchPhone     = regexp.MustCompile(regexPhoneSearch)
)

func matchRegex(regex string) *regexp.Regexp {
	return regexp.MustCompile("^" + regex + "$")
}

func MatchEmail(s string) bool {
	return regexMatchEmail.MatchString(s)
}
func MatchPhoneIntl(s string) bool {
	return regexMatchPhoneIntl.MatchString(s)
}
func MatchPhoneLocal(s string) bool {
	return regexMatchPhoneLocal.MatchString(s)
}
func MatchDomain(s string) bool {
	return regexMatchDomain.MatchString(s)
}
func MatchHostname(s string) bool {
	return regexMatchHostname.MatchString(s)
}
func MatchRootURL(s string) bool {
	return regexMatchRootURL.MatchString(s)
}
func Match64BytesHex(s string) bool {
	return regexMatch64BytesHex.MatchString(s)
}
func MatchMD5String(s string) bool {
	return regexMatch32BytesHex.MatchString(s)
}

func SearchIPv4(s string) (results []string) {
	return regexSearchIPv4.FindAllString(s, -1)
}
func SearchIPv6(s string) (results []string) {
	return regexSearchIPv6.FindAllString(s, -1)
}
func SearchEmail(s string) (results []string) {
	return regexSearchEmail.FindAllString(s, -1)
}
func SearchPhone(s string) (results []string) {
	return regexSearchPhone.FindAllString(s, -1)
}
