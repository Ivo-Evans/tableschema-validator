package schema

import (
	"reflect"
	"testing"
)


func TestConstructSchema (t *testing.T) {
	got := MakeSchema(SchemaOptions{
		Fields: Fields{
			StringFields: []StringField{
				{
					FieldBase: FieldBase{Name: "foo"}, 
					Constraints: StringConstraints{
						Required: requiredConstraint{selected: true, value: true}, 
						Enum: enumConstraint{selected: true, value: []string{"bar", "baz"}},
					},
				},
				{
					FieldBase: FieldBase{Name: "bar"}, 
					Constraints: StringConstraints{
						MinLength: minLengthConstraint{selected: true, value: 10},
						Required: requiredConstraint{selected: true, value: true}, 
						Enum: enumConstraint{selected: true, value: []string{"bar", "baz"}},
					},
				},
			},
		},
	})

	expected := Schema{
		schemaSchema: `https://datapackage.org/profiles/2.0/tableschema.json`,
		SchemaOptions: SchemaOptions{
			Fields: Fields{
				StringFields: []StringField{
					{
						FieldBase: FieldBase{Name: "foo", fieldType: "string"}, 
						Constraints: StringConstraints{
							Required: requiredConstraint{selected: true, value: true}, 
							Unique: uniqueContraint{selected: false, value: false},
							Pattern: patternConstraint{selected: false, value: ``},
							Enum: enumConstraint{selected: true, value: []string{"bar", "baz"}},
							MinLength: minLengthConstraint{selected: false, value: 0},
							MaxLength: maxLengthConstraint{selected: false, value: 0},
						},
					},
					{
						FieldBase: FieldBase{Name: "bar", fieldType: "string"}, 
						Constraints: StringConstraints{
							Required: requiredConstraint{selected: true, value: true}, 
							Unique: uniqueContraint{selected: false, value: false},
							Pattern: patternConstraint{selected: false, value: ``},
							Enum: enumConstraint{selected: true, value: []string{"bar", "baz"}},
							MinLength: minLengthConstraint{selected: true, value: 10},
							MaxLength: maxLengthConstraint{selected: false, value: 0},
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