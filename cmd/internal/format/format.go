package format

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// PrintTable writes a tab-aligned table from items using a row-mapping function.
// headers defines the column names. rowFn converts each item into a string slice
// matching the header order.
func PrintTable[T any](w io.Writer, headers []string, items []T, rowFn func(T) []string) {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, strings.Join(headers, "\t"))
	for _, item := range items {
		fmt.Fprintln(tw, strings.Join(rowFn(item), "\t"))
	}
	tw.Flush()
}

// EmptyToDash returns s, or "-" if s is empty.
func EmptyToDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
