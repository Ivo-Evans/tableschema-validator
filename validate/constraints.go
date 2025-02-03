package validate

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
	"tableschema-validator/schema"
	"tableschema-validator/util"
)

func EnforceStringConstraint() (CellValidationResult, error) {
	// the CSV comes through as a string, meaning every field can be interpreted as a string.
	// that makes this constraint a bit 'dumb', but I've included it for consistency so that
	// a field of any data type will have a validator for its datatype.
	return CellValidationResult{constraint: "String", isValid: true}, nil
}

func EnforceNumberConstraint(header string, field string) (CellValidationResult, error) {
	trimmed := strings.TrimSpace(field)
	validResponse := CellValidationResult{constraint: "Number", isValid: true}

	specialValues := []string{"NaN", "INF", "-INF"}
	if slices.Contains(specialValues, trimmed) {
		return validResponse, nil
	}

	// I find the use of a regex here a bit sus - but in the interests of time, since this is only a side project, I think it's worthwhile and robust enough
	isMatch, err := regexp.MatchString("^[+-]?\\d+\\.?\\d*(E[+-]?\\d+)?$", field)
	if err != nil {
		return CellValidationResult{}, err
	}

	if isMatch {
		return validResponse, nil
	}

	return CellValidationResult{constraint: "Number", isValid: false, header: header, value: field, reason: header + " was marked as a number, but its value " + field + " could not be parsed as a number"}, nil

}

func EnforceRequiredConstraint(requiredConstraint schema.Constraint[bool], header string, field string) (CellValidationResult, error) {
	validResponse := CellValidationResult{constraint: "required", isValid: true}
	// Why check for both Selected and Value? Selected tells us that the Value false is not to be interpreted as a 0 value bool - we can beliefe Value == false means the user has opted out
	if !(requiredConstraint.Selected && requiredConstraint.Value) {
		return validResponse, nil
	}

	if field == "" {
		return CellValidationResult{constraint: "required", isValid: false, header: header, value: field, reason: (header + " was marked as required, but not provided")}, nil
	} else {
		return validResponse, nil
	}
}

func EnforceUniqueConstraint(uniqueConstraint schema.Constraint[bool], header string, validatedRows *[]RowValidationResult) {
	// this could be optimised (for instance, we iterate through _every_ row for _each_ unique constraint, which is inefficient) but that isn't a priority right now

	dedupedValueIndices := make(map[string]int) // the index of the list time we saw a given value, keyed to the value
	dupes := make(map[string][]int)             // the indices of each time we've seen a given value, keyed to the value, if we saw it more than once. 
	// The code could probably be cleaned up if we collapsed these two maps into one. Then in the second for loop, we could just check if 
	// a given value had more than one index

	for index, row := range *validatedRows {
		value := row.Parsed[header]
		lastSeenIndex, hasBeenSeen := dedupedValueIndices[value]
		if hasBeenSeen {
			_, valueDupesExist := dupes[value]

			if !valueDupesExist {
				dupes[value] = append(dupes[value], lastSeenIndex)
			}

			dupes[value] = append(dupes[value], index)
		}
		dedupedValueIndices[value] = index
	}

	for duplicatedValue, indices := range dupes {
		for _, duplicatedValueIndex := range indices {
			newRow := (*validatedRows)[duplicatedValueIndex]

			failure := CellValidationResult{constraint: "unique", isValid: false, header: header, value: duplicatedValue, reason: header + " was marked as unique but its value " + duplicatedValue + " was found on rows " + util.CommaSeparatedList(indices) + " (this row: " + strconv.Itoa(duplicatedValueIndex) + ")"}

			newRow.Failures = append(newRow.Failures, failure)
			newRow.IsValid = false

			(*validatedRows)[duplicatedValueIndex] = newRow
		}
	}

	return
}
