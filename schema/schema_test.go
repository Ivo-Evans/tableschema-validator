package schema

import (
	"reflect"
	"testing"
)


func TestConstructSchema (t *testing.T) {
	expected := Schema{
		schemaSchema: `https://datapackage.org/profiles/2.0/tableschema.json`,
		SchemaOptions: SchemaOptions{
			Fields: Fields{
				StringFields: []StringField{
					{
						FieldBase: FieldBase{Name: "foo", fieldType: "string"}, 
						Constraints: StringConstraints{
							Required: true, 
							Enum: []string{"bar", "baz"},
						},
					},
				},
			},
		},
	} 

	got := MakeSchema(SchemaOptions{
		Fields: Fields{
			StringFields: []StringField{
				{
					FieldBase: FieldBase{Name: "foo"}, 
					Constraints: StringConstraints{
						Required: true, 
						Enum: []string{"bar", "baz"},
						// there are hidden 0-values here, like maxLength 0. A potential problem
					},
				},
			},
		},
	})

	

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nWanted %+v\n got %+v\n", expected, got)
	}
	
}