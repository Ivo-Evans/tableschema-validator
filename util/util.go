package util

import (
	"fmt"
	"strings"
)

func CommaSeparatedList[listItem any](list []listItem) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(list)), ", "), "[]")
}
