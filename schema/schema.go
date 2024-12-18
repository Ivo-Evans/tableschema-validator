package schema

type Schema struct {
	schemaSchema string
	SchemaOptions
}

func (fields *Fields) insertFieldTypes() {
	// can't use the value from the range iteration because it's a shallow clone - mutation there won't mutate the receiver
	for i := range fields.StringFields {
		fields.StringFields[i].fieldType = "string"
	}

	for i := range fields.NumberFields {
		fields.NumberFields[i].fieldType = "number"
	}

	for i := range fields.BooleanFields {
		fields.BooleanFields[i].fieldType = "boolean"
	}

	for i := range fields.ListFields {
		fields.ListFields[i].fieldType = "list"
	}
}

func MakeSchema(options SchemaOptions) Schema {
	options.Fields.insertFieldTypes()

	return Schema{
		schemaSchema:  `https://datapackage.org/profiles/2.0/tableschema.json`,
		SchemaOptions: options}
}
