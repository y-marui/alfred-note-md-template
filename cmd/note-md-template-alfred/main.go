// Command note-md-template-alfred is the Script Filter binary Alfred
// invokes (see workflow/info.plist). Alfred's Script Filter node runs it
// with the query typed after the "note" keyword as $1.
package main

import (
	"fmt"
	"os"

	"github.com/y-marui/alfred-note-md-template/internal/scriptfilter"
	"github.com/y-marui/alfred-note-md-template/internal/templatelist"
)

func main() {
	query := ""
	if len(os.Args) > 1 {
		query = os.Args[1]
	}
	writeResponse(dispatch(query))
}

// dispatch recovers from any panic in templatelist, so an unhandled failure
// still produces a visible Script Filter error item rather than
// empty/invalid output.
func dispatch(query string) (resp scriptfilter.Response) {
	defer func() {
		if r := recover(); r != nil {
			resp = errorResponse(fmt.Sprintf("%v", r))
		}
	}()
	return templatelist.List(query)
}

func writeResponse(resp scriptfilter.Response) {
	if err := resp.Write(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "note-md-template-alfred: writing response:", err)
		os.Exit(1)
	}
}

func errorResponse(message string) scriptfilter.Response {
	return scriptfilter.Response{
		Items: []scriptfilter.Item{
			{Title: "Workflow Error", Subtitle: message, Arg: message, Valid: scriptfilter.BoolPtr(false)},
		},
	}
}
