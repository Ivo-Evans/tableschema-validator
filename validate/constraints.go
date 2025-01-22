package validate

import (
	"regexp"
	"slices"
	"strings"
	"tableschema-validator/schema"
)

func EnforceStringConstraint() (CellValidationResult, error) {
	// the CSV comes through as a string, meaning every field can be interpreted as a string
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
	if !(requiredConstraint.Selected && requiredConstraint.Value) {
		return validResponse, nil
	}

	if field == "" {
		return CellValidationResult{constraint: "required", isValid: false, header: header, value: field, reason: (header + " was marked as required, but not provided")}, nil
	} else {
		return validResponse, nil
	}
}
