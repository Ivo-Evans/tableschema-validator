package schema

type Constraint [selection any] struct {
	// we use selected true/false to tell the difference between a selection's 0-value being deliberately set v.s. defaulting
	selected bool
	value selection
}

type requiredConstraint = Constraint[bool]
type uniqueContraint = Constraint[bool]
type patternConstraint = Constraint[string] // has to match XML schema
type enumConstraint = Constraint[[]string]  // min 1 item. All must be unique.
type minLengthConstraint = Constraint[int64]
type maxLengthConstraint = Constraint[int64]
type minConstraint = Constraint[int64]
type maxConstraint = Constraint[int64]

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