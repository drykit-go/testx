package docparser

import (
	"strings"
)

func ParseDocLines(rawdoc string, repl map[string]string) []string {
	lines := []string{}
	if rawdoc == "" {
		return lines
	}

	// Set up strings replacer
	replArgs := []string{}
	for k, v := range repl {
		replArgs = append(replArgs, k, v)
	}
	replacer := strings.NewReplacer(replArgs...)

	// Parse and add new lines
	for _, l := range strings.Split(rawdoc, "\n") {
		lines = append(lines, replacer.Replace(l))
	}

	// Strip last line, always empty due to final '\n'
	return lines[:len(lines)-1]
}