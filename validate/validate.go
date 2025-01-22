package validate

import (
	"tableschema-validator/schema"
)

type Readable interface {
	ReadAll() (records [][]string, err error)
}

type CellValidationResult struct {
	header     string
	value      string
	constraint string
	reason     string
	isValid    bool
}

type RowValidationResult struct {
	original []string
	isValid  bool
	failures []CellValidationResult
}

func mapRowCellsToHeaders(headers []string, rawRow []string) map[string]string {
	row := make(map[string]string)

	for headerIndex, header := range headers {
		row[header] = rawRow[headerIndex]
	}

	return row
}

func validateRow(rawRow []string, row map[string]string, schema schema.Schema) (RowValidationResult, error) {
	isValid := true
	var validationFailures []CellValidationResult

	handleValidationResult := func(result CellValidationResult) {
		if !result.isValid {
			isValid = false
			validationFailures = append(validationFailures, result)
		}
	}

	for _, stringField := range schema.Fields.StringFields {
		dataTypeValidationFailure, err := EnforceStringConstraint()
		if err != nil {
			return RowValidationResult{}, err
		}
		handleValidationResult(dataTypeValidationFailure)
		requiredValidationFailure, err := EnforceRequiredConstraint(stringField.Constraints.Required, stringField.Name, row[stringField.Name])
		if err != nil {
			return RowValidationResult{}, err
		}
		handleValidationResult(requiredValidationFailure)
	}

	for _, numberField := range schema.Fields.NumberFields {
		dataTypeValidationFailure, err := EnforceNumberConstraint(numberField.Name, row[numberField.Name])
		if err != nil {
			return RowValidationResult{}, err
		}
		handleValidationResult(dataTypeValidationFailure)
	}

	return RowValidationResult{original: rawRow, isValid: isValid, failures: validationFailures}, nil
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
		rowValidationResult, err := validateRow(rawRow, row, schema)
		if err != nil {
			return validationResults, err
		}
		validationResults = append(validationResults, rowValidationResult)
	}

	return validationResults, nil
}

