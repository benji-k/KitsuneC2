package web

import "regexp"

func isAlphaNumeric(s string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, s)
	return match
}

func isValidIpAdress(s string) bool {
	ipv4Pattern := `\b(?:(?:2[0-5][0-5]|1?\d{1,2})\.){3}(?:2[0-5][0-5]|1?\d{1,2})\b`
	ipv6Pattern := `\b(?:[A-Fa-f0-9]{1,4}:){7}[A-Fa-f0-9]{1,4}\b`

	ipv4Match, _ := regexp.MatchString(ipv4Pattern, s)
	if ipv4Match {
		return true
	}

	ipv6Match, _ := regexp.MatchString(ipv6Pattern, s)
	return ipv6Match
}

func isValidDomain(s string) bool {
	domainPattern := `^([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(domainPattern, s)
	return match
}

func isValidPort(p int) bool {
	return (p > 0) && (p <= 65535)
}

func isValidOs(s string) bool {
	validOs := []string{"linux", "windows", "aix", "android", "darwin", "dragonfly", "freebsd", "illumos", "ios", "js", "netbsd", "plan9", "solaris"}
	for _, v := range validOs {
		if v == s {
			return true
		}
	}
	return false
}

func isValidArch(s string) bool {
	validArch := []string{"386", "amd64", "arm", "arm64", "mips", "mips64", "mips64le", "mipsle", "ppc64", "ppc64le", "riscv64", "s390x", "wasm"}
	for _, v := range validArch {
		if v == s {
			return true
		}
	}
	return false
}
