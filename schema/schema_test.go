package schema

import (
	"reflect"
	"testing"
)

func TestConstructSchema(t *testing.T) {
	got := MakeSchema(SchemaOptions{
		Fields: Fields{
			StringFields: []StringField{
				{
					FieldBase: FieldBase{Name: "foo"},
					Constraints: StringConstraints{
						Required: RequiredConstraint{Selected: true, Value: true},
						Enum:     EnumConstraint{Selected: true, Value: []string{"bar", "baz"}},
					},
				},
				{
					FieldBase: FieldBase{Name: "bar"},
					Constraints: StringConstraints{
						MinLength: MinLengthConstraint{Selected: true, Value: 10},
						Required:  RequiredConstraint{Selected: true, Value: true},
						Enum:      EnumConstraint{Selected: true, Value: []string{"bar", "baz"}},
					},
				},
			},
		},
	})

	expected := Schema{
		SchemaSchema: `https://datapackage.org/profiles/2.0/tableschema.json`,
		SchemaOptions: SchemaOptions{
			Fields: Fields{
				StringFields: []StringField{
					{
						FieldBase: FieldBase{Name: "foo", FieldType: "string"},
						Constraints: StringConstraints{
							Required:  RequiredConstraint{Selected: true, Value: true},
							Unique:    UniqueContraint{Selected: false, Value: false},
							Pattern:   PatternConstraint{Selected: false, Value: ``},
							Enum:      EnumConstraint{Selected: true, Value: []string{"bar", "baz"}},
							MinLength: MinLengthConstraint{Selected: false, Value: 0},
							MaxLength: MaxLengthConstraint{Selected: false, Value: 0},
						},
					},
					{
						FieldBase: FieldBase{Name: "bar", FieldType: "string"},
						Constraints: StringConstraints{
							Required:  RequiredConstraint{Selected: true, Value: true},
							Unique:    UniqueContraint{Selected: false, Value: false},
							Pattern:   PatternConstraint{Selected: false, Value: ``},
							Enum:      EnumConstraint{Selected: true, Value: []string{"bar", "baz"}},
							MinLength: MinLengthConstraint{Selected: true, Value: 10},
							MaxLength: MaxLengthConstraint{Selected: false, Value: 0},
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nWanted %+v\n got %+v\n", expected, got)
	}

}
