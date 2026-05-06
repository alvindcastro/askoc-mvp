package privacy

import (
	"regexp"
	"strings"
)

var (
	emailPattern            = regexp.MustCompile(`[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}`)
	phoneSeparatedPattern   = regexp.MustCompile(`\b(?:\+?1[-.\s]?)?(?:\([0-9]{3}\)|[0-9]{3})[-.\s]+[0-9]{3}[-.\s]+[0-9]{4}\b`)
	phoneCompactPattern     = regexp.MustCompile(`\b(?:\+?1)?[0-9]{10}\b`)
	secretAssignmentPattern = regexp.MustCompile(`(?i)\b(password|passcode|token|api[_-]?key)\s*(?:is|=|:)\s*[^&\s]+`)
	realIDPattern           = regexp.MustCompile(`\b[0-9]{7,}\b`)
	spacePattern            = regexp.MustCompile(`[ \t]{2,}`)
)

func Redact(value string) string {
	value = secretAssignmentPattern.ReplaceAllStringFunc(value, redactSecretAssignment)
	value = emailPattern.ReplaceAllString(value, "[REDACTED_EMAIL]")
	value = phoneSeparatedPattern.ReplaceAllString(value, "[REDACTED_PHONE]")
	value = phoneCompactPattern.ReplaceAllString(value, "[REDACTED_PHONE]")
	value = realIDPattern.ReplaceAllString(value, "[REDACTED_ID]")
	return strings.TrimSpace(spacePattern.ReplaceAllString(value, " "))
}

func redactSecretAssignment(value string) string {
	for _, separator := range []string{" is ", "=", ":"} {
		if idx := strings.Index(strings.ToLower(value), separator); idx >= 0 {
			return value[:idx+len(separator)] + "[REDACTED_SECRET]"
		}
	}
	return "[REDACTED_SECRET]"
}
