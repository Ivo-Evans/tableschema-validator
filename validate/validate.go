package validate

import (
	"tableschema-validator/schema"
)

type Readable interface {
	ReadAll() (records [][]string, err error)
}

type CellValidationFailure struct {
	header     string
	value      string
	constraint string
	reason     string
}

type RowValidationResult struct {
	original []string
	isValid  bool
	failures []CellValidationFailure
}

func enforceRequiredConstraint(requiredConstraint schema.Constraint[bool], header string, field string) (isValid bool, validationResult CellValidationFailure) {

	if !(requiredConstraint.Selected && requiredConstraint.Value) {
		return true, CellValidationFailure{}
	}

	if field == "" {
		return false, CellValidationFailure{constraint: "required", header: header, value: field, reason: (header + " was marked as required, but not provided")}
	} else {
		return true, CellValidationFailure{}
	}
}

func mapRowCellsToHeaders(headers []string, rawRow []string) (map[string]string) {
	row := make(map[string]string)

		for headerIndex, header := range headers {
			row[header] = rawRow[headerIndex]
		}

		return row
}

func validateRow (rawRow []string, row map[string]string, schema schema.Schema) (RowValidationResult) {
	isValid := true
	var validationFailures []CellValidationFailure

	for _, stringField := range schema.Fields.StringFields {
		isRequiredValid, requiredValidationFailure := enforceRequiredConstraint(stringField.Constraints.Required, stringField.Name, row[stringField.Name])
		if !isRequiredValid {
			isValid = false
			validationFailures = append(validationFailures, requiredValidationFailure)
		}
	}

	return RowValidationResult{original: rawRow, isValid: isValid, failures: validationFailures}
}

func Validate(schema schema.Schema, sourceData Readable) ([]RowValidationResult, error) {
	data, err := sourceData.ReadAll()
	if err != nil {
		return nil, err
	}

	headers := data[0]

	var validationResults []RowValidationResult

	for _, rawRow := range data[1:] {
		row := mapRowCellsToHeaders(headers, rawRow)
		rowValidationResult := validateRow(rawRow, row, schema)
		validationResults = append(validationResults, rowValidationResult)
	}

	return validationResults, nil
}

