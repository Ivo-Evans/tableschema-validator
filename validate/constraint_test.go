package validate

import (
	"tableschema-validator/schema"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStringConstraint(t *testing.T) {
	validationResult, err := EnforceStringConstraint()
	if err != nil {
		t.Errorf("Error enforcing string constraint")
	}
	expectedValidationResult := CellValidationResult{constraint: "String", isValid: true}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

func TestNumberConstraint(t *testing.T) {
	validationResult, err := EnforceNumberConstraint("foo", "1000")
	if err != nil {
		t.Errorf("Error enforcing number constraint")
	}
	expectedValidationResult := CellValidationResult{constraint: "Number", isValid: true}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	validationResult, err = EnforceNumberConstraint("foo", "1.1")
	if err != nil {
		t.Errorf("Error enforcing number constraint")
	}
	expectedValidationResult = CellValidationResult{constraint: "Number", isValid: true}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	validationResult, err = EnforceNumberConstraint("foo", "-46")
	if err != nil {
		t.Errorf("Error enforcing number constraint")
	}
	expectedValidationResult = CellValidationResult{constraint: "Number", isValid: true}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	validationResult, err = EnforceNumberConstraint("foo", "+8")
	if err != nil {
		t.Errorf("Error enforcing number constraint")
	}
	expectedValidationResult = CellValidationResult{constraint: "Number", isValid: true}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	validationResult, err = EnforceNumberConstraint("foo", "-61.9E+56")
	if err != nil {
		t.Errorf("Error enforcing number constraint")
	}
	expectedValidationResult = CellValidationResult{constraint: "Number", isValid: true}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	validationResult, err = EnforceNumberConstraint("foo", "NaN")
	if err != nil {
		t.Errorf("Error enforcing number constraint")
	}
	expectedValidationResult = CellValidationResult{constraint: "Number", isValid: true}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	validationResult, err = EnforceNumberConstraint("foo", "INF")
	if err != nil {
		t.Errorf("Error enforcing number constraint")
	}
	expectedValidationResult = CellValidationResult{constraint: "Number", isValid: true}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	validationResult, err = EnforceNumberConstraint("foo", "-INF")
	if err != nil {
		t.Errorf("Error enforcing number constraint")
	}
	expectedValidationResult = CellValidationResult{constraint: "Number", isValid: true}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	// more than one + or -
	validationResult, err = EnforceNumberConstraint("foo", "++33333")
	if err != nil {
		t.Errorf("Error enforcing number constraint")
	}
	expectedValidationResult = CellValidationResult{constraint: "Number", isValid: false, header: "foo", value: "++33333", reason: "foo was marked as a number, but its value ++33333 could not be parsed as a number"}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	// exponent sign with no exponent value
	validationResult, err = EnforceNumberConstraint("foo", "-61.9E")
	if err != nil {
		t.Errorf("Error enforcing number constraint")
	}
	expectedValidationResult = CellValidationResult{constraint: "Number", isValid: false, header: "foo", value: "-61.9E", reason: "foo was marked as a number, but its value -61.9E could not be parsed as a number"}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	// non-number
	validationResult, err = EnforceNumberConstraint("foo", "foo")
	if err != nil {
		t.Errorf("Error enforcing number constraint")
	}
	expectedValidationResult = CellValidationResult{constraint: "Number", isValid: false, header: "foo", value: "foo", reason: "foo was marked as a number, but its value foo could not be parsed as a number"}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	// non-integer exponent
	validationResult, err = EnforceNumberConstraint("foo", "-61.9E+56.5")
	if err != nil {
		t.Errorf("Error enforcing number constraint")
	}
	expectedValidationResult = CellValidationResult{constraint: "Number", isValid: false, header: "foo", value: "-61.9E+56.5", reason: "foo was marked as a number, but its value -61.9E+56.5 could not be parsed as a number"}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

func TestRequiredConstraint(t *testing.T) {
	unselectedConstraint := schema.RequiredConstraint{Selected: false, Value: true}
	selectedConstraint := schema.RequiredConstraint{Selected: true, Value: true}

	validationResult, err := EnforceRequiredConstraint(unselectedConstraint, "example", "")
	if err != nil {
		t.Errorf("Error enforcing required constraint")
	}
	expectedValidationResult := CellValidationResult{constraint: "required", isValid: true}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	validationResult, err = EnforceRequiredConstraint(unselectedConstraint, "example", "hi there")
	if err != nil {
		t.Errorf("Error enforcing required constraint")
	}
	expectedValidationResult = CellValidationResult{constraint: "required", isValid: true}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	validationResult, err = EnforceRequiredConstraint(selectedConstraint, "example", "")
	if err != nil {
		t.Errorf("Error enforcing required constraint")
	}
	expectedValidationResult = CellValidationResult{constraint: "required", isValid: false, header: "example", value: "", reason: "example was marked as required, but not provided"}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	validationResult, err = EnforceRequiredConstraint(selectedConstraint, "example", "hi there")
	if err != nil {
		t.Errorf("Error enforcing required constraint")
	}
	expectedValidationResult = CellValidationResult{constraint: "required", isValid: true}
	if diff := cmp.Diff(expectedValidationResult, validationResult, cmp.AllowUnexported(CellValidationResult{})); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}
