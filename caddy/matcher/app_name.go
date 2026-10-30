package matcher

import "strings"

func AppName(serverName, domain string) string {
	name := strings.ToLower(strings.TrimSuffix(serverName, "."))
	suffix := "." + strings.ToLower(domain)
	if !strings.HasSuffix(name, suffix) {
		return ""
	}
	labels := strings.Split(strings.TrimSuffix(name, suffix), ".")
	return labels[len(labels)-1]
}
