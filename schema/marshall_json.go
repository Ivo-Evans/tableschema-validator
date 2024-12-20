package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func trimCharacters(array string, prefix string, suffix string) string {
	if strings.HasPrefix(array, prefix) {
		array = array[1:]
	}

	if strings.HasSuffix(array, suffix) {
		array = array[:len(array)-1]
	}
	return array
}

func (fields Fields) MarshalJSON() ([]byte, error) {
	// I think you could do the same thing here as in your constraintsMarshaller, with reflection
	var fieldStrings []string

	stringFieldbytes, err := json.Marshal(fields.StringFields)
	if err != nil {
		return nil, err
	}

	stringFields := string(stringFieldbytes)

	if stringFields != `null` {
		fieldStrings = append(fieldStrings, trimCharacters(stringFields, `[`, `]`))
	}

	numberFieldbytes, err := json.Marshal(fields.NumberFields)
	if err != nil {
		return nil, err
	}

	numberFields := string(numberFieldbytes)

	if numberFields != `null` {
		fieldStrings = append(fieldStrings, trimCharacters(numberFields, `[`, `]`))
	}

	booleanFieldbytes, err := json.Marshal(fields.BooleanFields)
	if err != nil {
		return nil, err
	}

	booleanFields := string(booleanFieldbytes)

	if booleanFields != `null` {
		fieldStrings = append(fieldStrings, trimCharacters(booleanFields, `[`, `]`))
	}

	listFieldbytes, err := json.Marshal(fields.ListFields)
	if err != nil {
		return nil, err
	}

	listFields := string(listFieldbytes)

	if listFields != `null` {
		fieldStrings = append(fieldStrings, trimCharacters(listFields, `[`, `]`))
	}

	marhsalled := fmt.Sprintf(`[%s]`, strings.Join(fieldStrings, `,`))

	return []byte(marhsalled), nil
}

func (constraint Constraint[any]) MarshalJSON() ([]byte, error) {
	if !constraint.Selected {
		return []byte(`null`), nil
	}
	return json.Marshal(constraint.Value)
}

func constraintsMarshaller[anyConstraintSet StringConstraints | NumberConstraints | BooleanConstraints | ListConstraints](constraints anyConstraintSet) ([]byte, error) {
	var fields []string

	val := reflect.ValueOf(constraints)

	for i := 0; i < val.NumField(); i++ {
		structKey := val.Type().Field(i).Name
		jsonKey := strings.ToLower(structKey[0:1]) + structKey[1:]

		constraint := val.Field(i).Interface()
		constraintMarshalled, err := json.Marshal(constraint)
		if err != nil {
			return nil, err
		}

		jsonValue := string(constraintMarshalled)

		if jsonValue != `null` {
			keyValuePair := fmt.Sprintf(`"%s": %s`, jsonKey, jsonValue)
			fields = append(fields, keyValuePair)
		}
	}

	marhsalled := fmt.Sprintf(`{%s}`, strings.Join(fields, `,`))

	return []byte(marhsalled), nil
}

func (constraints StringConstraints) MarshalJSON() ([]byte, error) {
	return constraintsMarshaller(constraints)
}

func (constraints NumberConstraints) MarshalJSON() ([]byte, error) {
	return constraintsMarshaller(constraints)
}

func (constraints BooleanConstraints) MarshalJSON() ([]byte, error) {
	return constraintsMarshaller(constraints)
}

func (constraints ListConstraints) MarshalJSON() ([]byte, error) {
	return constraintsMarshaller(constraints)
}
