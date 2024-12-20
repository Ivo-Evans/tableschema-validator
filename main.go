package main

import (
	"encoding/json"
	"fmt"
	"tableschema-validator/schema"
)

func main() {
	mySchema := schema.MakeSchema(schema.SchemaOptions{
		Fields: schema.Fields{
			NumberFields: []schema.NumberField{
				{
					FieldBase: schema.FieldBase{Name: "baz"},
					Constraints: schema.NumberConstraints{
						Min: schema.MinConstraint{Selected: true, Value: 11},
					},
				},
			},
			StringFields: []schema.StringField{
				{
					FieldBase: schema.FieldBase{Name: "foo"},
					Constraints: schema.StringConstraints{
						Required: schema.RequiredConstraint{Selected: true, Value: true},
						Enum:     schema.EnumConstraint{Selected: true, Value: []string{"bar", "baz"}},
						Pattern:  schema.PatternConstraint{Selected: true, Value: ".+|$"},
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

	t, err := json.MarshalIndent(mySchema, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(t))
}
