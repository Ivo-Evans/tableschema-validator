package schema

type requiredConstraint = bool
type uniqueContraint = bool
type patternConstraint = string // has to match XML schema
type enumConstraint = []string  // min 1 item. All must be unique.
type minLengthConstraint int64
type maxLengthConstraint int64
type minConstraint int64
type maxConstraint int64

type StringConstraints struct {
	Required  requiredConstraint
	Unique    uniqueContraint
	Pattern   patternConstraint
	Enum      enumConstraint
	MinLength minLengthConstraint
	MaxLength maxLengthConstraint
}

type NumberConstraints struct {
	Required requiredConstraint
	Unique   uniqueContraint
	Min      minConstraint
	Max      maxConstraint
}

type BooleanConstraints struct {
	Required requiredConstraint
	Enum     enumConstraint
}

type ListConstraints struct {
	Required  requiredConstraint
	MinLength minLengthConstraint
	MaxLength maxLengthConstraint
}

type FieldBase struct {
	fieldType   string
	Name        string
	Title       string
	Description string
	Example     string
}

type StringField struct {
	FieldBase
	Constraints StringConstraints
}

type NumberField struct {
	FieldBase
	Constraints NumberConstraints
}

type BooleanField struct {
	FieldBase
	Constraints BooleanConstraints
}

type ListField struct {
	FieldBase
	Constraints ListConstraints
}

type Fields struct {
	StringFields  []StringField
	NumberFields  []NumberField
	BooleanFields []BooleanField
	ListFields    []ListField
}

type SchemaOptions struct {
	Fields Fields
}