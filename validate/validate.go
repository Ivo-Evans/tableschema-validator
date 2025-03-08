// Package validate exposes the Validate function as well as various functions used by it.
// Validate takes a schema from [schema/schema.go] and something that implements `ReadAll()`,
// such as a `*csv.Reader` from Go's `encoding/csv`library, and returns a list of
// `RowValidationResult` structs.
package validate

import (
	"tableschema-validator/schema"
)

type Readable interface {
	ReadAll() (records [][]string, err error)
}

// A CellValidationResult is the `verdict` on a single cell-constraint combination. This means a single cell could have multiple
// CellValidationResults if it fails multiple validations (really? check this. Possibly rename). 
type CellValidationResult struct {
	header     string
	value      string
	constraint string
	reason     string
	isValid    bool
}

// A RowValidationResult is the 'verdict' on a single row. It includes the Validate package's internal representation of the row,
// as well as an `isValid` result which is false iff there is at least one item in the Failures slice. While a `CellValidationResult`
// may be produced for a valid or an invalid row, it will only exist in a `RowValidationResult` to indicate invalid data - if
// `CellValidationResult.isValid` is false.
type RowValidationResult struct {
	Original []string
	Parsed   map[string]string
	IsValid  bool
	Failures []CellValidationResult
}

func mapRowCellsToHeaders(headers []string, rawRow []string) map[string]string {
	row := make(map[string]string)

	for headerIndex, header := range headers {
		row[header] = rawRow[headerIndex]
	}

	return row
}

// validateRow is used for standard validations. It takes the raw row string, a map of input header to csv values, and then a row
// it then applies all validations to the row that aren't relational, i.e. don't depend on other rows.
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
		requiredValidationFailure, err := EnforceRequiredConstraint(numberField.Constraints.Required, numberField.Name, row[numberField.Name])
		if err != nil {
			return RowValidationResult{}, err
		}
		handleValidationResult(requiredValidationFailure)
	}

	return RowValidationResult{Original: rawRow, Parsed: row, IsValid: isValid, Failures: validationFailures}, nil
}

// validateColumns is used for validations which depend on comparong values from multiple cells in the same column. 
// Currently this only includes unique validations.
func validateColumns(schema schema.Schema, validatedRows *[]RowValidationResult) *[]RowValidationResult {
	for _, stringField := range schema.Fields.StringFields {
		if stringField.Constraints.Unique.Selected == true && stringField.Constraints.Unique.Value == true {
			EnforceUniqueConstraint(stringField.Constraints.Unique, stringField.Name, validatedRows)

		}
		continue
	}

	for _, numberField := range schema.Fields.NumberFields {
		if numberField.Constraints.Unique.Selected == true && numberField.Constraints.Unique.Value == true {
			EnforceUniqueConstraint(numberField.Constraints.Unique, numberField.Name, validatedRows)
		}
		continue
	}

	return validatedRows
}

// Validate takes a schema from [schema/schema.go] and something that implements `ReadAll()`,
// such as a `*csv.Reader` from Go's `encoding/csv`library, and returns a list of
// `RowValidationResult` structs.
func Validate(schema schema.Schema, sourceData Readable) ([]RowValidationResult, error) {
	data, err := sourceData.ReadAll()
	if err != nil {
		return nil, err
	}

	headers := data[0]

	var rowValidationResults []RowValidationResult

	for _, rawRow := range data[1:] {
		row := mapRowCellsToHeaders(headers, rawRow)
		rowValidationResult, err := validateRow(rawRow, row, schema)
		if err != nil {
			return rowValidationResults, err
		}
		rowValidationResults = append(rowValidationResults, rowValidationResult)
	}

	columnValidationResults := validateColumns(schema, &rowValidationResults)

	return *columnValidationResults, nil
}

// TODOs
// remaining checks
// package can be installed into another Go project
// Possibly add coerced values to output - though it could be a pain because there are no sum types, have to think about it. The answer is probably reflection and tags on a passed struct.
// you might want to have a think about, like, module boundaries and stuff to neaten up imports/ownership.
// interesting approach to a similar problem you have here https://www.reddit.com/r/golang/comments/1ijcaki/how_would_you_decodeencode_json_sum_types_in_go/
// godoc could be interesting too
