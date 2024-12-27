package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/xeipuuv/gojsonschema"
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

func TestMarshallSchemaToJSON(t *testing.T) {
	schema := MakeSchema(SchemaOptions{
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

	asJson, err := json.MarshalIndent(schema, "", "  ")

	if err != nil {
		t.Errorf("Failed to marshall schema to JSON with error %s", err.Error())
	}

	got := strings.TrimSpace(string(asJson))

	expected := strings.TrimSpace(`{
  "$schema": "https://datapackage.org/profiles/2.0/tableschema.json",
  "fields": [
    {
      "type": "string",
      "name": "foo",
      "constraints": {
        "required": true,
        "enum": [
          "bar",
          "baz"
        ]
      }
    },
    {
      "type": "string",
      "name": "bar",
      "constraints": {
        "required": true,
        "enum": [
          "bar",
          "baz"
        ],
        "minLength": 10
      }
    }
  ]
}
`)

	if got != expected {
		t.Errorf("\nWanted %s got %s", expected, got)
	}

	schemaLoader := gojsonschema.NewReferenceLoader("file://../schema.json")
	documentLoader := gojsonschema.NewStringLoader(got)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if err != nil {
		t.Errorf("Failed to validate schema, got %s", err.Error())
		panic(err.Error())
	}

	if !result.Valid() {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		t.Errorf("Generated an invalid tableschema. Generated %s", got)
	}

}
