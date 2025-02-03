package validate

import (
	"encoding/csv"
	"os"
	"tableschema-validator/schema"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValidate(t *testing.T) {
	schema := schema.MakeSchema(schema.SchemaOptions{
		Fields: schema.Fields{
			StringFields: []schema.StringField{
				{
					FieldBase: schema.FieldBase{Name: "foo"},
					Constraints: schema.StringConstraints{
						Required: schema.RequiredConstraint{Selected: true, Value: true},
						Enum:     schema.EnumConstraint{Selected: true, Value: []string{"bar", "baz"}},
						Unique:   schema.UniqueContraint{Selected: true, Value: true},
					},
				},
				{
					FieldBase: schema.FieldBase{Name: "bar"},
					Constraints: schema.StringConstraints{
						MinLength: schema.MinLengthConstraint{Selected: true, Value: 10},
						Required:  schema.RequiredConstraint{Selected: true, Value: true},
						Enum:      schema.EnumConstraint{Selected: true, Value: []string{"bar", "baz"}},
					},
				},
				{
					FieldBase: schema.FieldBase{Name: "php"},
					Constraints: schema.StringConstraints{
						Required: schema.RequiredConstraint{Selected: true, Value: true},
					},
				},
			},
		},
	})

	file, err := os.Open("../test-data/simple-example.csv")

	if err != nil {
		t.Error("Failed to read fixture CSV")
	}

	reader := csv.NewReader(file)

	expected := []RowValidationResult{
		{Original: []string{"baz", "baz", "0"}, Parsed: map[string]string{"bar": "baz", "foo": "baz", "php": "0"}, IsValid: false, Failures: []CellValidationResult{
			{header: "foo", value: "baz", constraint: "unique", isValid: false, reason: "foo was marked as unique but its value baz was found on rows 0, 4 (this row: 0)"},
		}},
		{Original: []string{"bar", "luhrman", "2"}, Parsed: map[string]string{"bar": "luhrman", "foo": "bar", "php": "2"}, IsValid: true, Failures: nil},
		{Original: []string{"100", "antidisestablishmentarianism", "3"}, Parsed: map[string]string{"bar": "antidisestablishmentarianism", "foo": "100", "php": "3"}, IsValid: true, Failures: nil},
		{Original: []string{"", "qux", ""}, Parsed: map[string]string{"bar": "qux", "foo": "", "php": ""}, IsValid: false, Failures: []CellValidationResult{
			{header: "foo", constraint: "required", isValid: false, value: "", reason: "foo was marked as required, but not provided"},
			{header: "php", constraint: "required", isValid: false, value: "", reason: "php was marked as required, but not provided"},
		}},
		{Original: []string{"baz", "ghgh1010101010101", "4"}, Parsed: map[string]string{"bar": "ghgh1010101010101", "foo": "baz", "php": "4"}, IsValid: false, Failures: []CellValidationResult{
			{header: "foo", value: "baz", constraint: "unique", isValid: false, reason: "foo was marked as unique but its value baz was found on rows 0, 4 (this row: 4)"},
		}}}

	got, err := Validate(schema, reader)

	if err != nil {
		t.Error("Failed to validate fixture CSV")
	}

	if diff := cmp.Diff(expected, got, cmp.AllowUnexported(RowValidationResult{}, CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

}
