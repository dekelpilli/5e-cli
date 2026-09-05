package main

import (
	"encoding/json"
	"fmt"
)

// This file defines the JSON contract 5e-cli speaks when driven as a
// non-interactive helper (e.g. by sns-companion). A request arrives on stdin
// and a view-model is written to stdout.

// Request is the context handed to a command on stdin. Every field is
// optional; a command that needs no inputs can be run with no stdin at all.
//
// State holds the store collections this plugin declared with
// `store/collections` or `store/manual` on its config entry, already read for
// us — sns-companion owns the persistence, so we never touch its files and
// never learn which backend it is using. Each collection is left raw so a
// command can unmarshal it into whatever shape it expects; see collection.
type Request struct {
	Inputs  map[string]any             `json:"inputs"`
	Session map[string]any             `json:"session"`
	State   map[string]json.RawMessage `json:"state"`
}

// collection unmarshals one of the request's store collections into T. An
// absent collection is the zero value, matching sns-companion, where a
// collection nothing has written to reads as empty.
func collection[T any](r Request, name string) (T, error) {
	var out T
	raw, ok := r.State[name]
	if !ok || len(raw) == 0 {
		return out, nil
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return out, fmt.Errorf("could not read the %q collection: %w", name, err)
	}
	return out, nil
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

// Item is a single generated line, mapped by sns-companion to `:item/*` keys.
type Item struct {
	Title    string   `json:"title,omitempty"`
	Body     string   `json:"body"`
	Metadata []string `json:"metadata,omitempty"`
}

// Section groups related items under an optional heading.
type Section struct {
	Heading string `json:"heading,omitempty"`
	Items   []Item `json:"items"`
}

// Action is an optional UI button (label + event vector).
type Action struct {
	Label string   `json:"label"`
	Event []string `json:"event"`
}

// ViewModel is the friendly, un-namespaced result sns-companion consumes.
//
// Mutations are writes to apply, shaped {collection: {key: value}} with a null
// retracting a key. sns-companion applies them only once this output has
// validated, so a command that ends in an error changes nothing.
type ViewModel struct {
	Title     string                    `json:"title"`
	Subtitle  string                    `json:"subtitle,omitempty"`
	Sections  []Section                 `json:"sections,omitempty"`
	Actions   []Action                  `json:"actions,omitempty"`
	Mutations map[string]map[string]any `json:"mutations,omitempty"`
}

// CommandFunc is a single loot generator: it reads the request and produces a
// view-model (or an error, which becomes a non-zero exit).
type CommandFunc func(Request) (ViewModel, error)

// sectionOf wraps items in a single, heading-less section.
func sectionOf(items ...Item) []Section {
	return []Section{{Items: items}}
}

// vmText builds a view-model with a title and a single body line.
func vmText(title, body string) ViewModel {
	return ViewModel{Title: title, Sections: sectionOf(Item{Body: body})}
}

// staticVM adapts a fixed title/body pair into a CommandFunc.
func staticVM(title, body string) CommandFunc {
	return func(Request) (ViewModel, error) { return vmText(title, body), nil }
}

// affixItem renders an affix as an item: substituted description plus its point
// value, upgrade note, and affinities as tags.
func affixItem(a Affix) Item {
	return Item{
		Body:     processString(a.Description),
		Metadata: append([]string{a.PointValue, a.Upgrade}, a.Affinities...),
	}
}
