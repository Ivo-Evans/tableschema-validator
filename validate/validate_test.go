package validate

import (
	"encoding/csv"
	"os"
	"tableschema-validator/schema"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRequiredConstraint(t *testing.T) {
	unselectedConstraint := schema.RequiredConstraint{Selected: false, Value: true}
	selectedConstraint := schema.RequiredConstraint{Selected: true, Value: true}
	
	isValid, validationResult := enforceRequiredConstraint(unselectedConstraint, "example", "")
	expectedValidationResult := CellValidationFailure{}

	if !isValid {
		t.Error("Considered a missing value invalid given an unselected required constraint")
	}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationFailure{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	isValid, validationResult = enforceRequiredConstraint(unselectedConstraint, "example", "hi there")
	expectedValidationResult = CellValidationFailure{}
	if !isValid {
		t.Error("Considered a provided value invalid given an unselected required constraint")
	}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationFailure{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	isValid, validationResult = enforceRequiredConstraint(selectedConstraint, "example", "")
	expectedValidationResult = CellValidationFailure{constraint: "required", header: "example", value: "", reason: "example was marked as required, but not provided"}
	if isValid {
		t.Error("Considered a missing value valid given a selected required constraint")
	}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationFailure{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	isValid, validationResult = enforceRequiredConstraint(selectedConstraint, "example", "hi there")
	expectedValidationResult = CellValidationFailure{}
	if !isValid {
		t.Error("Considered a provided value invalid given a selected required constraint")
	}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationFailure{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}


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

	if (err != nil) {
		t.Error("Failed to read fixture CSV")
	}

	reader := csv.NewReader(file)

	expected := []RowValidationResult{
		{original: []string{"baz", "baz"}, isValid: true, failures: nil},
		{original: []string{"bar", "luhrman"}, isValid: true, failures: nil},
		{original: []string{"100", "antidisestablishmentarianism"}, isValid: true, failures: nil},
		{original: []string{"", "qux"}, isValid: false, failures: []CellValidationFailure{{header: "foo", constraint: "required", value: "", reason: "foo was marked as required, but not provided" }}},
	
	}

	got, err := Validate(schema, reader)

	if (err != nil) {
		t.Error("Failed to validate fixture CSV")
	}

	if diff := cmp.Diff(expected, got, cmp.AllowUnexported(RowValidationResult{},CellValidationFailure{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}


}