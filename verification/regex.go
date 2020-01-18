package verification

import "regexp"

const (
	regexIPv4         = `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`
	regexIPv6         = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	regexEmail        = `[a-zA-Z0-9-_.]+@[a-zA-Z0-9-_.]+\.[a-zA-Z]{2,10}`
	regexPhoneSearch  = `(\+|( *0 *0)){0,1}[0-9][0-9 ]{7}[0-9 ]*[0-9]`
	regexPhoneIntl    = `(\+|00)[0-9]{9,15}`
	regexPhoneLocal   = `([1-9][0-9]{4,13}|0[1-9][0-9]{3,12})`
	regexDomain       = `(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})`
	regexHostname     = `([a-zA-Z0-9]|[a-zA-Z0-9_][a-zA-Z0-9\-_]{0,61}[a-zA-Z0-9_])(\.([a-zA-Z0-9]|[a-zA-Z0-9_][a-zA-Z0-9\-_]{0,61}[a-zA-Z0-9]))*`
	regexRootURL      = `\/[a-zA-Z0-9\-_/\+]*`
	regexMD5String    = `[a-fA-F0-9]{32}`
	regexSHA256String = `[a-fA-F0-9]{64}`
)

// Search functions
func search(s, regex string) []string {
	return regexp.MustCompile(regex).FindAllString(s, -1)
}

func SearchIPv4(s string) []string {
	return search(s, regexIPv4)
}

func SearchIPv6(s string) []string {
	return search(s, regexIPv6)
}

func SearchEmail(s string) []string {
	return search(s, regexEmail)
}

func SearchPhone(s string) []string {
	return search(s, regexPhoneSearch)
}

// Match functions
func match(s, regex string) bool {
	return regexp.MustCompile("^" + regex + "$").MatchString(s)
}

func MatchEmail(s string) bool {
	return match(s, regexEmail)
}

func MatchPhoneIntl(s string) bool {
	return match(s, regexPhoneIntl)
}

func MatchPhoneLocal(s string) bool {
	return match(s, regexPhoneLocal)
}

func MatchDomain(s string) bool {
	return match(s, regexDomain)
}

func MatchHostname(s string) bool {
	return match(s, regexHostname)
}

func MatchRootURL(s string) bool {
	return match(s, regexRootURL)
}

func MatchMD5String(s string) bool {
	return match(s, regexMD5String)
}

func MatchSHA256String(s string) bool {
	return match(s, regexSHA256String)
}
