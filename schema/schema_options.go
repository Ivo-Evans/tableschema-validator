package schema

// A Constraint object is a specification of a particular constraint
// If Selected is true, Value is read
// Value is then the constraint's value, e.g. for a Pattern constraint
// Value might be a regular expression. If Selected is true and Value
// is a 0 value, Value is taken to be the value. For instance, a
// minConstraint which was selected with a value of 0 would enforce
// that source data values were above 0.
type Constraint[selection any] struct {
	// we use Selected true/false to tell the difference between a selection's 0-value being deliberately set v.s. defaulting
	Selected bool
	Value    selection
}

type RequiredConstraint = Constraint[bool]
type UniqueContraint = Constraint[bool]
type PatternConstraint = Constraint[string] // has to match XML schema
type EnumConstraint = Constraint[[]string]  // min 1 item. All must be unique.
type MinLengthConstraint = Constraint[int64]
type MaxLengthConstraint = Constraint[int64]
type MinConstraint = Constraint[int64]
type MaxConstraint = Constraint[int64]

type MarshalableConstraintStruct struct {
	Required  bool     `json:"required,omitempty"`
	Unique    bool     `json:"unique,omitempty"`
	Pattern   string   `json:"pattern,omitempty"`
	Enum      []string `json:"enum,omitempty"`
	MinLength int64    `json:"minLength,omitempty"`
	MaxLength int64    `json:",omitempty"`
	Min       int64    `json:"min,omitempty"`
	Max       int64    `json:"max,omitempty"`
}

type StringConstraints struct {
	Required  RequiredConstraint  `json:"required"`
	Unique    UniqueContraint     `json:"unique"`
	Pattern   PatternConstraint   `json:"pattern"`
	Enum      EnumConstraint      `json:"enum"`
	MinLength MinLengthConstraint `json:"minLength"`
	MaxLength MaxLengthConstraint `json:"maxLength"`
}

type NumberConstraints struct {
	Required RequiredConstraint `json:""`
	Unique   UniqueContraint    `json:"unique"`
	Min      MinConstraint      `json:"min"`
	Max      MaxConstraint      `json:"max"`
}

type BooleanConstraints struct {
	Required RequiredConstraint `json:"required"`
	Enum     EnumConstraint     `json:"enum"`
}

type ListConstraints struct {
	Required  RequiredConstraint  `json:"required"`
	MinLength MinLengthConstraint `json:"minLength"`
	MaxLength MaxLengthConstraint `json:"maxLength"`
}

type FieldBase struct {
	FieldType   string `json:"type"`
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Example     string `json:"example,omitempty"`
}

type StringField struct {
	FieldBase
	Constraints StringConstraints `json:"constraints"`
}

type NumberField struct {
	FieldBase
	Constraints NumberConstraints `json:"constraints"`
}

type BooleanField struct {
	FieldBase
	Constraints BooleanConstraints `json:"constraints"`
}

type ListField struct {
	FieldBase
	Constraints ListConstraints `json:"constraints"`
}

// The fields, or columns, of the source data, that are to be included in the schema.
// Fields are split into their types rather than being a list like in the output json
// version of the schema or the csv header row because of the difficulty of modelling 
// sum types in Golang. An outcome of this is that as a user, you cannot control the 
// order of your field definitions.
type Fields struct {
	StringFields  []StringField
	NumberFields  []NumberField
	BooleanFields []BooleanField
	ListFields    []ListField
}

type SchemaOptions struct {
	Fields Fields `json:"fields"`
}
