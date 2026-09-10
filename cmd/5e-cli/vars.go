package main

import (
	"regexp"
	"strings"
)

// A `$placeholder` in the data files marks a value the companion draws for us.
// It is rewritten into a Handlebars tag plus an "item/vars" declaration naming
// the preset to draw from — the vocabularies themselves live in the companion's
// config under `randoms`, so the value arrives with its options attached and the
// DM can re-pick it without retyping the prose around it.
var placeholder = regexp.MustCompile(`\$[a-zA-Z]\w*`)

// templated rewrites a data string's placeholders into a companion template and
// the vars that fill it. Vars are nil when there is nothing to draw, so an item
// with no placeholders carries no "item/vars".
//
// ponytail: a placeholder repeated in one string shares a single draw (it used
// to roll twice); suffix the id if a mod ever needs two independent rolls.
func templated(s string) (string, map[string]any) {
	var vars map[string]any
	body := placeholder.ReplaceAllStringFunc(s, func(match string) string {
		preset := kebab(match[1:])
		if vars == nil {
			vars = map[string]any{}
		}
		vars[preset] = map[string]any{"random": preset}
		return "{{ " + preset + " }}"
	})
	return body, vars
}

// kebab converts the data files' camelCase placeholder names to the kebab-case
// the companion's presets are named in, which is also what the UI labels them by.
func kebab(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('-')
			}
			r += 'a' - 'A'
		}
		b.WriteRune(r)
	}
	return b.String()
}
