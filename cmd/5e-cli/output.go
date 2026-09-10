package main

import "strings"

// This file defines the JSON contract 5e-cli speaks when driven as a
// non-interactive helper (e.g. by sns-companion). A request arrives on stdin
// and a view-model is written to stdout.
//
// The view-model keys are namespaced ("loot/title", "item/body"), which is the
// shape sns-companion consumes directly — the same one its in-process plugins
// return, so nothing is remapped on the way in.

// Request is the context handed to a command on stdin. Both fields are
// optional; a command that needs no inputs can be run with no stdin at all.
type Request struct {
	Inputs  map[string]any `json:"inputs"`
	Session map[string]any `json:"session"`
}

// str returns a string input, or "" when absent/not a string.
func (r Request) str(key string) string {
	if s, ok := r.Inputs[key].(string); ok {
		return s
	}
	return ""
}

// num returns a numeric input. JSON numbers decode to float64.
func (r Request) num(key string) (float64, bool) {
	switch v := r.Inputs[key].(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	default:
		return 0, false
	}
}

// strs returns a `list?` string input as a []string. Absent, non-array, or
// non-string elements are dropped rather than erroring, matching str/num's
// permissive style (an empty result reads the same as "not provided").
func (r Request) strs(key string) []string {
	arr, ok := r.Inputs[key].([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// Item is a single generated line. Body is a template rendered by sns-companion
// in the browser; with nothing to interpolate it renders as itself. Vars are the
// values it interpolates, either finished or declared as a draw the companion
// makes (see vars.go) — sent separately from the text so the DM can edit a value
// without retyping the prose.
type Item struct {
	Title    string         `json:"item/title,omitempty"`
	Body     string         `json:"item/body"`
	Metadata []string       `json:"item/metadata,omitempty"`
	Vars     map[string]any `json:"item/vars,omitempty"`
}

// Section groups related items under an optional heading.
type Section struct {
	Heading string `json:"section/heading,omitempty"`
	Items   []Item `json:"section/items"`
}

// Action is an optional UI button. Not a built ":action/event" vector: that
// carries the id this command was registered under in sns-companion's config,
// which this side has no way to know, so the adapter builds it from these three
// fields and routes the click back to the same command.
type Action struct {
	Label  string         `json:"label"`
	Action string         `json:"action"`
	Params map[string]any `json:"params,omitempty"`
}

// ViewModel is the result sns-companion consumes.
type ViewModel struct {
	Title    string    `json:"loot/title"`
	Subtitle string    `json:"loot/subtitle,omitempty"`
	Sections []Section `json:"loot/sections,omitempty"`
	Actions  []Action  `json:"loot/actions,omitempty"`
}

// CommandFunc is a single loot generator: it reads the request and produces a
// view-model (or an error, which becomes a non-zero exit).
type CommandFunc func(Request) (ViewModel, error)

// sectionOf wraps items in a single, heading-less section.
func sectionOf(items ...Item) []Section {
	return []Section{{Items: items}}
}

// newItem builds an item from a body that may carry `$placeholder` markers,
// declaring a var for each so the companion draws it.
func newItem(body string) Item {
	b, vars := templated(body)
	return Item{Body: b, Vars: vars}
}

// vmText builds a view-model with a title and a single body line.
func vmText(title, body string) ViewModel {
	return ViewModel{Title: title, Sections: sectionOf(newItem(body))}
}

// drawVM is a command whose whole body is one companion-drawn value per named
// preset — the utility rollers that exist only to pick a word.
func drawVM(title string, presets ...string) CommandFunc {
	tags := make([]string, len(presets))
	vars := make(map[string]any, len(presets))
	for i, p := range presets {
		tags[i] = "{{ " + p + " }}"
		vars[p] = map[string]any{"random": p}
	}
	body := strings.Join(tags, " ")
	return func(Request) (ViewModel, error) {
		return ViewModel{Title: title, Sections: sectionOf(Item{Body: body, Vars: vars})}, nil
	}
}

// staticVM adapts a fixed title/body pair into a CommandFunc.
func staticVM(title, body string) CommandFunc {
	return func(Request) (ViewModel, error) { return vmText(title, body), nil }
}

// affixItem renders an affix as an item: its description plus its point value,
// upgrade note, and affinities as tags.
func affixItem(a Affix) Item {
	item := newItem(a.Description)
	item.Metadata = append([]string{a.PointValue, a.Upgrade}, a.Affinities...)
	return item
}
