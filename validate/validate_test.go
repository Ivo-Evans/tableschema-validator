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
			},
		},
	})

	file, err := os.Open("../test-data/simple-example.csv")

	if err != nil {
		t.Error("Failed to read fixture CSV")
	}

	reader := csv.NewReader(file)

	expected := []RowValidationResult{
		{original: []string{"baz", "baz", "0"}, isValid: true, failures: nil},
		{original: []string{"bar", "luhrman", "2"}, isValid: true, failures: nil},
		{original: []string{"100", "antidisestablishmentarianism", "3"}, isValid: true, failures: nil},
		{original: []string{"", "qux", ""}, isValid: false, failures: []CellValidationResult{{header: "foo", constraint: "required", isValid: false, value: "", reason: "foo was marked as required, but not provided"}}},
	}

	got, err := Validate(schema, reader)

	if err != nil {
		t.Error("Failed to validate fixture CSV")
	}

	if diff := cmp.Diff(expected, got, cmp.AllowUnexported(RowValidationResult{}, CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

}
