package params

import "strings"

func parseEnviron(environ []string) (kv map[string]string) {
	kv = make(map[string]string, len(environ))
	for _, s := range environ {
		i := strings.Index(s, "=")
		if i == -1 {
			kv[s] = ""
			continue
		}

		key := s[:i]
		value := s[i+1:]
		kv[key] = value
	}
	return kv
}
