// Package schema exposes a factory function that the user can use to create a Schema object
// A Schema object can then be used by the validate package to validate some source data.
package schema

type Schema struct {
	SchemaSchema string `json:"$schema"`
	SchemaOptions
}

func (fields *Fields) insertFieldTypes() {
	// can't use the value from the range iteration because it's a shallow clone - mutation there won't mutate the receiver
	for i := range fields.StringFields {
		fields.StringFields[i].FieldType = "string"
	}

	for i := range fields.NumberFields {
		fields.NumberFields[i].FieldType = "number"
	}

	for i := range fields.BooleanFields {
		fields.BooleanFields[i].FieldType = "boolean"
	}

	for i := range fields.ListFields {
		fields.ListFields[i].FieldType = "list"
	}
}

// Takes a set of SchemaOptions and converts them into a valid Schema
func MakeSchema(options SchemaOptions) Schema {
	options.Fields.insertFieldTypes()

	return Schema{
		SchemaSchema:  `https://datapackage.org/profiles/2.0/tableschema.json`,
		SchemaOptions: options}
}
