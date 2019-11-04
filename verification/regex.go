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
	regexRootURL      = `\/[a-zA-Z0-9\-_/\+]+`
	regexMD5String    = `[a-fA-F0-9]{32}`
	regexSHA256String = `[a-fA-F0-9]{64}`
)

func buildSearchFn(regex string) func(s string) []string {
	return func(s string) []string {
		return regexp.MustCompile(regex).FindAllString(s, -1)
	}
}

func buildMatchFn(regex string) func(s string) bool {
	return func(s string) bool {
		return regexp.MustCompile("^" + regex + "$").MatchString(s)
	}
}

// Search functions
var (
	SearchIPv4  = buildSearchFn(regexIPv4)
	SearchIPv6  = buildSearchFn(regexIPv6)
	SearchEmail = buildSearchFn(regexEmail)
	SearchPhone = buildSearchFn(regexPhoneSearch)
)

// Match functions
var (
	MatchEmail      = buildMatchFn(regexEmail)
	MatchPhoneIntl  = buildMatchFn(regexPhoneIntl)
	MatchPhoneLocal = buildMatchFn(regexPhoneLocal)
	MatchDomain     = buildMatchFn(regexDomain)
	MatchHostname   = buildMatchFn(regexHostname)
	// TODO add tests
	MatchRootURL      = buildMatchFn(regexRootURL)
	MatchMD5String    = buildMatchFn(regexMD5String)
	MatchSHA256String = buildMatchFn(regexSHA256String)
)
