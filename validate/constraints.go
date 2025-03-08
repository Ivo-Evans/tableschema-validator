package validate

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
	"tableschema-validator/schema"
	"tableschema-validator/util"
)

// EnforceStringConstraint reports whether a cell can be interpreted as a string.
// Since the raw data comes through as a string, this function marks every value
// as being able to be interpreted as a string. It is included for consistency with
// other data-type constraints, e.g. EnforceNumberConstraint, EnforceListConstraint.
func EnforceStringConstraint() (CellValidationResult, error) {
	return CellValidationResult{constraint: "String", isValid: true}, nil
}

// EnforceNumberConstraint reports whether a cell can be interpreted as a number,
// defined [here](https://datapackage.org/standard/table-schema/#number). There are
// a number of edge cases covered in the tests for this function. 
func EnforceNumberConstraint(header string, field string) (CellValidationResult, error) {
	trimmed := strings.TrimSpace(field)
	validResponse := CellValidationResult{constraint: "Number", isValid: true}

	specialValues := []string{"NaN", "INF", "-INF"}
	if slices.Contains(specialValues, trimmed) {
		return validResponse, nil
	}

	// I find the use of a regex here a bit sus - but in the interests of time, since this is only a side project, I think it's worthwhile and robust enough. There are relatively thorough tests for it. 
	isMatch, err := regexp.MatchString("^[+-]?\\d+\\.?\\d*(E[+-]?\\d+)?$", field)
	if err != nil {
		return CellValidationResult{}, err
	}

	if isMatch {
		return validResponse, nil
	}

	return CellValidationResult{constraint: "Number", isValid: false, header: header, value: field, reason: header + " was marked as a number, but its value " + field + " could not be parsed as a number"}, nil

}

// EnforceRequiredConstraint reports whether a cell is both required and absent.
// if the cell is both required and absent, EnforceRequiredConstraint marks the cell
// as invalid; if the cell is either not required, or has a value, it is reported as
// valid. 
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

// EnfoceUniqueConstraint applies a single column's unique constraint to the passed slice of RowValidationResult items, mutating the slice.
// Duplicates are marked as invalid, and the address of each invalid row is inserted into each invalid row's `RowValidationResult`.
func EnforceUniqueConstraint(uniqueConstraint schema.Constraint[bool], header string, validatedRows *[]RowValidationResult) {
	// this could be optimised (for instance, we iterate through _every_ row for _each_ unique constraint, which is inefficient) but that isn't a priority until there's a need to optimise

	uniqueValueIndices := make(map[string][]int)

	for index, row := range *validatedRows {
		value := row.Parsed[header]
		uniqueValueIndices[value] = append(uniqueValueIndices[value], index)
	}

	for sourceValue, sourceIndices := range uniqueValueIndices {
		isValueDuplicated := len(sourceIndices) > 1

		if !isValueDuplicated {
			continue
		}

		for _, sourceIndexOfDuplicate := range sourceIndices {
			newRow := (*validatedRows)[sourceIndexOfDuplicate]

			reason := header + " was marked as unique but its value " + sourceValue + " was found on rows " + util.CommaSeparatedList(sourceIndices) + " (this row: " + strconv.Itoa(sourceIndexOfDuplicate) + ")"

			failure := CellValidationResult{constraint: "unique", isValid: false, header: header, value: sourceValue, reason: reason}

			newRow.Failures = append(newRow.Failures, failure)
			newRow.IsValid = false

			(*validatedRows)[sourceIndexOfDuplicate] = newRow
		}
	}

	return
}
