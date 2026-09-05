// Package scriptfilter models Alfred's Script Filter JSON output format
// (https://www.alfredapp.com/help/workflows/inputs/script-filter/json/).
package scriptfilter

import (
	"encoding/json"
	"io"
)

// Item is one Alfred result row.
type Item struct {
	UID          string            `json:"uid,omitempty"`
	Title        string            `json:"title"`
	Subtitle     string            `json:"subtitle,omitempty"`
	Arg          string            `json:"arg,omitempty"`
	Valid        *bool             `json:"valid,omitempty"`
	Autocomplete string            `json:"autocomplete,omitempty"`
	Variables    map[string]string `json:"variables,omitempty"`
}

// Response is the top-level Script Filter JSON document.
type Response struct {
	Items []Item `json:"items"`
}

// Write encodes r as JSON to w, Alfred's expected Script Filter output.
func (r Response) Write(w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}

// BoolPtr is a small helper for the Valid field above, which must
// distinguish "absent" (Alfred defaults to valid) from an explicit false.
func BoolPtr(b bool) *bool {
	return &b
}
