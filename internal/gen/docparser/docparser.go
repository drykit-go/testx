package docparser

import (
	"strings"
)

// ParseDocLines reads the rawdoc got from go/doc parsing,
// adds a leading "// " to each line, applies the replacements,
// and returns a slice of strings for each line.
// The result can then be conveniently used in a go template.
//
// Example:
// 	ParseDocLines(
// 		"myFunc does this.\nIt is cool.\n",
// 		map[string]string{"myFunc": "MyFunc", "cool": "nice"}
// 	)
//
// 	// Output
// 	[]string{
// 		"// MyFunc does this.",
// 		"// It is nice.",
// 	}
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
