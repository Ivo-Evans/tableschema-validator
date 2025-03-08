// Package util is a collection of general-purpose utility functions
package util

import (
	"fmt"
	"strings"
)

// CommaSeparatedList takes a slice of listItems and joins them by ", " without the word "and".
func CommaSeparatedList[listItem any](list []listItem) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(list)), ", "), "[]")
}
